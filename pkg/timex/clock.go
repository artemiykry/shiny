package timex

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
}

func Now(ctx context.Context) time.Time {
	clock, ok := ClockFromContext(ctx)
	if !ok {
		panic("clock is not found in context")
	}
	return clock.Now()
}

type clockKey struct{}

func WithClock(ctx context.Context, clock Clock) context.Context {
	return context.WithValue(ctx, clockKey{}, clock)
}

func ClockFromContext(ctx context.Context) (Clock, bool) {
	clock, ok := ctx.Value(clockKey{}).(Clock)
	return clock, ok
}
