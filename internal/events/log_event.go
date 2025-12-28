package events

import (
	"fmt"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type LogEvent struct {
	at      time.Time
	name    string
	clockID string
}

func NewLogEvent(at time.Time, name string, clockID string) *LogEvent {
	return &LogEvent{at: at, name: name, clockID: clockID}
}

func (e *LogEvent) Time() time.Time {
	return e.at
}

func (e *LogEvent) Name() string {
	return GetTypeName(e)
}

func (e *LogEvent) ClockID() string {
	return e.clockID
}

func (e *LogEvent) Execute(timeProvider clock.TimeProvider) []engine.Event {
	fmt.Printf("[%s] Executed %s\n", timeProvider.Now().Format(time.RFC3339), e.name)

	// Example: schedule a follow-up event
	if e.name == "Event-B" {
		newEventTime := timeProvider.Now().Add(90 * time.Minute)
		return []engine.Event{
			NewLogEvent(newEventTime, "Event-2A", e.clockID),
		}
	}

	return nil
}
