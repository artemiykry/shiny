package logx

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
)

type testingT interface {
	Helper()
	Log(...any)
}

type TestingHandler struct {
	t       testingT
	buf     *bytes.Buffer
	mu      *sync.Mutex
	handler slog.Handler
}

func NewTestingLogger(t testingT, opts *slog.HandlerOptions) *TestingHandler {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, opts)

	return &TestingHandler{
		t:       t,
		buf:     buf,
		mu:      &sync.Mutex{},
		handler: handler,
	}
}

func (f *TestingHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (f *TestingHandler) Handle(ctx context.Context, record slog.Record) error {
	f.t.Helper()

	f.mu.Lock()
	defer f.mu.Unlock()

	err := f.handler.Handle(ctx, record)
	if err != nil {
		return err
	}

	output, err := io.ReadAll(f.buf)
	if err != nil {
		return err
	}

	output = bytes.TrimSuffix(output, []byte("\n"))
	f.t.Log(string(output))
	return nil
}

func (f *TestingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TestingHandler{
		t:       f.t,
		buf:     f.buf,
		mu:      f.mu,
		handler: f.handler.WithAttrs(attrs),
	}
}

func (f *TestingHandler) WithGroup(name string) slog.Handler {
	return &TestingHandler{
		t:       f.t,
		buf:     f.buf,
		mu:      f.mu,
		handler: f.handler.WithGroup(name),
	}
}
