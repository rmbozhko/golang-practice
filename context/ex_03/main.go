package main

import (
	"context"
	"fmt"
	"time"
)

func generator(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		var i int
		for {
			i++
			select {
			case <-ctx.Done():
				return
			case out <- i:
			}
		}
	}()
	return out
}

func GeneratorCancel() {
	timeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ch := generator(ctx)
	
	for i := range ch {
		fmt.Println("channel val", i)
	}
}