// test time is more than 4 seconds

package muctx

import (
	"context"
	"time"
)

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

func New() *Muctx {
	mx := Muctx{Much: make(chan struct{}, 1), Chqueue: make(chan chan struct{})}
	go mx.queue()
	return &mx
}

func (muctx *Muctx) Lock() bool {
	ctx := context.Background()
	return muctx.LockTryCtx(ctx)
}

func (muctx *Muctx) LockTry() bool {
	ctxTO, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()
	return muctx.LockTryCtx(ctxTO)
}

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

func (muctx *Muctx) Unlock() bool {
	select {
	case <-muctx.Much:
		return true
	default:
		return false
	}
}

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

// type muctx struct {
// 	mu   sync.Mutex
// 	chmu chan struct{}
// }

// // New returns new mutex with context
// func New(mu ...sync.Mutex) *muctx {
// 	if len(mu) == 1 {
// 		return &muctx{
// 			mu:   mu[0],
// 			chmu: make(chan struct{}, 1),
// 		}
// 	}
// 	return &muctx{
// 		chmu: make(chan struct{}, 1),
// 	}
// }

// // Lock returns true if mutex set. False is context Done. Will lock process till one of this results
// func (mx *muctx) Lock(ctx context.Context) bool {
// 	select {
// 	case <-ctx.Done():
// 		return false
// 	case mx.chmu <- struct{}{}:
// 		mx.mu.Lock()
// 		return true
// 	}
// }

// // Unlock release mutex
// func (mx *muctx) Unlock() {
// 	<-mx.chmu
// 	mx.mu.Unlock()
// }
