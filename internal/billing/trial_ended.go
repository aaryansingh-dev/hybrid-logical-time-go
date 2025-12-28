package billing

import (
	"time"
	"fmt"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)


type TrialEnded struct {
	scheduledAt time.Time
	customerID  string
	partitionID string
}

func NewTrialEnded(at time.Time, customerID string, partitionID string) *TrialEnded {
	return &TrialEnded{
		scheduledAt: at,
		customerID:  customerID,
		partitionID: partitionID,
	}
}

func (billingEvent *TrialEnded) Time() time.Time { return billingEvent.scheduledAt }
func (billingEvent *TrialEnded) Name() string    { return getEventName(billingEvent) }
func (billingEvent *TrialEnded) ClockID() string { return billingEvent.partitionID }

func (billingEvent *TrialEnded) Execute(timeProvider clock.TimeProvider) []engine.Event {
	// Logic: When trial ends, we immediately transition to the invoicing phase.
	// We schedule the invoice for 1 hour after the trial officially concludes.
	invoiceTime := timeProvider.Now().Add(1 * time.Hour)
	
	fmt.Printf("[%s] [BILLING] Trial expired for customer %s. Transitioning to invoicing.\n", 
		timeProvider.Now().Format(time.RFC3339), billingEvent.customerID)
	
	return []engine.Event{
		NewInvoiceCreated(invoiceTime, billingEvent.customerID, billingEvent.partitionID),
	}
}