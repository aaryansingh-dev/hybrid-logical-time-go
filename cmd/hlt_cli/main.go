package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/billing"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/engine"
)

func main() {
	// initialize the engine with a ConsoleLogger for real-time visibility
	eng := engine.NewEngine(&engine.ConsoleLogger{})

	// start the background system worker
	eng.StartRealTimeWorker(30 * time.Second)

	fmt.Println("\n🚀 HYBRID LOGICAL TIME ENGINE CLI")
	fmt.Println("=================================")
	fmt.Println("System Status: Real-Time Worker Active (30s ticks)")
	fmt.Println("\nCommands:")
	fmt.Println("  create-partition <id> <frozen_time_rfc3339>")
	fmt.Println("----- Example: create-partition user_123 2025-01-01T10:00:00Z")
	fmt.Println("  schedule <partitionID:str> <value:int> <s|h|d|m>")
	fmt.Println("  advance <partitionID:str> <value:int> <s|h|d|m>")
	fmt.Println("  status")
	fmt.Println("  quit")
	fmt.Println("---------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "create-partition":
			// Example: create-partition user_123 2025-01-01T10:00:00Z
			if len(args) < 3 {
				fmt.Println("❌ Usage: create-partition <id> <2025-01-01T10:00:00Z>")
				continue
			}
			id := args[1]
			if id == "SYSTEM" {
				fmt.Printf("❌ Error: Cannot create partition with this id. Reserved keyword.")
				continue
			}

			startTime, err := time.Parse(time.RFC3339, args[2])
			if err != nil {
				fmt.Printf("❌ Invalid time format: %v\n", err)
				continue
			}

			eng.RegisterPartition(id, clock.NewTestClock(startTime))
			fmt.Printf("✅ Registered partition '%s' starting at %s\n", id, startTime.Format(time.RFC1123))

		case "schedule":
			// Example: schedule SYSTEM 30 s
			// Example: schedule user_1 1 h 0 (0-day trial) -> 0 day trial defaults to 1 minute
			if len(args) < 4 {
				fmt.Println("❌ Usage: schedule <id> <val> <s|h|d|m> [trial_days]")
				continue
			}
			id := args[1]
			val, _ := strconv.Atoi(args[2])
			delay := parseDuration(val, args[3])

			// Default trial is 14 days unless specified
			trialDays := 14
			if len(args) > 4 {
				trialDays, _ = strconv.Atoi(args[4])
			}

			currentTime, err := eng.GetPartitionTime(id)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			startTime := currentTime.Add(delay)

			// if scheduling on SYSTEM with 0 trial, we ensure a minimum offset
			// so the worker has time to pick it up before the "trial expires" immediately.
			var trialDuration time.Duration
			if id == "SYSTEM" && trialDays == 0 {
				fmt.Println("⚠️  Note: 0-day trial on SYSTEM auto-adjusted to 1-minute for worker polling.")
				trialDuration = 1 * time.Minute
			} else {
				trialDuration = time.Duration(trialDays) * 24 * time.Hour
			}

			event := billing.NewSubscriptionCreated(startTime, "CUST-"+id, trialDuration, id)
			eng.Schedule(event)

			fmt.Printf("✅ Scheduled '%s' for %s (Trial Duration: %v)\n", id, startTime.Format(time.RFC1123), trialDuration)

		case "advance":
			// Example: advance user_123 30 d
			if len(args) < 4 {
				fmt.Println("❌ Usage: advance <partitionID> <value> <h|d|m>")
				continue
			}
			id := args[1]
			if id == "SYSTEM" {
				fmt.Printf("❌ Error: Cannot advance with this id. Reserved keyword.")
				continue
			}

			val, _ := strconv.Atoi(args[2])
			jump := parseDuration(val, args[3])

			// 1. Fetch current logical frozen time
			currentTime, err := eng.GetPartitionTime(id)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			// 2. Calculate target time relative to the partition
			target := currentTime.Add(jump)

			// 3. Execute the deterministic advance loop
			err = eng.Advance(id, target, nil)
			if err != nil {
				fmt.Printf("❌ Advance failed: %v\n", err)
			}

		case "status":
			stats := eng.GetStatus()
			fmt.Println("\n--- Engine Partition Status ---")
			for id, info := range stats {
				fmt.Printf("[%s] %s\n", id, info)
			}

		case "quit", "exit":
			fmt.Println("👋 Shutting down engine...")
			return

		default:
			fmt.Printf("❓ Unknown command: %s\n", args[0])
		}
	}
}

// parseDuration converts numeric values and unit strings into time.Duration
func parseDuration(val int, unit string) time.Duration {
	switch strings.ToLower(unit) {
	case "s":
		return time.Duration(val) * time.Second
	case "h":
		return time.Duration(val) * time.Hour
	case "d":
		return time.Duration(val) * 24 * time.Hour
	case "m":
		// Approximation: 1 month = 30 days
		return time.Duration(val) * 24 * 30 * time.Hour
	default:
		return time.Duration(val) * time.Second
	}
}
