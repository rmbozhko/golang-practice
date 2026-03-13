package main

import (
	"context"
	"fmt"
	"sync"
)

func generator[T int](ctx context.Context, start, stop T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := start; i < stop; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(channel <-chan T) {
			defer wg.Done()
			for val := range ch {
				select {
				case out <- val:
				case <-ctx.Done():
					return
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

func Pipeline() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := generator(ctx, 0, 3)
	b := generator(ctx, 10, 13)
	c := generator(ctx, 100, 103)

	out := merge(ctx, a, b, c)

	for v := range out {
		// if v == 2 {
		// 	cancel()
		// 	break
		// } else {
		// 	fmt.Println(v)
		// }
		fmt.Println(v)
	}
}

func main() {
	Pipeline()
}
