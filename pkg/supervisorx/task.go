package supervisorx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type TaskFactory func(context.Context) Task

//go:generate mockgen -destination=task_mock_test.go -package=supervisorx_test . Task
type Task interface {
	Run() error
	Shutdown() error
}

type taskWrapper struct {
	task         Task
	shutdownChan chan struct{}
	wg           *sync.WaitGroup
	isCancelled  atomic.Bool
}

func newTaskWrapper(task Task, shutdownChan chan struct{}, wg *sync.WaitGroup) *taskWrapper {
	return &taskWrapper{
		task:         task,
		shutdownChan: shutdownChan,
		isCancelled:  atomic.Bool{},
		wg:           wg,
	}
}

func (t *taskWrapper) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in task :%v", r)
		}
		t.isCancelled.Store(true)
		if err != nil {
			t.shutdownChan <- struct{}{}
		}
		t.wg.Done()
	}()

	return t.task.Run()
}

func (t *taskWrapper) Shutdown() error {
	defer t.wg.Done()

	if !t.isCancelled.CompareAndSwap(false, true) {
		return nil
	}

	return t.task.Shutdown()
}
