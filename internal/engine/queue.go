package engine

import "container/heap"

type EventQueue struct{
	events []Event	// simply a slice which will be turned into a heap
}

func (q EventQueue) Len() int{
	return len(q.events)
}

func (q EventQueue) Less(i, j int) bool{
	return q.events[i].Time().Before(q.events[j].Time())
}

func (q EventQueue) Swap(i, j int){
	q.events[i], q.events[j] = q.events[j], q.events[i]
}

func (q *EventQueue) Push(x interface{}){
	event := x.(Event)		// converting this back to event. Heap expects and interface{} type i.e. any datatype. But, we need to tell the compiler that we know that the x is Event.
	q.events = append(q.events, event)
}

func (q *EventQueue) Pop() interface{}{
	n := len(q.events)
	event := q.events[n-1]
	q.events = q.events[:n-1]
	return event
}

// helpers to be made

func NewEventQueue() *EventQueue {
    q := &EventQueue{
        events: []Event{},
    }
    heap.Init(q)
    return q
}

func (q *EventQueue) PushEvent(e Event) {
    heap.Push(q, e)
}

func (q *EventQueue) PopEvent() Event {
    return heap.Pop(q).(Event)
}

func (q *EventQueue) Peek() Event {
    if len(q.events) == 0 {
        return nil
    }
    return q.events[0]
}

