// Package engine provides the core orchestration logic for the Hybrid Logical Time system.
package engine

import (
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
)

// Event defines the behavioral contract for discrete actions within the simulation.
// Implementations of this interface represent stateful business logic that 
// executes at a specific point in logical or real time.
type Event interface {
	// Time returns the specific moment in time when this event is scheduled to run.
	// In a simulation partition, the engine will "teleport" the clock to this 
	// exact timestamp before execution.
	Time() time.Time

	// Name returns a human-readable identifier for the event.
	// This is primarily used for diagnostic tracing and execution logging.
	Name() string

	// ClockID returns the unique identifier of the partition this event belongs to.
	// This allows the engine to route the event to the correct multi-tenant heap.
	ClockID() string

	// Execute runs the domain logic associated with the event.
	// It accepts a TimeProvider to allow logic to be aware of the current 
	// logical time. It returns a slice of future events to be scheduled, 
	// enabling the modeling of causal chains (e.g., TrialEnd triggering InvoiceCreated).
	Execute(timeProvider clock.TimeProvider) []Event
}