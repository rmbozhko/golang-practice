package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

// execute runs fn in a separate goroutine
// and waits for the result unless canceled.
func execute(ctx context.Context, fn func() int) (int, error) {
	ch := make(chan int, 1)

	go func() {
		ch <- fn()
	}()

	select {
	case res := <-ch:
		return res, nil
	case <-ctx.Done():
		return 0, context.Cause(ctx)
	}
}

// work does something for 100 ms.
func work() int {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("work done")
	return 42
}

// maybeCancel waits for 50 ms and cancels with 50% probability.
func maybeCancel(cancel context.CancelCauseFunc, agent string) {
	time.Sleep(50 * time.Millisecond)
	if rand.Float32() < 0.5 {
		fmt.Println("canceling")
		cancel(fmt.Errorf("cancelled by %s", agent))
	}
}

func CancelRandomly() {
	parentCtx, parentCancel := context.WithCancelCause(context.Background())
	childCtx, childCancel := context.WithCancelCause(parentCtx)
	go maybeCancel(childCancel, "child")
	_ = childCtx // to avoid unused variable error
	_ = parentCancel

	// noCancelCtx := context.WithoutCancel(childCtx)
	res, err := execute(parentCtx, work)
	fmt.Println(res, err)
}
