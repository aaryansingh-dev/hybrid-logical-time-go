package events

import (
    "fmt"
    "time"

    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
    ctxpkg "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
)

type LogEvent struct {
    at   time.Time
    name string
}

func NewLogEvent(at time.Time, name string) *LogEvent {
    return &LogEvent{at: at, name: name}
}

func (e *LogEvent) Time() time.Time {
    return e.at
}

func (e *LogEvent) Name() string {
    return e.name
}

func (e *LogEvent) Execute(ctx *ctxpkg.Context) []engine.Event {
    fmt.Printf("[%s] Executed %s\n", ctx.Clock.Now().Format(time.RFC3339), e.name)
    return nil
}
