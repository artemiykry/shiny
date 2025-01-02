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

func TestTestingLogger(t *testing.T) {
	testingT := testingTMock{
		buf: &bytes.Buffer{},
	}

	handler := logx.NewTestingLogger(&testingT, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.String("time", "2025-01-02T17:07:05Z+03")
			}
			return a
		},
	})

	ctx := logx.WithHandler(context.Background(), handler)

	testingT.buf.Reset()
	logx.Debug(ctx, "This is an debug message")
	output := testingT.buf.String()
	expectedStr := "time=2025-01-02T17:07:05Z+03 level=DEBUG msg=\"This is an debug message\"\n"
	assert.Equal(t, expectedStr, output)

	testingT.buf.Reset()
	logx.Info(ctx, "This is an info message")
	output = testingT.buf.String()
	expectedStr = "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"This is an info message\"\n"
	assert.Equal(t, expectedStr, output)

	testingT.buf.Reset()
	logx.Warn(ctx, "This is an warn message")
	output = testingT.buf.String()
	expectedStr = "time=2025-01-02T17:07:05Z+03 level=WARN msg=\"This is an warn message\"\n"
	assert.Equal(t, expectedStr, output)

	testingT.buf.Reset()
	logx.Error(ctx, "This is an error message")
	output = testingT.buf.String()
	expectedStr = "time=2025-01-02T17:07:05Z+03 level=ERROR msg=\"This is an error message\"\n"
	assert.Equal(t, expectedStr, output)
}

type testingTMock struct {
	buf *bytes.Buffer
}

func (t *testingTMock) Log(args ...any) {
	s := fmt.Sprintln(args...)
	t.buf.WriteString(s)
}

func (t *testingTMock) Helper() {
}
