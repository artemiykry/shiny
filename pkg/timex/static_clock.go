package timex

import "time"

type StaticClock struct {
	Time time.Time
}

func (c *StaticClock) Now() time.Time {
	return c.Time
}

func NewStaticClock(t time.Time) *StaticClock {
	return &StaticClock{Time: t}
}
