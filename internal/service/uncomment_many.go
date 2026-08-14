package service

import (
	"runtime"
	"sync/atomic"
	"time"
)

// Pair is the result of processing a single file: its path and, if the
// operation failed, the error that occurred.
type Pair struct {
	Path string
	Err  error
}

func worker(path string) Pair {
	res := Pair{
		Path: path,
	}

	if err := Uncomment(path); err != nil {
		res.Err = err
	}

	return res
}

// UncommentMany strips comments from the given files concurrently.
//
// It runs a worker pool bounded by the number of CPU cores, calls Uncomment
// on every path and streams the results through the returned channel as a
// sequence of Pairs. The channel is closed once all files are processed.
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
