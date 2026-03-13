package main

import (
	"context"
	"fmt"
	"reflect"
)

func generator[T int](ctx context.Context, start, stop T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := start; i < stop; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func mergeWithReflect[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		cases := make([]reflect.SelectCase, len(channels)+1)

		for i, ch := range channels {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}

		// cancellation case
		cases[len(channels)] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		}

		active := len(channels)

		for active > 0 {
			chosen, value, ok := reflect.Select(cases)

			// cancellation
			if chosen == len(channels) {
				return
			}

			if !ok {
				// channel closed → disable case
				cases[chosen].Chan = reflect.Value{}
				active--
				continue
			}

			out <- value.Interface().(T)
		}
	}()

	return out
}

func PipelineWithReflect() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := generator(ctx, 0, 3)
	b := generator(ctx, 10, 13)
	c := generator(ctx, 100, 103)

	out := mergeWithReflect(ctx, a, b, c)

	for v := range out {
		// if v == 2 {
		// 	cancel()
		// 	break
		// } else {
		// 	fmt.Println(v)
		// }
		fmt.Println(v)
	}
}

func main() {
	PipelineWithReflect()
}
