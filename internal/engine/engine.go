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
	queues map[string]*EventQueue
	clocks map[string]*clock.TestClock
	mu     sync.RWMutex
	diag   Diagnostic
}

func NewEngine() *Engine {
	// Initialize the engine with empty maps for queues and clocks
	return &Engine{queues: make(map[string]*EventQueue), clocks: make(map[string]*clock.TestClock)}
}

func (engine *Engine) RegisterClock(id string, testClock *clock.TestClock) {
	// Useful for reigstering a new clock when a new simulation is started by the user.
	engine.mu.Lock()
	defer engine.mu.Unlock()

	engine.clocks[id] = testClock
	if _, exists := engine.queues[id]; !exists {
		engine.queues[id] = NewEventQueue()
	}

}

func (engine *Engine) getPartition(id string) (*EventQueue, *clock.TestClock, error) {

	engine.mu.RLock()
	defer engine.mu.RUnlock()

	queue, qOk := engine.queues[id]
	clock, cOk := engine.clocks[id]

	if !qOk || !cOk {
		return nil, nil, fmt.Errorf("clock id %s not registered", id)
	}

	return queue, clock, nil
}

func (engine *Engine) Advance(id string, to time.Time, ctx *context.Context) error {
	// Resolve the specific queue and clock for this tenant
	queue, clock, err := engine.getPartition(id)
	if err != nil {
		return err
	}

	if engine.diag != nil {
		engine.diag.OnAdvanceStart(id, clock.Now(), to)
	}

	for {
		next := queue.Peek()

		// EXIT CONDITION: If no more events exist OR the next event is
		// scheduled for a time after our target, we jump to target and stop.
		if next == nil || next.Time().After(to) {
			clock.Set(to)
			if engine.diag != nil {
				engine.diag.OnAdvanceFinish(id, clock.Now())
			}
			return nil
		}

		// teleport to the next event
		clock.Set(next.Time())
		event := queue.PopEvent()

		if engine.diag != nil {
			engine.diag.OnEventExecute(id, event.Name(), clock.Now())
		}

		// Execute logic and handle "Causality" (chained events)
		futureEvents := event.Execute(clock)
		for _, fe := range futureEvents {
			queue.PushEvent(fe)
			if engine.diag != nil {
				engine.diag.OnEventCreated(id, fe.Name(), clock.Now())
			}
		}
	}
}
