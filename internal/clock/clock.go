package clock

import "time"

// minimal implementation for now

type TestClock struct{
	now time.Time
}

func NewTestClock(start time.Time) *TestClock{
	return &TestClock{now: start}
} 

func (c *TestClock) Now() time.Time{
	return c.now
}

func (c *TestClock) Set(t time.Time){
	c.now = t
}


