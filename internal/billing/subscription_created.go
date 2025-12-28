package billing

import (
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type SubscriptionCreated struct {
	scheduledAt time.Time
	customerID  string
	trialDays   int
	partitionID string
}

func NewSubscriptionCreated(at time.Time, customerID string, trialDays int, partitionID string) *SubscriptionCreated {
	return &SubscriptionCreated{
		scheduledAt: at,
		customerID:  customerID,
		trialDays:   trialDays,
		partitionID: partitionID,
	}
}

func (billingEvent *SubscriptionCreated) Time() time.Time { return billingEvent.scheduledAt }
func (billingEvent *SubscriptionCreated) Name() string    { return getEventName(billingEvent) }
func (billingEvent *SubscriptionCreated) ClockID() string { return billingEvent.partitionID }

func (billingEvent *SubscriptionCreated) Execute(timeProvider clock.TimeProvider) []engine.Event {
	// Calculate trial end date based on current logical time
	trialEnd := timeProvider.Now().AddDate(0, 0, billingEvent.trialDays)

	return []engine.Event{
		NewTrialEnded(trialEnd, billingEvent.customerID, billingEvent.partitionID),
	}
}
