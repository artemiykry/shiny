package logx_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/artemiykry/shiny/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	attr := logx.String("key", "value")

	assert.Equal(t, slog.String("key", "value"), attr)
}

func TestInt64(t *testing.T) {
	attr := logx.Int64("key", 123)

	assert.Equal(t, slog.Int64("key", 123), attr)
}

func TestInt(t *testing.T) {
	attr := logx.Int("key", 456)

	assert.Equal(t, slog.Int("key", 456), attr)
}

func TestUint64(t *testing.T) {
	attr := logx.Uint64("key", 789)

	assert.Equal(t, slog.Uint64("key", 789), attr)
}

func TestFloat64(t *testing.T) {
	attr := logx.Float64("key", 123.45)

	assert.Equal(t, slog.Float64("key", 123.45), attr)
}

func TestBool(t *testing.T) {
	attr := logx.Bool("key", true)

	assert.Equal(t, slog.Bool("key", true), attr)
}

func TestTime(t *testing.T) {
	timeNow := time.Now()
	attr := logx.Time("key", timeNow)

	assert.Equal(t, slog.Time("key", timeNow), attr)
}

func TestDuration(t *testing.T) {
	duration := time.Second
	attr := logx.Duration("key", duration)

	assert.Equal(t, slog.Duration("key", duration), attr)
}

func TestGroup(t *testing.T) {
	attr := logx.Group("key", "value1", "value2")

	assert.Equal(t, slog.Group("key", "value1", "value2"), attr)
}

func TestAny(t *testing.T) {
	attr := logx.Any("key", map[string]string{"foo": "bar"})

	assert.Equal(t, slog.Any("key", map[string]string{"foo": "bar"}), attr)
}
