package waitgroup

import "sync"

func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Go(task)
	}
	wg.Wait()
}
