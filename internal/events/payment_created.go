package events

import (
    "fmt"
    "time"

    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type PaymentAttempt struct {
    at time.Time
	name string
    customer string
}

func NewPaymentCreated(at time.Time, customer string, trialDays int) *SubscriptionCreated {
    return &SubscriptionCreated{at: at, customer: customer, trialDays: trialDays}
}

func (e *PaymentAttempt) Time() time.Time {
    return e.at
}

func (e *PaymentAttempt) Name() string{
	return GetTypeName(e)
}

func (e *PaymentAttempt) Execute(timeProvider clock.TimeProvider) []engine.Event {
    fmt.Printf("[%s] Payment attempted for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

	return nil
}
