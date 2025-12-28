package engine

import (
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
)

type Event interface {
	Time() time.Time // time stamp at which the event will execute
	Name() string    // Name of the event for logging purposes
	ClockID() string
	Execute(timeProvider clock.TimeProvider) []Event // executes the event -> called by the engine running the logic
}
