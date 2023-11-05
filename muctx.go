package muctx

import (
	"context"
	"time"
)

// Muctx supports Try method
type Muctx struct {
	Much     chan struct{}
	Chqueue  chan chan struct{}
	FirstJob *job
	LastJob  *job
}

type job struct {
	handchan chan struct{}
	next     *job
}

// New returns mu
func New() *Muctx {
	mx := Muctx{Much: make(chan struct{}, 1), Chqueue: make(chan chan struct{})}
	go mx.queue()
	return &mx
}

// Lock set this mu as locked
func (muctx *Muctx) Lock() bool {
	ctx := context.Background()
	return muctx.LockTryCtx(ctx)
}

// LockTry returns false if mu is locked now
func (muctx *Muctx) LockTry() bool {
	ctxTO, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()
	return muctx.LockTryCtx(ctxTO)
}

// LockTryCtx retiurns false if mu not locked till contect Done. Returns true if success
func (muctx *Muctx) LockTryCtx(ctx context.Context) bool {
	ch := make(chan struct{})
	muctx.Chqueue <- ch
	select {
	case <-ctx.Done():
		close(ch)
		return false
	case ch <- struct{}{}:
		close(ch)
		return true
	}
}

// Unlock releases mu
func (muctx *Muctx) Unlock() bool {
	select {
	case <-muctx.Much:
		return true
	default:
		return false
	}
}

// queue implements queue
func (muctx *Muctx) queue() {
	for {
		select {
		case newch := <-muctx.Chqueue:
			newjob := &job{handchan: newch}
			if muctx.FirstJob == nil {
				muctx.FirstJob = newjob
				muctx.LastJob = newjob
			} else {
				muctx.LastJob.next = newjob
				muctx.LastJob = newjob
			}
		case muctx.Much <- struct{}{}:

			if !muctx.nextjob() {
				<-muctx.Much
			}
		}
	}
}

// nextjob implements Next
func (muctx *Muctx) nextjob() bool {

	if muctx.FirstJob == nil {
		return false
	}

	_, ok := <-muctx.FirstJob.handchan
	if !ok {
		muctx.FirstJob = muctx.FirstJob.next
		return muctx.nextjob()
	}
	muctx.FirstJob = muctx.FirstJob.next
	return true
}
