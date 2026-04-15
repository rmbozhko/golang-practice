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

/*
"user:123"
"order:998"
"cache:item:abc"
*/
type Event struct {
	Key string
}

type Queue struct {
	queue chan Event
}

type ParameterizedQueue struct {
	Queue
	registeredEventsMu sync.Mutex
	registeredEvents   map[string]bool
	throttleDuration   time.Duration
	throttlingEventsMu sync.Mutex
	throttlingEvents   map[string]*time.Timer
}

func NewQueue(throttle time.Duration, capacity int) *ParameterizedQueue {
	q := &ParameterizedQueue{
		throttleDuration:   throttle,
		registeredEventsMu: sync.Mutex{},
		registeredEvents:   make(map[string]bool),
		throttlingEventsMu: sync.Mutex{},
		throttlingEvents:   make(map[string]*time.Timer),
		Queue: Queue{
			queue: make(chan Event, capacity),
		},
	}
	return q
}

// true  -> event accepted
// false -> event dropped
func (q *ParameterizedQueue) TryEnqueue(ctx context.Context, event Event) bool {
	q.registeredEventsMu.Lock()
	if q.registeredEvents[event.Key] {
		q.registeredEventsMu.Unlock()
		deduplicatedEvents.Add(1)
		return false
	}
	q.registeredEventsMu.Unlock()

	select {
	case <-ctx.Done():
		return false
	default:
		q.throttlingEventsMu.Lock()
		if timer, exists := q.throttlingEvents[event.Key]; exists {
			log.Printf("throttling event key=%s...", event.Key)
			throttledEvents.Add(1)
			<-timer.C
			delete(q.throttlingEvents, event.Key)
		}
		q.throttlingEventsMu.Unlock()
		select {
		case <-ctx.Done():
			return false
		case q.queue <- event:
			q.registeredEventsMu.Lock()
			q.registeredEvents[event.Key] = true
			q.registeredEventsMu.Unlock()
			return true
		default:
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
		q.throttlingEvents[event.Key] = time.NewTimer(q.throttleDuration)
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
				log.Printf("enqueue id=%s, accepted=%t", event.Key, accepted)
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
				// if rand.IntN(10) == 5 {
				// 	log.Printf("consumer error, cancelling context")
				// 	cancel() // Simulate consumer cancellation
				// 	return
				// }
			}
		}()
	}

	producerWG.Wait()
	queue.Close()

	consumerWG.Wait()

	log.Printf("acceptedEvents=%d, droppedEvents=%d, deduplicatedEvents=%d, throttledEvents=%d, processedEvents=%d",
		acceptedEvents.Load(),
		droppedEvents.Load(),
		deduplicatedEvents.Load(),
		throttledEvents.Load(),
		processedEvents.Load(),
	)
}
