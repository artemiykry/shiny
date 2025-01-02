package logx

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type handlerKey struct{}

func WithHandler(ctx context.Context, l slog.Handler) context.Context {
	return context.WithValue(ctx, handlerKey{}, l)
}

func HandlerFromContext(ctx context.Context) (slog.Handler, bool) {
	h, ok := ctx.Value(handlerKey{}).(slog.Handler)
	return h, ok
}

func WithAttrs(ctx context.Context, v ...Attr) context.Context {
	l, ok := HandlerFromContext(ctx)
	if !ok {
		panic("handler is not set in the context")
	}
	l = l.WithAttrs(v)
	return WithHandler(ctx, l)
}

func WithGroup(ctx context.Context, groupName string) context.Context {
	l, ok := HandlerFromContext(ctx)
	if !ok {
		panic("handler is not set in the context")
	}
	l = l.WithGroup(groupName)
	return WithHandler(ctx, l)
}

func Debug(ctx context.Context, msg string, args ...Attr) {
	log(ctx, msg, slog.LevelDebug, args...)
}

func Info(ctx context.Context, msg string, args ...Attr) {
	log(ctx, msg, slog.LevelInfo, args...)
}

func Warn(ctx context.Context, msg string, args ...Attr) {
	log(ctx, msg, slog.LevelWarn, args...)
}

func Error(ctx context.Context, msg string, args ...Attr) {
	log(ctx, msg, slog.LevelError, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...Attr) {
	log(ctx, msg, level, args...)
}

func log(ctx context.Context, msg string, level slog.Level, args ...Attr) {
	l, ok := HandlerFromContext(ctx)
	if !ok {
		panic("handler is not set in the context")
	}
	if !l.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc := pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(args...)
	_ = l.Handle(ctx, r)
}
