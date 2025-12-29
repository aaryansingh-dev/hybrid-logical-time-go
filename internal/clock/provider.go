package clock

import "time"

// TimeProvider defines an interface for retrieving the current time.
// This abstraction allows the engine to treat real-time and virtual-time
// providers interchangeably.
type TimeProvider interface {
	// Now returns the current time relative to the provider's implementation.
	// In a real-time context, this returns the system wall-clock. In a
	// simulation context, this returns the currently held logical time.
	Now() time.Time
}
