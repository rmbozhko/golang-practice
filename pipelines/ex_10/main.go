package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		for v := range ch {
			go func() {
				fmt.Println(v)
			}()
		}
	}()

	ch <- 1
	ch <- 2
	close(ch)

	time.Sleep(time.Second)
}
