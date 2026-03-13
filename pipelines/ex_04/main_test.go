package main_test

import (
	"fmt"
	main "practice/ex_04"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

// func TestSyncTest(t *testing.T) {
// 	synctest.Test(t, func(t *testing.T) {
// 		main.RunGenerator()
// 		synctest.Wait()
// 	})
// }

// run with -v flag to see the output about leaked goroutines
func TestPrintLeaks(t *testing.T) {
	printLeaks(main.RunGenerator)
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
