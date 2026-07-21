package iterator

import (
	"iter"
	"sync"
)

func RunAll(tasks iter.Seq[func()]) {
	var wg sync.WaitGroup
	for task := range tasks {
		wg.Go(task)
	}
	wg.Wait()
}

func RunAllManual(tasks iter.Seq[func()]) {
	var wg sync.WaitGroup
	tasks(func(task func()) bool {
		wg.Go(task)
		return true
	})
	wg.Wait()
}
