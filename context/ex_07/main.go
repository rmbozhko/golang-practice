package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Worker struct {
	// можеш додати поля (ресурси, логгер, etc.)
}

func (w *Worker) Do(ctx context.Context) error {
	d := 50 * time.Millisecond
	timer := time.NewTimer(d)
	for range 4 {
		timer.Reset(d)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// виконуємо роботу
		}
	}
	return nil
}

func main() {
	cleanup := func() {
		fmt.Println("cleanup resources")
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("deferring the main"))

	stopCleanup := context.AfterFunc(ctx, cleanup)
	w := &Worker{}
	go func() {
		time.Sleep(100 * time.Millisecond)
		// cancel(errors.New("manual stop"))
		// return
		// if rand.IntN(2) == 0 {
		// 	fmt.Println("canceling context")
		// 	cancel(errors.New("rand: context canceled by worker"))
		// }
	}()
	err := w.Do(ctx)
	if err != nil {
		fmt.Println("error occurred:", err)
		fmt.Println("error cause: ", context.Cause(ctx))
		fmt.Println("could stop cleaning up?", stopCleanup())
	} else {
		fmt.Println("work completed successfully")
		fmt.Println("could stop cleaning up?", stopCleanup())
	}
	time.Sleep(100 * time.Millisecond)

	timeoutCtx, timeoutCancel := context.WithTimeoutCause(
		context.Background(),
		100*time.Millisecond,
		errors.New("timeout exceeded"),
	)
	defer timeoutCancel(errors.New("normal canceling"))
	

}
