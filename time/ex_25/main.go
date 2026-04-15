package main

import (
	"fmt"
	"runtime"
	"time"
)

type token struct{}

func consumer(cancel <-chan token, in <-chan token) {
	const timeout = time.Hour
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		timer.Reset(timeout)
		select {
		case <-in:
			// do stuff
		case <-timer.C:
			// log warning
		case <-cancel:
			fmt.Println("consumer exiting by cancel signal")
			return
		}
	}
}

// measure returns the number of bytes allocated
// and the number of allocations performed by the function fn.
func measure(fn func()) {
	var m runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m)
	allocBefore, mallocsBefore := m.TotalAlloc, m.Mallocs

	fn()

	runtime.GC()
	runtime.ReadMemStats(&m)
	allocAfter, mallocsAfter := m.TotalAlloc, m.Mallocs

	alloc := allocAfter - allocBefore
	mallocs := mallocsAfter - mallocsBefore
	fmt.Printf("Memory used: %d KB, # allocations: %d\n", alloc/1024, mallocs)
}

func main() {
	cancel := make(chan token)
	defer close(cancel)

	tokens := make(chan token)
	go consumer(cancel, tokens)

	measure(func() {
		for range 100000 {
			tokens <- token{}
		}
	})
}
