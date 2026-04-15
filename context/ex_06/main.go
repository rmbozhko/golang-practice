package main

import (
	"context"
	"errors"
	"fmt"
	"practice/context/ex_06/err_type"
	"time"
)

// work does something for 100 ms.
func work(ctx context.Context) {
	select {
	case <-ctx.Done():
	}
}

func main() {
	rootCtx := context.Background()

	// Cancellation example
	cancelCtx, cancelFunc := context.WithCancelCause(rootCtx)

	context.AfterFunc(cancelCtx, func() {
		fmt.Println("Context cancelled #1")
	})

	context.AfterFunc(cancelCtx, func() {
		fmt.Println("Context cancelled #2")
	})

	context.AfterFunc(cancelCtx, func() {
		fmt.Println("Context cancelled #3")
	})

	stopFunc := context.AfterFunc(cancelCtx, func() {
		fmt.Println("Context cancelled #4")
	})
	stopFunc() // detach the function

	cancelFunc(err_type.CustomError)

	fmt.Println(context.Cause(cancelCtx))

	fmt.Println(errors.Is(context.Cause(cancelCtx), err_type.CustomError))

	// Timeout example
	parentCtx, cancel := context.WithTimeoutCause(rootCtx, 500*time.Millisecond, errors.New("parent timeout"))
	defer cancel()

	childCtx, childCancel := context.WithTimeoutCause(parentCtx, 100*time.Millisecond, errors.New("child timeout"))
	defer childCancel()

	_ = childCtx
	time.Sleep(200 * time.Millisecond)

	fmt.Println(childCtx.Err())
	fmt.Println(parentCtx.Err())

	fmt.Println(context.Cause(childCtx))
	fmt.Println(context.Cause(parentCtx))


	workCtx, cancel := context.WithTimeoutCause(rootCtx, 100 * time.Millisecond, errors.New("workCtx deadline exceeded"))
	// go work(workCtx) // also try to run in main goroutine
	// defer cancel()
	fmt.Println(workCtx.Err())
	fmt.Println(context.Cause(workCtx))
}
