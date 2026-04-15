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
	ctx   context.Context
	tasks chan Task
}

func NewQueue(ctx context.Context, capacity int) *Queue {
	return &Queue{
		ctx:   ctx,
		tasks: make(chan Task, capacity),
	}
}

func (q *Queue) TryEnqueue(task Task) bool {
	select {
	case <-q.ctx.Done():
		return false
	case q.tasks <- task:
		return true
	default:
		return false
	}
}

func (q *Queue) Enqueue(task Task) {
	q.tasks <- task
}

func (q *Queue) Dequeue() (Task, error) {
	select {
	case <-q.ctx.Done():
		return Task{}, q.ctx.Err()
	case task, ok := <-q.tasks:
		if !ok {
			return Task{}, errors.New("queue closed")
		}
		return task, nil
	}
}

func (q *Queue) Close() {
	close(q.tasks)
}

func RunQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue := NewQueue(ctx, 2)

	var processed_tasks atomic.Int32
	var dropped_tasks atomic.Int32

	var wg sync.WaitGroup
	// Start a worker to process tasks.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			task, err := queue.Dequeue()
			if err != nil {
				log.Println("Worker stopped:", err)
				break
			}
			// Process the task (simulate work).
			log.Printf("Processing task %d\n", task.ID)
			processed_tasks.Add(1)
			time.Sleep(50 * time.Millisecond)
			if rand.IntN(10) == 2 {
				cancel() // Simulate worker stopping randomly.
				log.Println("Simulating worker stop")
				break
			}
		}
	}()

	// Enqueue tasks.
	for i := 1; i <= 6; i++ {
		res := queue.TryEnqueue(Task{ID: i})
		log.Printf("Enqueuing task %d, %v\n", i, res)
		if !res {
			dropped_tasks.Add(1)
		}
		time.Sleep(10 * time.Millisecond)
	}

	queue.Close()

	res := queue.TryEnqueue(Task{ID: 42})
	log.Printf("Enqueuing task %d, %v\n", 42, res)

	wg.Wait()

	log.Printf("Processed tasks: %d, Dropped tasks: %d\n", processed_tasks.Load(), dropped_tasks.Load())
}

func main() {
	RunQueue()
}
