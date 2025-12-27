package events

import (
	"fmt"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine")

    
type InvoiceCreated struct {
	at       time.Time
	name     string
	customer string
}

func NewInvoiceCreated(at time.Time, customer string) *InvoiceCreated {
	return &InvoiceCreated{at: at, customer: customer}
}

func (e *InvoiceCreated) Time() time.Time {
	return e.at
}

func (e *InvoiceCreated) Name() string {
	return GetTypeName(e)
}

func (e *InvoiceCreated) Execute(timeProvider clock.TimeProvider) []engine.Event {
	fmt.Printf("[%s] Invoice created for %s\n", timeProvider.Now().Format(time.RFC3339), e.customer)

	paymentTime := timeProvider.Now().Add(10 * time.Minute)
	return []engine.Event{
		NewPaymentAttempt(paymentTime, e.customer),
	}
}
