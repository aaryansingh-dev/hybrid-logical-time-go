package engine

import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
import "time"

type Engine struct {
    queue *EventQueue
}


func NewEngine(q *EventQueue) *Engine{
	return &Engine{queue: q}
}

func (engine *Engine) Advance(to time.Time, ctx *context.Context) {
    for {
        next := engine.queue.Peek()

        // No more events: jump clock to target time and stop
        if next == nil {
            ctx.Clock.Set(to)
            return
        }

        // Next event is after target time: stop advancing
        if next.Time().After(to) {
            ctx.Clock.Set(to)
            return
        }

        // Advance clock to event time
        ctx.Clock.Set(next.Time())

        // Remove event from queue
        event := engine.queue.PopEvent()

        // Execute event and schedule any new events
        futureEvents := event.Execute(ctx)
        for _, newEvent := range futureEvents {
            engine.queue.PushEvent(newEvent)
        }
    }
}

