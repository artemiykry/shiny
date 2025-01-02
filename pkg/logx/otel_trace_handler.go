package logx

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

type OtelTraceHandler struct {
	handler slog.Handler
}

func NewOtelTraceHandler(handler slog.Handler) *OtelTraceHandler {
	return &OtelTraceHandler{handler: handler}
}

func (o *OtelTraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return o.handler.Enabled(ctx, level)
}

func (o *OtelTraceHandler) Handle(ctx context.Context, record slog.Record) (_ error) {
	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return o.handler.Handle(ctx, record)
	}

	spanContext := span.SpanContext()
	if spanContext.HasTraceID() {
		traceID := spanContext.TraceID().String()
		record.AddAttrs(slog.String(traceIDKey, traceID))
	}

	if spanContext.HasSpanID() {
		spanID := spanContext.SpanID().String()
		record.AddAttrs(slog.String(spanIDKey, spanID))
	}

	return o.handler.Handle(ctx, record)
}

func (o *OtelTraceHandler) WithAttrs(attrs []slog.Attr) (_ slog.Handler) {
	return &OtelTraceHandler{handler: o.handler.WithAttrs(attrs)}
}

func (o *OtelTraceHandler) WithGroup(name string) (_ slog.Handler) {
	return &OtelTraceHandler{handler: o.handler.WithGroup(name)}
}
