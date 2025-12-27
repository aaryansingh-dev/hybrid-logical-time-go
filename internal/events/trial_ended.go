package events

import (
	"fmt"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

type TrialEnded struct {
	at       time.Time
	name     string
	customer string
}

func NewTrialEnded(at time.Time, customer string) *TrialEnded {
	return &TrialEnded{at: at, customer: customer}
}

func (e *TrialEnded) Time() time.Time {
	return e.at
}

func (e *TrialEnded) Name() string {
	return e.name
}

func (e *TrialEnded) Execute(timeProvider clock.TimeProvider) []engine.Event {
	fmt.Printf("[%s] Trial ended for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

	// Schedule invoice
	invoiceTime := timeProvider.Now().Add(1 * time.Hour)
	return []engine.Event{
		NewInvoiceCreated(invoiceTime, e.customer),
	}
}
