package clock

import "time"

type TimeProvider interface{
	Now() time.Time
}
