package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var start time.Time

	work := func() {
		fmt.Printf("work done after %dms\n", time.Since(start).Milliseconds())
	}

	// run work after 10 milliseconds
	timeout := 10 * time.Millisecond
	start = time.Now() // ignore the data race for simplicity
	t := time.AfterFunc(timeout, work)

	// wait for 5 to 15 milliseconds
	delay := time.Duration(5+rand.Intn(11)) * time.Millisecond
	time.Sleep(delay)
	fmt.Printf("%dms has passed...\n", delay.Milliseconds())

	// Reset behavior depends on whether the timer has expired
	t.Reset(timeout)
	start = time.Now()

	time.Sleep(50 * time.Millisecond)
}
