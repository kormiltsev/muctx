package muctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMuctx(t *testing.T) {
	mu := New()
	assert.NotNil(t, mu)
}

func TestLock(t *testing.T) {
	mu := New()
	assert.True(t, mu.Lock())
}

func TestUnlock(t *testing.T) {
	mu := New()
	require.True(t, mu.Lock())
	assert.True(t, mu.Unlock())
}

func TestQueue(t *testing.T) {
	mu := New()
	mu.Lock()

	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go func(ch chan int, mu *Muctx, i int) {
			mu.Lock()
			ch <- i
			mu.Unlock()
		}(ch, mu, i)
		time.Sleep(time.Second)
	}
	time.Sleep(500 * time.Millisecond)
	mu.Unlock()

	count := 0
	for count < 3 {
		// require same order. 0->1->2->....
		require.Equal(t, count, <-ch)
		count++
	}
}

func TestUnlockMulti(t *testing.T) {
	req := require.New(t)

	mu := New()
	mu2 := New()
	req.True(mu.Lock())
	req.False(mu2.Unlock())
	req.True(mu2.Lock())

	// f := mu.Lock
	// req.False(dedfunc(f))

	// as.Never(func() bool { return mu.Lock() }, time.Duration(2*time.Second), time.Second)
	req.True(mu.Unlock())
	req.True(mu.Lock())
	req.True(mu.Unlock())
	req.False(mu.Unlock())
	req.True(mu2.Unlock())
}

func TestTryCtx(t *testing.T) {
	req := require.New(t)

	mu := New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()

	mu.Lock()
	req.False(mu.LockTryCtx(ctx))
}

func TestTry(t *testing.T) {
	req := require.New(t)

	mu := New()

	mu.Lock()
	req.False(mu.LockTry())

	req.True(mu.Unlock())
	req.True(mu.LockTry())

}

/// not working??
// func dedfunc(f func() bool) bool {
// 	res := make(chan bool)
// 	go func(chan bool) {
// 		time.Sleep(time.Second)
// 		a := f()
// 		res <- a
// 	}(res)
// 	select {
// 	case <-time.After(1500 * time.Millisecond):
// 		return false
// 	case b := <-res:
// 		return b
// 	}
// }
