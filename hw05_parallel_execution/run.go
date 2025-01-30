package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		errorCount     int
		currentTaskIdx int
		wg             sync.WaitGroup
		mx             sync.Mutex
	)

	runGoRoutine := func() {
		defer wg.Done()

		for {
			mx.Lock()

			if errorCount >= m || currentTaskIdx >= len(tasks) {
				mx.Unlock()
				return
			}

			task := tasks[currentTaskIdx]
			currentTaskIdx++

			mx.Unlock()

			if err := task(); err != nil {
				mx.Lock()
				errorCount++
				mx.Unlock()
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go runGoRoutine()
	}

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
