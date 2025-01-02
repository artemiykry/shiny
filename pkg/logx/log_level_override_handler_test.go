package logx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/artemiykry/shiny/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestLogLevelOverrideHandler(t *testing.T) {
	var handler slog.Handler
	buf, handler := getTestHandler()
	handler = logx.NewLogLevelOverrideHandler(handler, logx.LevelWarn)
	ctx := logx.WithHandler(context.Background(), handler)

	buf.Reset()
	logx.Debug(ctx, "test message 1")
	logStr := buf.String()
	assert.Empty(t, logStr)

	buf.Reset()
	logx.Info(ctx, "test message 2")
	logStr = buf.String()
	assert.Empty(t, logStr)

	logx.Warn(ctx, "test message 3")
	logStr = buf.String()
	assert.Equal(t, "time=2025-01-02T17:07:05Z+03 level=WARN msg=\"test message 3\"\n", logStr)

	buf.Reset()
	logx.Error(ctx, "test message 4")
	logStr = buf.String()
	assert.Equal(t, "time=2025-01-02T17:07:05Z+03 level=ERROR msg=\"test message 4\"\n", logStr)
}
