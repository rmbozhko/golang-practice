package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Resource struct {
	closed atomic.Bool
}

type Worker struct {
	// можеш додати поля (ресурси, логгер, etc.)
	resource *Resource
}

func (w *Worker) Do(ctx context.Context) error {
	d := 50 * time.Millisecond
	timer := time.NewTimer(d)
	for range 4 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// виконуємо роботу
		}
		timer.Reset(d)
	}
	return nil
}

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	w := &Worker{
		resource: &Resource{},
	}
	cleanup := func() {
		defer wg.Done()
		fmt.Println("cleanup resources")
		if w.resource.closed.CompareAndSwap(false, true) {
			fmt.Println("resources cleaned up")
		} else {
			fmt.Println("resources were already cleaned up")
		}
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("deferring the main"))

	stopCleanup := context.AfterFunc(ctx, cleanup)
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel(errors.New("manual stop"))
		return
		// if rand.IntN(2) == 0 {
		// 	fmt.Println("canceling context")
		// 	cancel(errors.New("rand: context canceled by worker"))
		// }
	}()
	err := w.Do(ctx)
	go cleanup()
	if err != nil {
		fmt.Println("error occurred:", err)
		fmt.Println("error cause: ", context.Cause(ctx))
		fmt.Println("could stop cleaning up?", stopCleanup())
	} else {
		fmt.Println("work completed successfully")
		fmt.Println("could stop cleaning up?", stopCleanup())
	}
	cleanup()
	wg.Wait()
	fmt.Printf("is worker resource closed: %t\n", w.resource.closed.Load())

	fmt.Println("--------------TIMEOUT PART------------------")
	timeoutCtx, timeoutCancel := context.WithTimeoutCause(
		context.Background(),
		100*time.Millisecond,
		errors.New("timeout exceeded"),
	)
	defer timeoutCancel()

	wg.Add(1)
	stop := context.AfterFunc(timeoutCtx, func() {
		defer wg.Done()
		fmt.Println("cleanup resources for timeout context")
		if w.resource.closed.CompareAndSwap(false, true) {
			fmt.Println("resources cleaned up")
		} else {
			fmt.Println("resources were already cleaned up")
		}
	})
	w.resource.closed.Store(false)
	err = w.Do(timeoutCtx)
	if err != nil {
		fmt.Println("error occurred:", err)
		fmt.Println("error cause: ", context.Cause(timeoutCtx))
		fmt.Println("could stop cleaning up?", stop())
	} else {
		fmt.Println("work completed successfully")
		fmt.Println("could stop cleaning up?", stop())
	}
	fmt.Printf("is worker resource closed: %t\n", w.resource.closed.Load())
	wg.Wait()


	fmt.Println("--------------CANCELABLE CLEANUP PART------------------")
	ctx = context.WithoutCancel(ctx)
	cleanupCtx, cleanupCancel := context.WithTimeoutCause(ctx, 10*time.Second, errors.New("cleanup timeout"))
	defer cleanupCancel()
	
	ctx, cancelFunc := context.WithCancelCause(context.Background())
	
	wg.Add(1)
	context.AfterFunc(ctx, func() {
		defer wg.Done()
		select {
		case <-time.After(5 * time.Minute):
		case <-cleanupCtx.Done():
			fmt.Println("cleanup canceled:", cleanupCtx.Err(), context.Cause(ctx))
		}
		fmt.Println("cleanup resources for cancelable cleanup context")
	})
	cancelFunc(errors.New("canceling cleanup"))
	
	fmt.Println("cleanup should have been canceled")

	wg.Wait()
}

/*
If cleanup can block indefinitely, AfterFunc can leak a goroutine. Since Go cannot forcibly kill goroutines, cleanup must be cooperative:

- accept its own context
- have its own timeout
- avoid blocking syscalls without deadlines
- avoid holding locks across cancelable waits
- be idempotent
*/

select {
case err := <-db.Close():
    return err

case <-time.After(5 * time.Second):
    return errors.New("close timed out")
}
