package main

import (
	"fmt"
	"time"
)

func work(at time.Time) {
	fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
}

func main() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range 5 {
		at := <-ticker.C
		work(at)
	}

	// enough for 5 ticks
	time.Sleep(300 * time.Millisecond)
}
