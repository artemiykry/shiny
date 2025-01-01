package timex_test

import (
	"testing"
	"time"

	"github.com/artemiykry/shiny/pkg/timex"
	"github.com/stretchr/testify/assert"
)

func TestStaticClockNow(t *testing.T) {
	mockTime := time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
	clock := timex.NewStaticClock(mockTime)

	now := clock.Now()

	assert.Equal(t, mockTime, now)
}
