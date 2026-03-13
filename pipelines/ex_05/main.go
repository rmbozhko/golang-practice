package main

import (
	"context"
	"fmt"
	"time"
)

func rangeGen(ctx context.Context, start, stop int) <-chan int {
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


// select works the best only when you know the amount of channels in advance otherwise you cannot predict the amount of cases needed if the function accepts channels as a variadic parameter. 
func merge(ctx context.Context, ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for ch1 != nil || ch2 != nil {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch1:
				if !ok {
					ch1 = nil
				} else {
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				}
			case val, ok := <-ch2:
				if !ok {
					ch2 = nil
				} else {
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				}
			}
		}
	}()
	return out
}

func Pipeline() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in2 := rangeGen(ctx, 21, 25)
	in1 := rangeGen(ctx, 11, 15)

	start := time.Now()
	merged := merge(ctx, in1, in2)
	for val := range merged {
		fmt.Print(val, " ")
		// if val == 23 {
		// 	cancel()
		// 	break
		// }
	}
	fmt.Println()
	fmt.Println("Took", time.Since(start))
}

func main() {
	Pipeline()
}
