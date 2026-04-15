package main_test

import (
	"context"
	"fmt"
	main "practice/context/ex_05"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	printLeaks(func() {
		timeout := 50 * time.Second
		main.Pipeline(timeout)
	})
}

func TestWithLongTimeout(t *testing.T) {
	expectedLength := 100
	timeout := 50 * time.Second
	result := main.Pipeline(timeout)
	assert.Nil(t, result.Err)
	assert.Equal(t, len(result.Numbers), expectedLength)
}

func TestWithShortTimeout(t *testing.T) {
	expectedLength := 100
	timeout := 5 * time.Millisecond

	result := main.Pipeline(timeout)

	assert.Equal(t, result.Err, context.DeadlineExceeded)
	assert.NotEqual(t, len(result.Numbers), expectedLength)
	assert.Less(t, len(result.Numbers), expectedLength)
}

func printLeaks(f func()) {
	prof := pprof.Lookup("goroutineleak")

	defer func() {
		time.Sleep(5000 * time.Millisecond)
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
