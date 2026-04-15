package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func work() {
	fmt.Println("work done")
}

func main() {
	start := time.Now()
	timer := time.NewTimer(100 * time.Millisecond)

	var timeAfterDelay time.Time
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(timeAfterDelay *time.Time) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Printf("function canceled after %v\n", time.Since(start))
			return
		case t, ok := <-timer.C:
			if !ok {
				fmt.Printf("timer stopped before firing after %v\n", time.Since(start))
				return
			}
			*timeAfterDelay = t
			work()
		}
	}(&timeAfterDelay)

	time.Sleep(10 * time.Millisecond)
	fmt.Println("main goroutine sleep has passed...")

	// the timer hasn't expired yet
	if timer.Stop() {
		fmt.Printf("delayed function canceled after %v\n", time.Since(start))
		cancel()
	} else {
		fmt.Printf("delayed function already executed after %v\n", timeAfterDelay.Sub(start))
	}
	wg.Wait()
	

	// Second option to run the delayed function
	afterFuncTimer := time.AfterFunc(100*time.Millisecond, work)
	time.Sleep(1000 * time.Millisecond)

	if afterFuncTimer.Stop() {
		fmt.Printf("option #2: delayed function canceled")
	} else {
		fmt.Printf("option #2: delayed function already executed")
	}
}
