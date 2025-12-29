package billing

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type PaymentAttempt struct {
	scheduledAt  time.Time
	customerID   string
	partitionID  string
	currentRetry int
}

func NewPaymentAttempt(at time.Time, customerID string, partitionID string, retryCount int) *PaymentAttempt {
	return &PaymentAttempt{
		scheduledAt:  at,
		customerID:   customerID,
		partitionID:  partitionID,
		currentRetry: retryCount,
	}
}

func (billingEvent *PaymentAttempt) Time() time.Time { return billingEvent.scheduledAt }
func (billingEvent *PaymentAttempt) Name() string    { return getEventName(billingEvent) }
func (billingEvent *PaymentAttempt) ClockID() string { return billingEvent.partitionID }

func (billingEvent *PaymentAttempt) Execute(timeProvider clock.TimeProvider) []engine.Event {
	const MaxRetries = 3

	// Simulation: 20% failure rate to demonstrate error handling causality
	if rand.Float64() < 0.2 {
		if billingEvent.currentRetry >= MaxRetries {
			fmt.Printf("[BILLING] CRITICAL: Payment permanently failed for customer %s after %d retries\n",
				billingEvent.customerID, billingEvent.currentRetry)
			return nil
		}

		// Implement Linear Backoff: Each retry waits longer (1hr, 2hr, 3hr...)
		backoffDuration := time.Duration(billingEvent.currentRetry+1) * time.Hour
		retryTime := timeProvider.Now().Add(backoffDuration)

		return []engine.Event{
			NewPaymentAttempt(retryTime, billingEvent.customerID, billingEvent.partitionID, billingEvent.currentRetry+1),
		}
	}

	fmt.Printf("[BILLING] SUCCESS: Payment processed for %s at %s\n",
		billingEvent.customerID, timeProvider.Now().Format(time.RFC3339))

	nextInvoiceCycle := timeProvider.Now().AddDate(0, 1, 0) // 1 month later
	
	return []engine.Event{
		NewInvoiceCreated(nextInvoiceCycle, billingEvent.customerID, billingEvent.partitionID),
	}
}
