package engine

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"
)

// MockDiagnostic captures engine hooks for verification
type MockDiagnostic struct {
	eventsExecuted []string
	createdEvents  []string
	mu             sync.Mutex
}

func (m *MockDiagnostic) OnAdvanceStart(id string, start, target time.Time) {}
func (m *MockDiagnostic) OnEventExecute(id string, name string, t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventsExecuted = append(m.eventsExecuted, name)
}
func (m *MockDiagnostic) OnEventCreated(id string, name string, t, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createdEvents = append(m.createdEvents, name)
}
func (m *MockDiagnostic) OnAdvanceFinish(id string, current time.Time) {}

// MockEvent satisfies the engine.Event interface
type MockEvent struct {
	executionTime time.Time
	name          string
	clockID       string
	onExecute     func(tp clock.TimeProvider) []Event
}

func (e *MockEvent) Time() time.Time { return e.executionTime }
func (e *MockEvent) Name() string    { return e.name }
func (e *MockEvent) ClockID() string { return e.clockID }
func (e *MockEvent) Execute(tp clock.TimeProvider) []Event {
	if e.onExecute != nil {
		return e.onExecute(tp)
	}
	return nil
}

func TestEngine_RegisterAndGetPartition(t *testing.T) {
	eng := NewEngine(nil)
	id := "tenant_a"
	start := time.Now()
	tc := clock.NewTestClock(start)

	eng.RegisterPartition(id, tc)

	q, p, err := eng.getPartition(id)
	if err != nil {
		t.Fatalf("Expected partition to be registered, got error: %v", err)
	}
	if q == nil || p != tc {
		t.Error("Returned partition or provider does not match registered values")
	}
}

func TestEngine_Schedule_LazyRegistration(t *testing.T) {
	eng := NewEngine(nil)
	id := "lazy_tenant"
	ev := &MockEvent{executionTime: time.Now(), name: "LazyEvent", clockID: id}

	eng.Schedule(ev)

	eng.mu.RLock()
	q, exists := eng.queues[id]
	eng.mu.RUnlock()

	if !exists {
		t.Fatal("Expected partition queue to be lazy-initialized")
	}
	if q.Len() != 1 {
		t.Errorf("Expected 1 event in lazy queue, got %d", q.Len())
	}
}

func TestEngine_Advance_CausalChain(t *testing.T) {
	diag := &MockDiagnostic{}
	eng := NewEngine(diag)
	id := "chain_tenant"
	start := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	tc := clock.NewTestClock(start)
	eng.RegisterPartition(id, tc)

	ev1 := &MockEvent{
		executionTime: start.Add(time.Hour),
		name:          "First",
		clockID:       id,
		onExecute: func(tp clock.TimeProvider) []Event {
			return []Event{&MockEvent{
				executionTime: tp.Now().Add(time.Hour),
				name:          "Second",
				clockID:       id,
			}}
		},
	}

	eng.Schedule(ev1)

	target := start.Add(3 * time.Hour)
	err := eng.Advance(id, target, nil)
	if err != nil {
		t.Fatalf("Advance failed: %v", err)
	}

	if !tc.Now().Equal(target) {
		t.Errorf("Clock did not land on target. Got %v, want %v", tc.Now(), target)
	}

	diag.mu.Lock()
	defer diag.mu.Unlock()
	if len(diag.eventsExecuted) != 2 {
		t.Errorf("Expected 2 events executed, got %d", len(diag.eventsExecuted))
	}
}

func TestEngine_SimulatedLifecycle_NoCyclicImport(t *testing.T) {
	diag := &MockDiagnostic{}
	eng := NewEngine(diag)
	id := "lifecycle_test"
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tc := clock.NewTestClock(start)
	eng.RegisterPartition(id, tc)

	// Define a chain within the test to avoid importing 'billing'
	// SubscriptionCreated -> TrialEnded -> InvoiceCreated
	subEvent := &MockEvent{
		executionTime: start,
		name:          "SubscriptionCreated",
		clockID:       id,
		onExecute: func(tp clock.TimeProvider) []Event {
			return []Event{&MockEvent{
				executionTime: tp.Now().Add(14 * 24 * time.Hour),
				name:          "TrialEnded",
				clockID:       id,
				onExecute: func(tp2 clock.TimeProvider) []Event {
					return []Event{&MockEvent{
						executionTime: tp2.Now().Add(1 * time.Hour),
						name:          "InvoiceCreated",
						clockID:       id,
					}}
				},
			}}
		},
	}

	eng.Schedule(subEvent)

	target := start.Add(20 * 24 * time.Hour)
	err := eng.Advance(id, target, nil)
	if err != nil {
		t.Fatalf("Advance failed: %v", err)
	}

	diag.mu.Lock()
	defer diag.mu.Unlock()

	expectedEvents := []string{"SubscriptionCreated", "TrialEnded", "InvoiceCreated"}
	if len(diag.eventsExecuted) != len(expectedEvents) {
		t.Errorf("Expected %d events, executed %d: %v", len(expectedEvents), len(diag.eventsExecuted), diag.eventsExecuted)
	}

	for i, name := range expectedEvents {
		if diag.eventsExecuted[i] != name {
			t.Errorf("Event mismatch at index %d: got %s, want %s", i, diag.eventsExecuted[i], name)
		}
	}
}

func TestEngine_SimultaneousEvents(t *testing.T) {
	eng := NewEngine(nil)
	id := "simultaneous"
	start := time.Now().UTC()
	tc := clock.NewTestClock(start)
	eng.RegisterPartition(id, tc)

	// Schedule 3 events at the exact same time
	execTime := start.Add(time.Hour)
	var mu sync.Mutex
	count := 0
	increment := func(tp clock.TimeProvider) []Event {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	}

	for i := 0; i < 3; i++ {
		eng.Schedule(&MockEvent{executionTime: execTime, name: "SameTime", clockID: id, onExecute: increment})
	}

	err := eng.Advance(id, execTime.Add(time.Minute), nil)
	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Errorf("Expected 3 events executed at the same timestamp, got %d", count)
	}
}

func TestEngine_Advance_Validation_Errors(t *testing.T) {
	eng := NewEngine(nil)

	// Test: Advancing non-existent partition
	err := eng.Advance("GHOST", time.Now(), nil)
	if err == nil || !strings.Contains(err.Error(), "not registered") {
		t.Errorf("Expected registration error, got: %v", err)
	}

	// Test: Advancing SYSTEM
	err = eng.Advance("SYSTEM", time.Now(), nil)
	if err == nil || !strings.Contains(err.Error(), "invalid operation") {
		t.Errorf("Expected system error, got: %v", err)
	}
}

func TestEngine_StatusReporting(t *testing.T) {
	eng := NewEngine(nil)
	id := "status_check"
	start := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	eng.RegisterPartition(id, clock.NewTestClock(start))

	eng.Schedule(&MockEvent{executionTime: start.Add(time.Hour), clockID: id})

	status := eng.GetStatus()
	info, ok := status[id]
	if !ok {
		t.Fatal("Partition missing from status")
	}

	if !strings.Contains(info, "Pending Events: 1") {
		t.Errorf("Status string incorrect: %s", info)
	}
}
