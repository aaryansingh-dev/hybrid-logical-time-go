package engine

import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
import "time"

type Engine struct {
    queue *EventQueue
    clock *clock.TestClock
}


func NewEngine(q *EventQueue, clock *clock.TestClock) *Engine{
	return &Engine{queue: q, clock: clock}
}

func (engine *Engine) Advance(to time.Time, ctx *context.Context) {
    for {
        next := engine.queue.Peek()

        // No more events: jump clock to target time and stop
        if next == nil {
            engine.clock.Set(to)
            return
        }

        // Next event is after target time: stop advancing
        if next.Time().After(to) {
            engine.clock.Set(to)
            return
        }

        // Advance clock to event time
        engine.clock.Set(next.Time())

        // Remove event from queue
        event := engine.queue.PopEvent()

        // Execute event and schedule any new events
        futureEvents := event.Execute(engine.clock)
        for _, newEvent := range futureEvents {
            engine.queue.PushEvent(newEvent)
        }
    }
}

