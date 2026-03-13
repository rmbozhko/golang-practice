package main

import (
	"fmt"
)



func main() {
	var stream chan int

	select {
	case event := <-stream:
		fmt.Println("Received event:", event)
	default:
		fmt.Println("No events received")
	}
	fmt.Println("There")
	close(stream)
}