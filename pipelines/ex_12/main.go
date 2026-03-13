package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// generate produces numbers from 1 to stop inclusive.
func generate(stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range stop {
			out <- i + 1
		}
	}()
	return out
}

// Answer represents the result of a calculation.
type Answer struct {
	x, y int
}

// Retrieve an answer for a given number from a remote API.
func fetchAnswer(n int) (Answer, error) {
	// Simulate some work.
	time.Sleep(100 * time.Millisecond)
	// resp, err := http.Get(fmt.Sprintf("https://example.com/api/answer/%d", n))
	// if err != nil {
	// 	return Answer{}, err
	// }
	// defer resp.Body.Close()
	if rand.IntN(2) == 1 {
		return Answer{}, fmt.Errorf("fetchAnswer failed to retrieve")
	}

	return Answer{x: n, y: n * n}, nil
}

// calculate produces answers for the given numbers.
func calculate(in <-chan int, errc chan<- error) <-chan Answer {
	out := make(chan Answer)
	go func() {
		defer close(out)
		for n := range in {
			ans, err := fetchAnswer(n)
			if err != nil {
				errc <- err
			} else {
				out <- ans
			}
		}
	}()
	return out
}

func errCollector(in <-chan error) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for err := range in {
			fmt.Println("error:", err)
		}
	}()
	return done
}

func main() {
	errc := make(chan error)
	errcDone := errCollector(errc)
	inputs := generate(5)
	results := calculate(inputs, errc)

	for res := range results {
		fmt.Printf("%v\n", res)
	}
	close(errc)
	<-errcDone
}
