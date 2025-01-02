package logx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/artemiykry/shiny/pkg/logx"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestTraceHandler(t *testing.T) {
	var handler slog.Handler
	buf, handler := getTestHandler()
	handler = logx.NewOtelTraceHandler(handler)

	ctx := logx.WithHandler(context.Background(), handler)

	logx.Info(ctx, "test message")
	logStr := buf.String()
	assert.Equal(t, "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"test message\"\n", logStr)

	ctx = trace.ContextWithSpan(ctx, noopSpanForTest{})

	buf.Reset()
	logx.Info(ctx, "test message 2")
	logStr = buf.String()
	assert.Equal(t, "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"test message 2\" trace_id=01010101010101010101010101010101 span_id=0101010101010101\n", logStr)

	ctx = logx.WithAttrs(ctx, logx.String("key", "value"))
	ctx = logx.WithGroup(ctx, "groupName")
	buf.Reset()
	logx.Info(ctx, "test message 3", slog.Int("int", 42))
	logStr = buf.String()

	// TODO: groupName.span_id should be span_id
	assert.Equal(t, "time=2025-01-02T17:07:05Z+03 level=INFO msg=\"test message 3\" key=value groupName.int=42 groupName.trace_id=01010101010101010101010101010101 groupName.span_id=0101010101010101\n", logStr)
}

type noopSpanForTest struct {
	embedded.Span
}

func (se noopSpanForTest) End(options ...trace.SpanEndOption) {
}

func (se noopSpanForTest) AddEvent(name string, options ...trace.EventOption) {
}

func (se noopSpanForTest) AddLink(link trace.Link) {
}

func (se noopSpanForTest) IsRecording() bool {
	return true
}

func (se noopSpanForTest) RecordError(err error, options ...trace.EventOption) {
}

func (se noopSpanForTest) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    [16]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		SpanID:     [8]byte{1, 1, 1, 1, 1, 1, 1, 1},
		TraceFlags: 0,
		TraceState: trace.TraceState{},
		Remote:     false,
	})
}

func (se noopSpanForTest) SetStatus(code codes.Code, description string) {
}

func (se noopSpanForTest) SetName(name string) {
}

func (se noopSpanForTest) SetAttributes(kv ...attribute.KeyValue) {

}

func (se noopSpanForTest) TracerProvider() trace.TracerProvider {
	return &noop.TracerProvider{}
}
