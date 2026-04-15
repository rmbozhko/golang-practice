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

var (
	acceptedEvents     atomic.Int64
	droppedEvents      atomic.Int64
	deduplicatedEvents atomic.Int64
	throttledEvents    atomic.Int64
	processedEvents    atomic.Int64
)

type Event struct {
	Key string
}

type Queue struct {
	queue chan Event
}

type ParameterizedQueue struct {
	Queue
	registeredEventsMu sync.Mutex
	registeredEvents   map[string]struct{}
	throttleDuration   time.Duration
	throttlingEventsMu sync.Mutex
	throttlingEvents   map[string]time.Time
}

func NewQueue(throttle time.Duration, capacity int) *ParameterizedQueue {
	q := &ParameterizedQueue{
		throttleDuration:   throttle,
		registeredEventsMu: sync.Mutex{},
		registeredEvents:   make(map[string]struct{}),
		throttlingEventsMu: sync.Mutex{},
		throttlingEvents:   make(map[string]time.Time),
		Queue: Queue{
			queue: make(chan Event, capacity),
		},
	}
	return q
}

// true  -> event accepted
// false -> event dropped
func (q *ParameterizedQueue) TryEnqueue(ctx context.Context, event Event) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		q.throttlingEventsMu.Lock()
		if timer, exists := q.throttlingEvents[event.Key]; exists {
			if time.Since(timer) < q.throttleDuration {
				log.Printf("throttling event key=%s", event.Key)
				throttledEvents.Add(1)
				q.throttlingEventsMu.Unlock()
				return false
			}
			delete(q.throttlingEvents, event.Key)
		}
		q.throttlingEventsMu.Unlock()

		q.registeredEventsMu.Lock()
		if _, exists := q.registeredEvents[event.Key]; exists {
			deduplicatedEvents.Add(1)
			q.registeredEventsMu.Unlock()
			return false
		}
		q.registeredEvents[event.Key] = struct{}{}
		q.registeredEventsMu.Unlock()

		select {
		case q.queue <- event:
			return true
		default:
			q.registeredEventsMu.Lock()
			delete(q.registeredEvents, event.Key)
			q.registeredEventsMu.Unlock()
			return false
		}
	}
}

func (q *ParameterizedQueue) Dequeue(ctx context.Context) (Event, error) {
	select {
	case <-ctx.Done():
		return Event{}, ctx.Err()
	case event, ok := <-q.queue:
		if !ok {
			return Event{}, errors.New("queue closed")
		}

		q.registeredEventsMu.Lock()
		delete(q.registeredEvents, event.Key)
		q.registeredEventsMu.Unlock()

		q.throttlingEventsMu.Lock()
		q.throttlingEvents[event.Key] = time.Now()
		q.throttlingEventsMu.Unlock()

		processedEvents.Add(1)

		return event, nil
	}
}

func (q *ParameterizedQueue) Close() {
	close(q.queue)
}

func main() {
	queue := NewQueue(500*time.Millisecond, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerCount := 1
	producerCount := 5
	var consumerWG sync.WaitGroup
	var producerWG sync.WaitGroup

	producerWG.Add(producerCount)
	for range producerCount {
		go func() {
			defer producerWG.Done()
			events := []Event{
				{Key: "user:1"},
				{Key: "user:2"},
				{Key: "user:3"},
				{Key: "user:4"},
				{Key: "user:5"},
				{Key: "user:6"},
				{Key: "user:7"},
				{Key: "user:8"},
				{Key: "user:9"},
				{Key: "user:10"},
			}
			for i := 0; i < 5; i++ {
				event := events[rand.IntN(len(events))]
				accepted := queue.TryEnqueue(ctx, event)
				log.Printf("enqueue id=%s accepted=%t", event.Key, accepted)
				if accepted {
					acceptedEvents.Add(1)
				} else {
					droppedEvents.Add(1)
				}
				producerInterval := rand.IntN(31)
				time.Sleep(time.Duration(producerInterval) * time.Millisecond)
			}
		}()
	}

	consumerWG.Add(consumerCount)
	for range consumerCount {
		go func() {
			defer consumerWG.Done()
			for {
				event, err := queue.Dequeue(ctx)
				if err != nil {
					log.Print("error dequeuing event:", err.Error())
					return
				}
				log.Printf("dequeue id=%s", event.Key)
				time.Sleep(100 * time.Millisecond) // Simulate processing time
				if rand.IntN(10) == 5 {
					log.Printf("consumer error, cancelling context")
					cancel() // Simulate consumer cancellation
					return
				}
			}
		}()
	}

	producerWG.Wait()
	queue.Close()

	consumerWG.Wait()

	log.Printf("acceptedEvents=%d droppedEvents=%d deduplicatedEvents=%d throttledEvents=%d processedEvents=%d",
		acceptedEvents.Load(),
		droppedEvents.Load(),
		deduplicatedEvents.Load(),
		throttledEvents.Load(),
		processedEvents.Load(),
	)
}
