package waitgroup

import "sync"

func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Go(task)
	}
	wg.Wait()
}

func RunAllManual(tasks []func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
}
