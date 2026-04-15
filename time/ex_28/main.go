package main

import (
	"fmt"
	"time"
)

func work(at time.Time) {
	fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	time.Sleep(100 * time.Millisecond)
}

func main() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for {
			at := <-ticker.C
			work(at)
		}
	}()

	// enough for 3 ticks because of the slow work()
	time.Sleep(360 * time.Millisecond)
}
