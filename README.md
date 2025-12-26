# Hybrid Logical Time (HLT)
### A Stripe-inspired virtualized time orchestration engine for billing simulations.

**Hybrid Logical Time (HLT)** is a high-fidelity simulation engine designed to solve the "Observability of Future States" problem in complex billing systems. Inspired by Stripe’s "Test Clocks," this project allows developers to teleport through months of subscription lifecycles in milliseconds. This is a small scale project, attempting to simulate what engineers at Stripe did to overcome the challenge.

Stripe Blog Link: https://stripe.com/blog/test-clocks-how-we-made-it-easier-to-test-stripe-billing-integrations

Author: Aaryan Singh
---

## The Problem: The "Causality Loop"
In billing systems, time is not linear—it's reactive. 
1. You jump 30 days into the future.
2. A **Trial Ended** event fires at Day 14.
3. This triggers an **Invoice Created** event.
4. If the payment fails, a **Retry Scheduled** event is created for Day 17.
5. **Waiting is expensive**: You can't wait 12 months to test a annual renewal.

A simple "Date.now()" override fails because it doesn't account for these **cascading events** created *during* the time jump. HLT solves this using a recursive **Advance Loop** and a **Min-Heap Priority Queue**.

## Architecture

The system is built on three core pillars:

1. **Abstract Time Provider:** Decouples business logic from the system clock (`time.Now()`).
2. **Min-Heap Event Queue:** Stores scheduled tasks (Invoices, Retries, Webhooks) and ensures they are executed in strict chronological order ($O(\log n)$ efficiency).
3. **The Advance Loop:** A recursive orchestrator that processes events until the target "frozen time" is reached.


### Project Structure
```text
├── cmd/clock/          # Entry point for the CLI demo
├── internal/
│   ├── clock/          # HLC (Hybrid Logical Clock) implementation
│   ├── engine/         # The Advance Loop & Priority Queue logic
│   ├── billing/        # Mock business logic (Subscriptions, Invoices)
│   └── events/         # Event definitions and handlers