package main

import (
	"context"
	"errors"
	"fmt"
)

type Result struct {
	val int
	err error
}

// Generator – produces numbers 1..N.
func generator(ctx context.Context, upperBound int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range upperBound {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Worker – squares each number.
func worker(ctx context.Context, in <-chan int) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for n := range in {
			if n == 7 {
				out <- Result{err: errors.New("unlucky number")}
				return
			}
			select {
			case out <- Result{val: n * n}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func sink(in <-chan Result) error {
	for r := range in {
		if r.err != nil {
			return r.err
		} else {
			fmt.Println("Value:", r.val)
		}
	}
	return nil
}

func runPipeline(ctx context.Context, n int) error {
	in := generator(ctx, n)
	out := worker(ctx, in)
	return sink(out)
}

func Pipeline() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := runPipeline(ctx, 10); err != nil {
		cancel()
		fmt.Println("error:", err)
	}
}

func main() {
	Pipeline()
}
