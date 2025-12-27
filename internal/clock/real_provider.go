package clock

import "time"

type RealTimeProvider struct{}

func NewRealTimeProvider() *RealTimeProvider{
	return &RealTimeProvider{}
}

func (realTimeProvider *RealTimeProvider) Now() time.Time{
	return time.Now().UTC()
}
