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

// Result contains a value or an error.
type Result[T any] struct {
	val T
	err error
}

func (r Result[T]) OK() bool {
	return r.err == nil
}
func (r Result[T]) Val() T {
	return r.val
}
func (r Result[T]) Err() error {
	return r.err
}

// calculate produces answers for the given numbers.
func calculate(in <-chan int) <-chan Result[Answer] {
	out := make(chan Result[Answer])
	go func() {
		defer close(out)
		for n := range in {
			ans, err := fetchAnswer(n)
			out <- Result[Answer]{val: ans, err: err}
		}
	}()
	return out
}

// func main() {
// 	inputs := generate(5)
// 	answers, errs := calculate(inputs)

// 	for ans := range answers {
// 		fmt.Printf("%d -> %d\n", ans.x, ans.y)
// 	}
// 	if err := <-errs; err != nil {
// 		fmt.Println("error:", err)
// 	}
// }

func main() {
	inputs := generate(5)
	results := calculate(inputs)

	for res := range results {
		if res.OK() {
			fmt.Println(res.Val())
		} else {
			fmt.Println("error:", res.Err())
		}
	}
}
