package timex

import "time"

type SystemClock struct{}

func (s *SystemClock) Now() time.Time {
	return time.Now()
}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}
