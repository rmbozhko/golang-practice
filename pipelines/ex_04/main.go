package main

import (
	"context"
	"fmt"
	"runtime/pprof"
	"strings"
	"time"
)

func generator(ctx context.Context) <-chan int {
	i := 0
	out := make(chan int, 10)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case out <- i:
				i++
			}
		}
	}()
	return out
}

func RunGenerator() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := generator(ctx)
	for v := range out {
		if v == 10 {
			break
		}
		fmt.Println(v)
	}
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

func main() {
	printLeaks(RunGenerator)
}
