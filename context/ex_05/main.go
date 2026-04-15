package main

import (
	"context"
	"fmt"
	"time"
)

func generator(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- i:
			}
		}
	}()
	return out
}

func worker(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if ok {
					time.Sleep(1 * time.Millisecond)
					select {
					case <-ctx.Done():
						return
					case out <- i * 2:
					}
				} else {
					return
				}
			}
		}
	}()
	return out
}

func collector(ctx context.Context, in <-chan int) ([]int, error) {
	res := make([]int, 0)
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case i, ok := <-in:
			if ok {
				res = append(res, i)
			} else {
				return res, nil
			}
		}
	}
}

type Result struct {
	Numbers []int
	Err     error
}

func Pipeline(timeout time.Duration) *Result {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	n := 100
	gen := generator(ctx, n)
	work := worker(ctx, gen)
	res, err := collector(ctx, work)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	return &Result{
		Numbers: res,
		Err:     err,
	}
}
