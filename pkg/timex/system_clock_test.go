package timex_test

import (
	"testing"
	"time"

	"github.com/artemiykry/shiny/pkg/timex"
	"github.com/stretchr/testify/assert"
)

func TestSystemClock(t *testing.T) {
	clock := timex.NewSystemClock()

	clockNow := clock.Now()

	assertTimeInRange(t, clockNow, time.Now().Add(-100*time.Millisecond), time.Now())
}

func assertTimeInRange(t assert.TestingT, actual time.Time, start, end time.Time) {
	assert.Truef(t, actual.After(start), "expected %v to be after %v", actual, start)
	assert.Truef(t, actual.Before(end), "expected %v to be before %v", actual, end)
}
