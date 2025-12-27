package engine

import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
import "time"

type Event interface {
	Time() time.Time              // time stamp at which the event will execute
	Name() string                 // Name of the event for logging purposes
	Execute(timeProvider clock.TimeProvider) []Event // executes the event -> called by the engine running the logic
}
