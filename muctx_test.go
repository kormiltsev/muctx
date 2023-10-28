package muctx_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kormiltsev/library/muctx"
)

func TestLock(t *testing.T) {
	var testCases = []struct {
		locked     bool
		ctx        context.Context
		expectBool bool
	}{
		{false, context.Background, true},
		{true, context.WithTimeout(ctx, time.Duration(time.Second*2)), false},
		{true, context.WithTimeout(ctx, time.Duration(time.Second*10)), true},
	}
	for _, tc := range testCases {
		t.Run(tc.param, func(t *testing.T) {
			muctx := New()
			if tc.locked {
				muctx.Lock()
			}

go func(){
	muctx.Unlock()
}

			assert.Equal(t, tc.expectBool, muctx.Lock(tc.ctx))
		})
	}
}

func main() {
	ctx := context.Background()

	a := muctx.New()

	ctx2, cancel := context.WithTimeout(ctx, time.Duration(time.Second*2))
	defer cancel()

	a.Lock(ctx)

	fmt.Println("a locked, try lock b w/o ctx")

	b := muctx.New()
	b.Lock(context.Background())
	fmt.Println("b locked")
	b.Unlock()
	fmt.Println("b unlocked, try lock a with ctxTimeout 2 sec")

	if a.Lock(ctx2) {
		fmt.Println("error: can't lock a with ctxTimeout")
		os.Exit(1)
	}

	fmt.Println("a lock cancel (context done)\nTry to unlock a")

	a.Unlock()

	fmt.Println("a unlocked")

	ctx3, cancel3 := context.WithTimeout(ctx, time.Duration(time.Second))
	defer cancel3()

	if !a.Lock(ctx3) {
		fmt.Println("error: can't lock a")
		os.Exit(1)
	}

	a.Unlock()

	fmt.Println("DONE")

}
