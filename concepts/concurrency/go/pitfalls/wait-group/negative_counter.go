package waitgroup

import "sync"

func BadNegativeCounter() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Done()
}

func BadNegativeCounterManual() {
	var wg sync.WaitGroup
	wg.Add(-1)
}

func GoodCounter(n int) {
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()
}
