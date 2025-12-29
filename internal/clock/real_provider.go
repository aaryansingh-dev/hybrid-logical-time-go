package clock

import "time"

// RealTimeProvider satisfies the TimeProvider interface using the host's
// actual wall-clock. It is used in production or "SYSTEM" partitions
// where events must follow the real passage of time.
type RealTimeProvider struct{}

// NewRealTimeProvider initializes and returns a new RealTimeProvider.
func NewRealTimeProvider() *RealTimeProvider {
	return &RealTimeProvider{}
}

// Now returns the current UTC time from the system clock.
// This ensures that all real-time event comparisons remain consistent
// regardless of the host's local timezone settings.
func (realTimeProvider *RealTimeProvider) Now() time.Time {
	return time.Now().UTC()
}
