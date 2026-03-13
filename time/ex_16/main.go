package main

import (
	"errors"
	"fmt"
	"time"
)

func throttle(n int, fn func()) (handle func() error, wait func()) {
	// Semaphore for n goroutines.
	sema := make(chan struct{}, n)

	// Execute fn functions concurrently, but not more than n at a time.
	handle = func() error {
		select {
		case sema <- struct{}{}:
			go func() {
				fn()
				<-sema
			}()
		default:
			return errors.New("busy")
		}
		return nil
	}

	// Wait until all functions have finished.
	wait = func() {
		for range n {
			sema <- struct{}{}
		}
	}

	return handle, wait
}

func work() {
	// Something very important, but not very fast.
	time.Sleep(100 * time.Millisecond)
}

func main() {
	handle, wait := throttle(2, work)
	start := time.Now()

	err := handle()
	fmt.Println("handle 1:", err)
	
	err = handle()
	fmt.Println("handle 2:", err)

	err = handle()
	fmt.Println("handle 3:", err)

	err = handle()
	fmt.Println("handle 4:", err)

	wait()

	fmt.Println("4 calls took", time.Since(start))
}
