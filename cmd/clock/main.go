package main

import (
    "time"

    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/events"
)

func main() {
    // Start of virtual time
    start := time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC)

    // Create test clock
    clk := clock.NewTestClock(start)

    // Create event queue + engine
    queue := engine.NewEventQueue()
    eng := engine.NewEngine(queue)

    // Execution context
    ctx := &context.Context{
        Clock: clk,
    }

    // Add events (out of order on purpose)
    queue.PushEvent(events.NewLogEvent(start.Add(2*time.Hour), "Event-B"))
    queue.PushEvent(events.NewLogEvent(start.Add(1*time.Hour), "Event-A"))

    // Advance virtual time
    eng.Advance(start.Add(3*time.Hour), ctx)
}
