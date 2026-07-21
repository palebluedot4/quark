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

func RunUntilError(tasks iter.Seq[func() error]) error {
	for task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}
	return nil
}

func RunUntilErrorManual(tasks iter.Seq[func() error]) error {
	var err error
	tasks(func(task func() error) bool {
		if e := task(); e != nil {
			err = e
			return false
		}
		return true
	})
	return err
}
