package supervisorx_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/artemiykry/shiny/pkg/logx"
	"github.com/artemiykry/shiny/pkg/supervisorx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunSuccessful(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTask := NewMockTask(mockCtrl)

	mockTask.EXPECT().Run().Return(nil)

	err := supervisorx.Run(
		*slog.New(logx.NewTestingLogger(t, &slog.HandlerOptions{})),
		map[string]supervisorx.TaskFactory{
			"task": func(context.Context) supervisorx.Task { return mockTask },
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
	require.NoError(t, err)
}

func TestRunError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTask := NewMockTask(mockCtrl)

	mockTask.EXPECT().Run().Return(fmt.Errorf("somethign went wrong"))

	err := supervisorx.Run(
		*slog.New(logx.NewTestingLogger(t, &slog.HandlerOptions{})),
		map[string]supervisorx.TaskFactory{
			"task": func(context.Context) supervisorx.Task { return mockTask },
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
	require.NoError(t, err)
}

func TestRunPanic(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTask := NewMockTask(mockCtrl)

	mockTask.EXPECT().Run().DoAndReturn(func() error {
		panic("panic")
	})

	err := supervisorx.Run(
		*slog.New(logx.NewTestingLogger(t, &slog.HandlerOptions{})),
		map[string]supervisorx.TaskFactory{
			"task": func(context.Context) supervisorx.Task { return mockTask },
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
	require.NoError(t, err)
}

func TestRunIgnoreCancellation2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTask := NewMockTask(mockCtrl)
	mockTask2 := NewMockTask(mockCtrl)

	var ctx context.Context
	mockTask.EXPECT().Run().DoAndReturn(func() error {
		<-ctx.Done()
		return nil
	})
	mockTask.EXPECT().Shutdown().Return(nil)

	mockTask2.EXPECT().Run().DoAndReturn(func() error {
		time.Sleep(10 * time.Millisecond)
		return fmt.Errorf("something went wrong")
	})

	err := supervisorx.Run(
		*slog.New(logx.NewTestingLogger(t, &slog.HandlerOptions{})),
		map[string]supervisorx.TaskFactory{
			"task": func(inctx context.Context) supervisorx.Task {
				ctx = inctx
				return mockTask
			},
			"task2": func(context.Context) supervisorx.Task {
				return mockTask2
			},
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
	require.NoError(t, err)
}

func TestRunIgnoreCancellation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTask := NewMockTask(mockCtrl)
	mockTask2 := NewMockTask(mockCtrl)

	mockTask.EXPECT().Run().DoAndReturn(func() error {
		time.Sleep(time.Second)
		return nil
	})
	mockTask.EXPECT().Shutdown().Return(nil)

	mockTask2.EXPECT().Run().DoAndReturn(func() error {
		return fmt.Errorf("something went wrong")
	})

	err := supervisorx.Run(
		*slog.New(logx.NewTestingLogger(t, &slog.HandlerOptions{})),
		map[string]supervisorx.TaskFactory{
			"task": func(context.Context) supervisorx.Task {
				return mockTask
			},
			"task2": func(context.Context) supervisorx.Task {
				return mockTask2
			},
		},
		50*time.Millisecond,
		10*time.Millisecond,
	)
	require.Error(t, err, "some task failed to shut down after 60ms")
}
