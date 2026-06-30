package semaphorex

import "sync"

func RunAll(n int, tasks []func()) {
	if n < 1 {
		panic("semaphorex.RunAll: limit must be positive")
	}
	sem := make(chan struct{}, n)
	var wg sync.WaitGroup
	for _, task := range tasks {
		sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-sem }()
			task()
		})
	}
	wg.Wait()
}
