package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func rangeGenWithWgGo(ctx context.Context, start, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < stop; i++ {
			time.Sleep(50 * time.Millisecond)
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func mergeWithWgGo(ctx context.Context, ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Go(func() {
		for val := range ch1 {
			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	})

	wg.Go(func() {
		for val := range ch2 {
			select {
			case out <- val:
			case <-ctx.Done():
				return
			}
		}
	})

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func PipelineWithWgGo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in2 := rangeGenWithWgGo(ctx, 21, 25)
	in1 := rangeGenWithWgGo(ctx, 11, 15)

	start := time.Now()
	merged := mergeWithWgGo(ctx, in1, in2)
	for val := range merged {
		fmt.Print(val, " ")
		if val == 23 {
			cancel()
			break
		}
	}
	fmt.Println()
	fmt.Println("Took", time.Since(start))
}

func main() {
	PipelineWithWgGo()
}
