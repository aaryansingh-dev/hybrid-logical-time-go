# ⏳ HLT (Hybrid Logical TestClocks)

**Deterministic Temporal Engine for Multi-Tenant State Simulations**

HLT is a high-fidelity simulation engine built in Go. Inspired by the engineering principles behind Stripe Test Clocks, this project provides a technical solution for decoupling application logic from the linear progression of wall-clock time.

It enables engineering teams to execute years of stateful business transitions—such as subscription billing lifecycles and complex dunning sequences—in milliseconds with full causal determinism.

Stripe Blog Link: https://stripe.com/blog/test-clocks-how-we-made-it-easier-to-test-stripe-billing-integrations

---

## 🎯 The Engineering Challenge: Temporal Friction

In modern SaaS infrastructure, testing time-dependent systems traditionally involves significant friction. Historically, verifying billing integrations meant waiting for real-world time to pass or using "shaky foundations" like 10-second trials and shortened subscription cycles that do not perfectly mirror production. HLT addresses these specific hurdles:

### Common Failure Modes

**Temporal Bottleneck**  
Development cycles stall when teams cannot instantly verify long-term transitions such as 14-day trials or 30-day billing resets inside a standard CI/CD pipeline.

**Non-Deterministic Flakiness**  
Relying on `time.Now()` makes integration tests dependent on host CPU scheduling and system interrupts, leading to unstable results and race conditions.

**Causal Side Effects**  
Direct database timestamp manipulation often misses cascading effects where one execution (e.g. trial expiration) must logically trigger another (e.g. invoice generation).

---

## The Solution: Out-of-Band Temporal Testing

HLT moves beyond basic mocking by implementing a **Discrete Event Simulation (DES)** runtime. Time is treated as a controllable, partitioned dependency rather than a global constant.

### Core Capabilities

**Temporal Partitioning**  
Create isolated timelines per tenant. One partition may remain frozen for debugging while another warps through a one-year simulation.

**Deterministic Causal Walks**  
Events are processed in strict chronological order. If Event A schedules a side effect at `T+2`, the engine guarantees that the side effect is discovered and executed before time advances further.

**High-Efficiency Scaling**  
By moving time out-of-band, HLT achieves orders-of-magnitude speedups in test execution. Performance is bounded by event density rather than temporal duration.

---

## 📐 Architecture

HLT uses a Min-Heap priority queue as the source of truth for event ordering, combined with interface-based clock providers.

```[ Orchestrator ]
       |
       |--> [ SYSTEM Partition ] ----> RealTimeProvider (Wall-Clock)
       |
        --> [ TEST_CLOCK Partitions ] ----> TestClockProvider (Virtual-Time)
                    |
                     --> [ Event Heap ] --> [ Logic Execution ] --> [ Future Event Injection ]
```
---

## 🔄 The Causal Walk Algorithm

HLT does not simply “skip” to a future date. It performs a recursive execution loop to preserve state integrity.

1. **Boundary Selection**  
   The orchestrator defines a target time for the temporal advance.

2. **Sequential Polling**  
   The engine retrieves the earliest event `E` such that `E.Time ≤ TargetTime`.

3. **Clock Teleportation**  
   The partition’s internal clock is advanced exactly to `E.Time`.

4. **Execution & Injection**  
   Event logic executes. Any newly generated causal events are injected back into the heap for immediate re-ordering.

5. **Recursion**  
   The loop repeats until the heap is empty or the next event exceeds the target time.

This guarantees deterministic discovery of all causal side effects.

---

## 💻 CLI Commands

The HLT binary includes an interactive shell for managing temporal partitions.

| Command | Arguments | Description |
|------|---------|-------------|
| `create-partition` | `<id> <iso_timestamp>` | Initialize a new virtual clock for a tenant |
| `schedule` | `<id> <delay> <unit> <trial_days>` | Inject a subscription event into a partition |
| `advance` | `<id> <value> <unit>` | Perform a deterministic causal walk |
| `status` | `<id>` | Display current virtual time and pending events |
| `list` | — | List all active partitions |

### Supported Time Units

- `s` — seconds   
- `h` — hours  
- `d` — days  
- `m` — months (30 days: approximation)

---

## ✨ Technical Characteristics

- **Zero Clock Drift**  
  Memory-safe state transitions using `sync.RWMutex` across concurrent partitions.

- **Recursive Event Discovery**  
  Events created during execution are discovered and processed within the same causal sweep.

- **O(log N) Scheduling**  
  Min-heap scheduling ensures efficient operation even with large event volumes.

---

## 🚀 Quick Start

### 1. Build and Run

```bash
go run ./cmd/hlt_cli/main.go
```

---

## 2. Run a Billing Simulation

Dry-run a monthly billing lifecycle in a single session:

```text
# Initialize a sandbox starting on Jan 1st
> create-partition user_42 2025-01-01T00:00:00Z

# Schedule a subscription with a 14-day trial
> schedule user_42 0 s 14

# Warp 30 days forward to observe the full causal chain
# Trial End → Invoice Created → Payment Processed
> advance user_42 30 d
```
---
## Test
Run unit tests on the core engine logic and implementation:

```bash
go test -v ./internal/engine/
```
---
## Project Structure
```
/internal/engine   # Core DES engine and scheduler
/internal/clock    # TimeProvider abstractions
/internal/billing  # Subscription state machines
/cmd/hlt_cli       # CLI entrypoint
```