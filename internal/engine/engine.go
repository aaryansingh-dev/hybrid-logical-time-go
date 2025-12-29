package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
)

// Now moving to multi-tenancy which can handle multiple event queues and clocks.
// This will allow for simulating multiple independent systems in the same engine. Real Time Prod events and
// Virtual Time Test events can co-exist without interfering with each other.

type Engine struct {
	queues      map[string]*EventQueue
	clocks      map[string]clock.TimeProvider
	systemQueue *EventQueue

	mu   sync.RWMutex
	diag Diagnostic
}

func NewEngine(diag Diagnostic) *Engine {
	// Initialize the engine with empty maps for queues and clocks
	return &Engine{
		queues:      make(map[string]*EventQueue),
		clocks:      make(map[string]clock.TimeProvider),
		diag:        diag,
		systemQueue: NewEventQueue(),
	}
}

func (engine *Engine) RegisterPartition(partitionID string, timeProvider clock.TimeProvider) {
	// Useful for reigstering a new clock when a new simulation is started by the user.
	engine.mu.Lock()
	defer engine.mu.Unlock()

	engine.clocks[partitionID] = timeProvider
	if _, exists := engine.queues[partitionID]; !exists {
		engine.queues[partitionID] = NewEventQueue()
	}

}

func (engine *Engine) getPartition(id string) (*EventQueue, clock.TimeProvider, error) {

	// design decision: returns a clock.TimeProvider instead of a TestClock to better handle more different
	// clocks in the future: Liskov Substitution Principle, and Open Closed Principle

	engine.mu.RLock()
	defer engine.mu.RUnlock()

	queue, qOk := engine.queues[id]
	clock, cOk := engine.clocks[id]

	if !qOk || !cOk {
		return nil, nil, fmt.Errorf("clock id %s not registered", id)
	}

	return queue, clock, nil
}

func (engine *Engine) Schedule(event Event) {
	partitionID := event.ClockID()

	if partitionID == "SYSTEM" {
		engine.systemQueue.PushEvent(event)
		return
	}

	// If it doesn't exist, we auto-register.
	engine.mu.RLock()
	queue, exists := engine.queues[partitionID]
	engine.mu.RUnlock()

	// lazy registry pattern in case the queue was not made
	if !exists {
		engine.mu.Lock()
		// double-check pattern to handle concurrent initialization racing.
		if q, ok := engine.queues[partitionID]; ok {
			queue = q
		} else {
			queue = NewEventQueue()
			engine.queues[partitionID] = queue
		}
		engine.mu.Unlock()
	}

	queue.PushEvent(event)
}

func (engine *Engine) StartRealTimeWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	realTime := clock.NewRealTimeProvider()

	go func() {
		for range ticker.C {
			now := time.Now().UTC()
			for {
				next := engine.systemQueue.Peek()

				if next == nil || next.Time().After(now) {
					break
				}

				event := engine.systemQueue.PopEvent()
				if engine.diag != nil {
					engine.diag.OnEventExecute("SYSTEM", event.Name(), now)
				}

				futureEvents := event.Execute(realTime)
				for _, futureEvent := range futureEvents {
					engine.Schedule(futureEvent)
				}

			}
		}
	}()

}

// this is specifically used by the Test users to test their timelines
func (engine *Engine) Advance(partitionID string, to time.Time, ctx *context.Context) error {

	if partitionID == "SYSTEM" {
		return fmt.Errorf("invalid operation: the SYSTEM partition follows wall-clock time and cannot be advanced manually")
	}

	// Resolve the specific queue and clock for this tenant
	queue, provider, err := engine.getPartition(partitionID)
	if err != nil {
		return err
	}

	if engine.diag != nil {
		engine.diag.OnAdvanceStart(partitionID, provider.Now(), to)
	}

	testClock, ok := provider.(*clock.TestClock)
	if !ok {
		return fmt.Errorf("Partition %s is not a TestClock; manual time warping is only supported for simulation partitions", partitionID)
	}

	for {
		next := queue.Peek()

		// EXIT CONDITION: If no more events exist OR the next event is
		// scheduled for a time after our target, we jump to target and stop.
		if next == nil || next.Time().After(to) {
			testClock.Set(to)
			if engine.diag != nil {
				engine.diag.OnAdvanceFinish(partitionID, testClock.Now())
			}
			return nil
		}

		// teleport to the next event
		testClock.Set(next.Time())
		event := queue.PopEvent()

		if engine.diag != nil {
			engine.diag.OnEventExecute(partitionID, event.Name(), testClock.Now())
		}

		// Execute logic and handle "Causality" (chained events)
		futureEvents := event.Execute(testClock)
		for _, futureEvent := range futureEvents {
			engine.Schedule(futureEvent)
			if engine.diag != nil {
				engine.diag.OnEventCreated(partitionID, futureEvent.Name(), testClock.Now())
			}
		}
	}
}
