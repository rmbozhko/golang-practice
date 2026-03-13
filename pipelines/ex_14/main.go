package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
				select {
				case out <- Result{err: errors.New("unlucky number")}:
				case <-ctx.Done():
				}
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

func merge[T any](ctx context.Context, chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	wg.Add(len(chs))
	for _, ch := range chs {
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func runPipeline(ctx context.Context, n int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	workerCount := 4
	workerChs := make([]<-chan Result, workerCount)
	
	in := generator(ctx, n)
	// fan-out
	for i := range workerCount {
		workerChs[i] = worker(ctx, in)
	}
	// fan-in
	out := merge(ctx, workerChs...)

	return sink(out)
}

func Pipeline() {
	if err := runPipeline(context.Background(), 10); err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	Pipeline()
}
