package channel

import "sync"

func RunAll(n int, tasks []func()) {
	if n < 1 {
		panic("channel.RunAll: limit must be positive")
	}
	ch := make(chan func())
	var wg sync.WaitGroup
	for range min(n, len(tasks)) {
		wg.Go(func() {
			for task := range ch {
				task()
			}
		})
	}
	for _, task := range tasks {
		ch <- task
	}
	close(ch)
	wg.Wait()
}
