package context

import "github.com/aaryansingh-dev/hybrid-logical-time-go/internal/clock"


type Context struct{
	Clock *clock.TestClock
}
