package main

import (
	"fmt"
)

// Here's a function that sends numbers within a specified range to a channel:
func rangeGen(cancel <-chan struct{}, start, stop int) <-chan int {
	if start > stop {
		return nil
	}
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= stop; i++ {
			select {
			case out <- i:
			case <-cancel:
				fmt.Println("Received cancel signal, stopping generation")
				return
			}
		}
	}()
	return out
}

func printRange() {
	cancel := make(chan struct{})
	defer close(cancel)

	generated := rangeGen(cancel, 41, 46)
	for val := range generated {
		fmt.Println(val)
		// if val == 43 {
		// 	break
		// }
	}
}

func main() {
	printRange()

	// ch := make(chan int)
	// go func() {
	// 	for i := 0; i < 5; i++ {
	// 		ch <- i
	// 	}
	// 	close(ch)
	// }()
	// time.Sleep(5 * time.Second)
	// for val := range ch {
	// 	fmt.Println(val)
	// }
}
