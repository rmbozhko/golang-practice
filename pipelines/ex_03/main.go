package main

import "fmt"

func generator(cancel <-chan struct{}) <-chan int {
	i := 0
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				return
			case out <- i:
				i++
			}
		}
	}()
	return out
}

func main() {
	cancel := make(chan struct{})

	out := generator(cancel)
	for v := range out {
		if v == 10 {
			close(cancel)
			break
		}
		fmt.Println(v)
	}
}
