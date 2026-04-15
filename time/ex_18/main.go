package main

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

/*
The queue has a fixed capacity.

Producers can enqueue tasks into the queue.

A worker dequeues and processes tasks sequentially.

If the queue is full, the producer must block until space becomes available.
*/

type Task struct {
	ID int
}

type Queue struct {
	tasks chan Task
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		tasks: make(chan Task, capacity),
	}
}

func (q *Queue) EnqueueTimeout(task Task, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()
	
	select {
	case q.tasks <- task:
		return nil
	case <-t.C:
		return errors.New("enqueue timeout")
	}
}

func (q *Queue) TryEnqueue(ctx context.Context, task Task) bool {
	select {
	case <-ctx.Done():
		return false
	case q.tasks <- task:
		return true
	default:
		return false
	}
}

func (q *Queue) Enqueue(ctx context.Context, task Task) {
	select {
	case <-ctx.Done():
		return
	case q.tasks <- task:
	}
}

func (q *Queue) Dequeue(ctx context.Context) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	case task, ok := <-q.tasks:
		if !ok {
			return Task{}, errors.New("queue closed")
		}
		return task, nil
	}
}

func (q *Queue) Close(cancel context.CancelFunc) {
	cancel()
}

func RunQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue := NewQueue(2)

	var processedTasks atomic.Int32
	var droppedTasks atomic.Int32

	var wg sync.WaitGroup
	// Start a worker to process tasks.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			task, err := queue.Dequeue(ctx)
			if err != nil {
				log.Println("worker stopped:", err)
				break
			}
			// Process the task (simulate work).
			log.Printf("process task=%d\n", task.ID)
			processedTasks.Add(1)
			time.Sleep(50 * time.Millisecond)
			if rand.IntN(10) == 2 {
				cancel() // Simulate worker stopping randomly.
				log.Println("simulating worker stop")
				break
			}
		}
	}()

	// Enqueue tasks.
	for i := 1; i <= 6; i++ {
		res := queue.TryEnqueue(ctx, Task{ID: i})
		log.Printf("enqueue task=%d accepted=%v\n", i, res)
		if !res {
			droppedTasks.Add(1)
		}
		time.Sleep(10 * time.Millisecond)
	}

	queue.Close(cancel)

	res := queue.TryEnqueue(ctx, Task{ID: 42})
	log.Printf("enqueue task=%d accepted=%v\n", 42, res)

	wg.Wait()

	log.Printf("processed tasks=%d dropped tasks=%d\n", processedTasks.Load(), droppedTasks.Load())
}

func main() {
	RunQueue()
}
