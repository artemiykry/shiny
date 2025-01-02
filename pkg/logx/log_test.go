package logx_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/artemiykry/shiny/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestWithHandlerAndHandlerFromContext(t *testing.T) {
	_, handler := getTestHandler()
	ctx := context.Background()
	ctx = logx.WithHandler(ctx, handler)
	ctxHandler, ok := logx.HandlerFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, handler, ctxHandler)
}

func TestHandlerFromContextWithoutHandlerInContext(t *testing.T) {
	ctx := context.Background()
	value, ok := logx.HandlerFromContext(ctx)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestLog(t *testing.T) {
	type testCase struct {
		logxLevel     logx.Level
		expectedLevel string
	}
	tests := []testCase{
		{logx.LevelDebug, "DEBUG"},
		{logx.LevelInfo, "INFO"},
		{logx.LevelWarn, "WARN"},
		{logx.LevelError, "ERROR"},
	}

	for _, tc := range tests {
		buf, handler := getTestHandler()
		ctx := context.Background()
		ctx = logx.WithHandler(ctx, handler)
		logx.Log(ctx, tc.logxLevel, "hello world", logx.String("key", "value"))

		str := buf.String()
		expected := fmt.Sprintf("time=2025-01-02T17:07:05Z+03 level=%s msg=\"hello world\" key=value\n", tc.expectedLevel)
		assert.Equal(t, expected, str)
	}
}

func TestLogFuncs(t *testing.T) {
	type testCase struct {
		logFn         func(ctx context.Context, msg string, args ...logx.Attr)
		expectedLevel string
	}
	tests := []testCase{
		{logx.Debug, "DEBUG"},
		{logx.Info, "INFO"},
		{logx.Warn, "WARN"},
		{logx.Error, "ERROR"},
	}

	for _, tc := range tests {
		buf, handler := getTestHandler()
		ctx := context.Background()
		ctx = logx.WithHandler(ctx, handler)

		tc.logFn(ctx, "hello world", logx.String("key", "value"))

		str := buf.String()
		expected := fmt.Sprintf("time=2025-01-02T17:07:05Z+03 level=%s msg=\"hello world\" key=value\n", tc.expectedLevel)
		assert.Equal(t, expected, str)
	}
}

func TestWithAttrs(t *testing.T) {
	buf, handler := getTestHandler()
	ctx := context.Background()
	ctx = logx.WithHandler(ctx, handler)
	ctx = logx.WithAttrs(ctx, logx.Int("int", 42))

	logx.Info(ctx, "hello world", logx.String("key", "value"))

	str := buf.String()
	expected := "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"hello world\" int=42 key=value\n"
	assert.Equal(t, expected, str)
}

func TestWithGroup(t *testing.T) {
	buf, handler := getTestHandler()
	ctx := context.Background()
	ctx = logx.WithHandler(ctx, handler)
	ctx = logx.WithGroup(ctx, "testGroup")
	ctx = logx.WithAttrs(ctx, logx.Int("int", 42))

	logx.Info(ctx, "hello world", logx.String("key", "value"))

	str := buf.String()
	expected := "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"hello world\" testGroup.int=42 testGroup.key=value\n"
	assert.Equal(t, expected, str)
}

func TestWithAttrsWithoutHandlerInContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "handler is not set in the context", r)
			return
		}
	}()

	ctx := context.Background()
	_ = logx.WithAttrs(ctx, logx.String("key", "value"))
	t.Errorf("expected panic because no handler is set in the context")
}

func TestWithGroupWithoutHandlerInContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "handler is not set in the context", r)
			return
		}
	}()

	ctx := context.Background()
	_ = logx.WithGroup(ctx, "testGroup")
	t.Errorf("expected panic because no handler is set in the context")
}

func TestLogWithoutHandlerInContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "handler is not set in the context", r)
			return
		}
	}()

	ctx := context.Background()
	logx.Info(ctx, "hello world")
	t.Errorf("expected panic because no handler is set in the context")
}

func getTestHandler() (*bytes.Buffer, slog.Handler) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.String("time", "2025-01-02T17:07:05Z+03")
			}
			return a
		}},
	)
	return buf, handler
}
