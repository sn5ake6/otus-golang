package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Consume(taskCh <-chan Task, errorCh chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range taskCh {
		err := t()
		if err != nil {
			select {
			case errorCh <- struct{}{}:
			default:
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorCh := make(chan struct{}, m)
	defer close(errorCh)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	taskCh := make(chan Task)
	defer close(taskCh)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go Consume(taskCh, errorCh, &wg)
	}

	errorsCount := 0
	for _, t := range tasks {
		select {
		case <-errorCh:
			errorsCount++
			if m > 0 && errorsCount >= m {
				return ErrErrorsLimitExceeded
			}
		default:
		}

		taskCh <- t
	}

	return nil
}
