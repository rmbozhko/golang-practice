package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func work() int {
	if rand.Intn(10) < 8 {
		time.Sleep(10 * time.Millisecond)
	} else {
		time.Sleep(200 * time.Millisecond)
	}
	return 42
}

func withTimeout(timeout time.Duration, f func() int) (int, error) {
	startTimerCh := make(chan struct{})
	timeoutCh := make(chan struct{})
	
	var wg sync.WaitGroup
	// defer wg.Wait()

	wg.Go(func() {
		defer close(timeoutCh)
		<-startTimerCh
		time.Sleep(timeout)
	})

	close(startTimerCh)

	done := make(chan int, 1)
	wg.Go(func() {
		defer close(done)
		done <- f()
	})

	select {
	case result := <-done:
		return result, nil
	case <-timeoutCh:
		return 0, fmt.Errorf("timeout")
	}
}

func main() {
	for range 10 {
		start := time.Now()
		timeout := 50 * time.Millisecond
		if answer, err := withTimeout(timeout, work); err != nil {
			fmt.Printf("Took longer than %v. Error: %v\n", time.Since(start), err)
		} else {
			fmt.Printf("Took %v. Result: %v\n", time.Since(start), answer)
		}
	}
}
