// Package clock provides abstractions for time measurement.
// It allows the system to switch between real wall-clock time and
// deterministic virtual time for simulations.
package clock

import "time"

// TestClock is a stateful virtual clock that only advances when explicitly
// commanded. It is used in deterministic simulations to "teleport" between
// scheduled events without waiting for real-world time to pass.
type TestClock struct {
	now time.Time
}

// NewTestClock creates and returns a new TestClock initialized to the
// provided starting timestamp.
func NewTestClock(start time.Time) *TestClock {
	return &TestClock{now: start}
}

// Now returns the current logical time held by the TestClock.
// This satisfies the TimeProvider interface.
func (c *TestClock) Now() time.Time {
	return c.now
}

// Set updates the internal logical time of the clock to the provided timestamp.
// This is typically called by the engine during a "temporal jump" or "causal walk."
func (c *TestClock) Set(t time.Time) {
	c.now = t
}
