package billing

import (
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type InvoiceCreated struct {
	scheduledAt time.Time
	customerID  string
	partitionID string
}

func NewInvoiceCreated(at time.Time, customerID string, partitionID string) *InvoiceCreated {
	return &InvoiceCreated{
		scheduledAt: at,
		customerID:  customerID,
		partitionID: partitionID,
	}
}

func (billingEvent *InvoiceCreated) Time() time.Time { return billingEvent.scheduledAt }
func (billingEvent *InvoiceCreated) Name() string    { return getEventName(billingEvent) }
func (billingEvent *InvoiceCreated) ClockID() string { return billingEvent.partitionID }

func (billingEvent *InvoiceCreated) Execute(timeProvider clock.TimeProvider) []engine.Event {
	// Attempt payment 10 minutes after invoice generation
	paymentTime := timeProvider.Now().Add(10 * time.Minute)

	return []engine.Event{
		NewPaymentAttempt(paymentTime, billingEvent.customerID, billingEvent.partitionID, 0),
	}
}
