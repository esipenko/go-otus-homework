package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRunNoErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	testCases := []struct {
		name         string
		tasksCount   int
		workersCount int
	}{
		{"workers less than tasks count without errors", 50, 5},
		{"workers more than tasks count without errors", 10, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks := make([]Task, 0, tc.tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			for i := 0; i < tc.tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				sumTime += taskSleep

				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			maxErrorsCount := 1

			start := time.Now()
			err := Run(tasks, tc.workersCount, maxErrorsCount)
			elapsedTime := time.Since(start)
			require.NoError(t, err)

			require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
			require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
		})
	}
}

func TestRunWithErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	testCases := []struct {
		name           string
		tasksCount     int
		workersCount   int
		maxErrorsCount int
	}{
		{"if were errors in first M tasks, than finished not more N+M tasks, len(tasks) < n", 50, 5, 5},
		{"if were errors in first M tasks, than finished not more N+M tasks, len(tasks) > n", 10, 20, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasksCount := 50
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32

			for i := 0; i < tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
		})
	}
}

func TestRunAdditional(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("zero tasks", func(t *testing.T) {
		var tasks []Task
		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("zero errors allowed", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(0), "no tasks started")
	})
}
