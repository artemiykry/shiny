package timex_test

import (
	"context"
	"testing"
	"time"

	"github.com/artemiykry/shiny/pkg/timex"
	"github.com/stretchr/testify/assert"
)

type mockClock struct{}

func (m mockClock) Now() time.Time {
	return mockTime
}

var mockTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNow(t *testing.T) {
	ctx := context.Background()
	ctx = timex.WithClock(ctx, mockClock{})

	clockTime := timex.Now(ctx)

	assert.Equal(t, mockTime, clockTime)
}

func TestClockFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = timex.WithClock(ctx, mockClock{})

	clock, ok := timex.ClockFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, mockTime, clock.Now())
}

func TestClockFromContextNo(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	ctx := context.Background()
	timex.Now(ctx)
	t.Errorf("expected panic because no clock is set in the context")
}
