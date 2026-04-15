package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

// execute runs fn in a separate goroutine
// and waits for the result unless canceled.
func execute(cancel <-chan struct{}, fn func() int) (int, error) {
	ch := make(chan int, 1)

	go func() {
		ch <- fn()
	}()

	select {
	case res := <-ch:
		return res, nil
	case <-cancel:
		return 0, errors.New("canceled")
	}
}

// work does something for 100 ms.
func work() int {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("work done")
	return 42
}

// maybeCancel waits for 50 ms and cancels with 50% probability.
func maybeCancel(cancel chan struct{}) {
	time.Sleep(50 * time.Millisecond)
	if rand.Float32() < 0.5 {
		close(cancel)
	}
}

func CancelRandomly() {
	cancel := make(chan struct{})
	go maybeCancel(cancel)

	res, err := execute(cancel, work)
	fmt.Println(res, err)
}
