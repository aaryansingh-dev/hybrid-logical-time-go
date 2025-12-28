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
	OnEventCreated(id string, eventName string, t time.Time)
	OnAdvanceFinish(id string, current time.Time)
}

// ConsoleLogger is one specific implementation of the Diagnostic interface.
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) OnAdvanceStart(id string, start, target time.Time) {
	fmt.Printf("\n[START]  %-12s | Advancing from %s -> %s\n", id, start.Format("15:04:05"), target.Format("15:04:05"))
	fmt.Println(strings.Repeat("-", 80))
}

func (c *ConsoleLogger) OnEventExecute(id string, eventName string, t time.Time) {
	fmt.Printf("[%s] EXEC  | %-12s | Event: %s\n", t.Format("15:04:05.000"), id, eventName)
}

func (c *ConsoleLogger) OnEventCreated(id string, eventName string, t time.Time) {
	fmt.Printf("[%s] CHAIN | %-12s | Created: %s\n", t.Format("15:04:05.000"), id, eventName)
}

func (c *ConsoleLogger) OnAdvanceFinish(id string, current time.Time) {
	fmt.Printf("[FINISH] %-12s | Simulation paused at %s\n", id, current.Format("15:04:05"))
}
