package service

import (
	"runtime"
	"sync/atomic"
	"time"
)

type Pair struct {
	path string
	err  error
}

func worker(path string) Pair {
	res := Pair{
		path: path,
	}

	if err := Uncomment(path); err != nil {
		res.err = err
	}

	return res
}

func UncommentMany(paths ...string) <-chan Pair {
	cores := runtime.NumCPU()
	sem := make(chan struct{}, cores)

	var active int32

	done := make(chan struct{})

	go func(active *int32, done chan struct{}) {
		for atomic.LoadInt32(active) > 0 {
			time.Sleep(10 * time.Millisecond)
		}
		close(done)
	}(&active, done)

	tasks := len(paths)

	errCh := make(chan Pair, tasks)

	for i := range tasks {
		sem <- struct{}{}
		atomic.AddInt32(&active, 1)

		go func(path string, active *int32) {
			defer func(active *int32) {
				<-sem
				atomic.AddInt32(active, -1)
			}(active)

			p := worker(path)
			errCh <- p

		}(paths[i], &active)
	}

	<-done
	close(errCh)

	return errCh
}
