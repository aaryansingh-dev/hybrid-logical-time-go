package events

import (
    "fmt"
    "time"

    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type SubscriptionCreated struct {
    at time.Time
	name string
    customer string
    trialDays int
}

func NewSubscriptionCreated(at time.Time, customer string, trialDays int) *SubscriptionCreated {
    return &SubscriptionCreated{at: at, customer: customer, trialDays: trialDays}
}

func (e *SubscriptionCreated) Time() time.Time {
    return e.at
}

func (e *SubscriptionCreated) Name() string{
	return e.name
}

func (e *SubscriptionCreated) Execute(timeProvider clock.TimeProvider) []engine.Event {
    fmt.Printf("[%s] Subscription created for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

    return []engine.Event{}
}
