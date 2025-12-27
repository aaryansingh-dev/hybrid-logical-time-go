package events

import (
    "fmt"
    "time"
    "math/rand"

    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
    "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type PaymentAttempt struct {
    at time.Time
	name string
    customer string
}

func NewPaymentAttempt(at time.Time, customer string) *PaymentAttempt {
    return &PaymentAttempt{at: at, customer: customer}
}

func (e *PaymentAttempt) Time() time.Time {
    return e.at
}

func (e *PaymentAttempt) Name() string{
	return GetTypeName(e)
}

func (e *PaymentAttempt) Execute(timeProvider clock.TimeProvider) []engine.Event {
    fmt.Printf("[%s] Payment attempted for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

	// 20% chance of failure → schedule retry
    if rand.Float64() < 0.2 {
        retryTime := timeProvider.Now().Add(1 * time.Hour) // retry after 1 hour
        fmt.Printf("[%s] Payment failed for %s, retry scheduled at %s\n", 
            timeProvider.Now().Format(time.RFC3339), e.customer, retryTime.Format(time.RFC3339))

        return []engine.Event{
            NewPaymentAttempt(retryTime, e.customer),
        }
    }

    fmt.Printf("[%s] Payment succeeded for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)
    return nil
}
