package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context) <-chan int {
	out := make(chan int)

	go func() {
		defer fmt.Println("doWork anon goroutine is done")
		timeoutDuration := 2 * time.Second
		select {
		case <-time.After(timeoutDuration):
			select {
			case out <- 42:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}()

	return out
}

func main() {
	timeoutDuration := time.Second
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer timeoutCancel()

	select {
	case v := <-doWork(timeoutCtx):
		fmt.Println(v)
	case <-time.After(timeoutDuration):
		fmt.Println("timeout")
		timeoutCancel()
	}
}

package main

import (
	"fmt"
	"time"
)

func doWork() <-chan int {
	out := make(chan int)

	go func() {
		defer fmt.Println("doWork anon goroutine is done")
		time.Sleep(2 * time.Second)
		out <- 42
	}()

	return out
}

func main() {
	select {
	case v := <-doWork():
		fmt.Println(v)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
	time.Sleep(time.Second)
}
