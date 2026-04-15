package main

import (
	"fmt"
	"time"
)

func work() {
	fmt.Println("work done")
}

// func main() {
// 	var eventTime time.Time

// 	start := time.Now()
// 	timer := time.NewTimer(100 * time.Millisecond) // (1)
// 	go func() {
// 		eventTime = <-timer.C // (2)
// 		work()
// 	}()

// 	// enough time for the timer to expire
// 	time.Sleep(150 * time.Millisecond)
// 	if timer.Stop() {
// 		fmt.Printf("delayed function canceled after %v\n", time.Since(start))
// 	} else {
// 		fmt.Printf("delayed function started after %v\n", eventTime.Sub(start))
// 	}
// }

// func main() {
// 	t := time.AfterFunc(100*time.Millisecond, work)

// 	// enough time for the timer to expire
// 	time.Sleep(10 * time.Millisecond)
// 	fmt.Println(t.Reset(50 * time.Millisecond))
// 	if t.Stop() {
// 		fmt.Println("delayed function canceled")
// 	} else {
// 		fmt.Println("delayed function started")
// 	}
// }

func customAfterFunc(timeout time.Duration, f func()) {
	go func() {
		time.Sleep(timeout)
		f()
	}()
}

func main() {
	time.AfterFunc(100 * time.Millisecond, work)
	customAfterFunc(100 * time.Millisecond, work)

	time.Sleep(120 * time.Millisecond)
}
