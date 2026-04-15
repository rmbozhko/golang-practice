package main

import (
	"context"
	"errors"
	"fmt"
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
		return 0, ctx.Err()
	}
}

// work does something for 100 ms.
func work() int {
	time.Sleep(100 * time.Millisecond)
	return 42
}

// slow does something for 300 ms.
func slow() int {
	time.Sleep(300 * time.Millisecond)
	return 13
}

func main() {

	ctx := context.Background()
	fmt.Println(ctx.Deadline())

	const dur200ms = 200 * time.Millisecond
	parentCtx, cancel := context.WithTimeoutCause(context.Background(), dur200ms, errors.New("parent context deadline exceeded"))
	defer cancel()

	const dur50ms = 50 * time.Millisecond
	childCtx, cancel := context.WithTimeoutCause(parentCtx, dur50ms, errors.New("child context deadline exceeded"))
	defer cancel()


	// completes in time
	res, err := execute(childCtx, slow)
	fmt.Printf("%d %v (%v)", res, err, context.Cause(childCtx))
}
