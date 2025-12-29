package engine

import (
	"fmt"
	"strings"
	"time"
)

// Diagnostic defines the hooks for system observability.
// This allows the engine to remain agnostic of how logs are formatted.
type Diagnostic interface {
	OnAdvanceStart(id string, start, target time.Time)
	OnEventExecute(id string, eventName string, t time.Time)
	OnEventCreated(id string, eventName string, eventTime time.Time, currentTime time.Time)
	OnAdvanceFinish(id string, current time.Time)
}

// ConsoleLogger implements the Diagnostic interface with formatted stdout output.
// It provides a high-fidelity visual trace of the simulation's causal walks.
type ConsoleLogger struct{}

const logTimeFormat = "2006-01-02 15:04:05"

// OnAdvanceStart prints the initialization header for a temporal advance.
func (c *ConsoleLogger) OnAdvanceStart(id string, start, target time.Time) {
	fmt.Printf("\n[START]  %-12s | Advancing from [%s] -> [%s]\n",
		id,
		start.Format(logTimeFormat),
		target.Format(logTimeFormat))
	fmt.Println(strings.Repeat("-", 100))
}

// OnEventExecute prints the execution trace for a specific event.
func (c *ConsoleLogger) OnEventExecute(id string, eventName string, t time.Time) {
	fmt.Printf("[%s] EXEC  | %-12s | Event: %s\n",
		t.Format(logTimeFormat),
		id,
		eventName)
}

// OnEventCreated prints the details of a causal event injection.
func (c *ConsoleLogger) OnEventCreated(id string, eventName string, eventTime time.Time, currentTime time.Time) {
	//
	fmt.Printf("[%s] CHAIN | %-12s | Created: %-18s (scheduled for %s)\n\n",
		currentTime.Format(logTimeFormat),
		id,
		eventName,
		eventTime.Format(logTimeFormat))
}

// OnAdvanceFinish prints the conclusion footer for a temporal advance.
func (c *ConsoleLogger) OnAdvanceFinish(id string, current time.Time) {
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("[FINISH] %-12s | Simulation paused at [%s]\n\n",
		id,
		current.Format(logTimeFormat))
}
