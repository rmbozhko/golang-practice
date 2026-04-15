package main_test

import (
	"fmt"
	main "practice/context/ex_02"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

func TestPrintLeaks(t *testing.T) {
	printLeaks(main.CancelRandomly)
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

	f()
}
