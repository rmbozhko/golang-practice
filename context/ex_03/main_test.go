package main_test

import (
	"fmt"
	main "practice/context/ex_03"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

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

func TestPrintLeaks(t *testing.T) {
	printLeaks(main.GeneratorCancel)
}

func printLeaks(f func()) {
	prof := pprof.Lookup("goroutineleak")

	defer func() {
		time.Sleep(50 * time.Millisecond)
		var content strings.Builder
		prof.WriteTo(&content, 2)
		// Print only the leaked goroutines.
		goros := strings.Split(content.String(), "\n\n")
		for _, goro := range goros {
			if strings.Contains(goro, "(leaked)") {
				fmt.Println(goro + "\n")
			}
		}
	}()

	measure(f)
}
