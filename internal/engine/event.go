package engine

import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/context"
import "time"

type Event interface {
	Time() time.Time              // time stamp at which the event will execute
	Name() string                 // Name of the event for logging purposes
	Execute(ctx *context.Context) []Event // executes the event -> called by the engine running the logic
}
