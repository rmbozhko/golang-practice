package main

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, jobs <-chan int, results chan<- int) {
	defer close(results) // FIX: worker is producer that outputs to result channel, so it should close it when done
	for job := range jobs {
		select {
		case results <- job * 2:
		case <-ctx.Done(): // FIX: add case for context cancellation to prevent goroutine leak
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int)
	results := make(chan int)

	go worker(ctx, jobs, results)

	go func() {
		for i := 0; i < 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	fmt.Println(<-results)
	cancel()
}
