package logx

import (
	"context"
	"log/slog"
)

type LogLevelOverrideHandler struct {
	next  slog.Handler
	level Level
}

func NewLogLevelOverrideHandler(next slog.Handler, level Level) *LogLevelOverrideHandler {
	return &LogLevelOverrideHandler{
		next:  next,
		level: level,
	}
}

func (l *LogLevelOverrideHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= slog.Level(l.level)
}

func (l *LogLevelOverrideHandler) Handle(ctx context.Context, record slog.Record) error {
	return l.next.Handle(ctx, record)
}

func (l *LogLevelOverrideHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogLevelOverrideHandler{next: l.next.WithAttrs(attrs)}
}

func (l *LogLevelOverrideHandler) WithGroup(name string) slog.Handler {
	return &LogLevelOverrideHandler{next: l.next.WithGroup(name)}
}
