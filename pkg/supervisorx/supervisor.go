package supervisorx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(
	logger slog.Logger,
	tasks map[string]TaskFactory,
	gracefultShutdownDuration time.Duration,
	shutdownDuration time.Duration,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wgChan := wgAwait(wg)

	startShutdownChan := createShutdownChannel(len(tasks) + 1)
	wrappedTasks := make(map[string]*taskWrapper)
	for name, tf := range tasks {
		wrappedTasks[name] = newTaskWrapper(tf(ctx), startShutdownChan, wg)
	}

	for taskName, task := range wrappedTasks {
		wg.Add(1)
		go func() {
			err := task.Run()
			if err != nil {
				logger.Error(
					"task failed",
					slog.String("supervisor_task_error", err.Error()),
					slog.String("supervisor_task_name", taskName),
				)
			}
		}()
	}

	select {
	case <-startShutdownChan:
	case <-wgChan:
		return nil
	}

	for taskName, task := range wrappedTasks {
		wg.Add(1)
		go func() {
			err := task.Shutdown()
			if err != nil {
				logger.Error(
					"shutdown failed",
					slog.String("supervisor_task_error", err.Error()),
					slog.String("supervisor_task_name", taskName),
				)
			}
		}()
	}

	gracefulTimeout := time.After(gracefultShutdownDuration)
	select {
	case <-gracefulTimeout:
		cancel()
	case <-wgChan:
		return nil
	}

	shutdownTimeout := time.After(shutdownDuration)
	select {
	case <-shutdownTimeout:
		return fmt.Errorf("some task failed to shut down after %s", shutdownDuration+gracefultShutdownDuration)
	case <-wgChan:
		return nil
	}
}

func createShutdownChannel(n int) chan struct{} {
	shutdownChan := make(chan struct{}, n)

	go func() {
		osSignalChan := make(chan os.Signal, 1)
		signal.Notify(osSignalChan, syscall.SIGTERM)
		<-osSignalChan

		shutdownChan <- struct{}{}
	}()

	return shutdownChan
}

func wgAwait(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}
