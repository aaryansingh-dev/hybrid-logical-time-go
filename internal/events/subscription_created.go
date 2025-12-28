package events

import (
	"fmt"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type SubscriptionCreated struct {
	at        time.Time
	name      string
	customer  string
	trialDays int
	clockID   string
}

func NewSubscriptionCreated(at time.Time, customer string, trialDays int, clockID string) *SubscriptionCreated {
	return &SubscriptionCreated{at: at, customer: customer, trialDays: trialDays, clockID: clockID}
}

func (e *SubscriptionCreated) Time() time.Time {
	return e.at
}

func (e *SubscriptionCreated) Name() string {
	return GetTypeName(e)
}

func (e *SubscriptionCreated) ClockID() string {
	return e.clockID
}

func (e *SubscriptionCreated) Execute(timeProvider clock.TimeProvider) []engine.Event {
	fmt.Printf("[%s] Subscription created for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

	// Add a trial ended event
	trialEnd := timeProvider.Now().Add(time.Duration(e.trialDays*24) * time.Hour)
	return []engine.Event{
		NewTrialEnded(trialEnd, e.customer, e.clockID),
	}
}
