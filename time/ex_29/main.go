package main

import (
	"fmt"
	"time"
)

/*
avoid timer leaks
behave correctly if:
the in channel is very active
the timeout rarely fires


Rewrite the function so that:

It does NOT allocate a new timer on every iteration
It behaves the same (timeout after 1 second of inactivity)
*/

func process(t time.Time) {
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("%s: work done\n", t.Format("15:04:05.000"))
}

func main() {
	time.AfterFunc(10*time.Millisecond, func() {
		fmt.Println("fired after func")
	})

	time.Sleep(20 * time.Millisecond)
	// t.Reset(10 * time.Millisecond)

	timer := time.NewTimer(10 * time.Millisecond)
	go func() {
		<-timer.C
		fmt.Println("fired timer")
	}()

	time.Sleep(20 * time.Millisecond)

	ticker := time.NewTicker(100 * time.Millisecond)

	for t := range ticker.C {
		process(t) // takes 300ms
	}
}
