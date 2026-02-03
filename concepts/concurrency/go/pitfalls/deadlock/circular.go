package deadlock

import "sync"

func BadInconsistentLockOrder() {
	var (
		wg  sync.WaitGroup
		mu1 sync.Mutex
		mu2 sync.Mutex
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		defer mu1.Unlock()
		mu2.Lock()
		defer mu2.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu2.Lock()
		defer mu2.Unlock()
		mu1.Lock()
		defer mu1.Unlock()
	}()
	wg.Wait()
}

func GoodConsistentLockOrder() {
	var (
		wg  sync.WaitGroup
		mu1 sync.Mutex
		mu2 sync.Mutex
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		defer mu1.Unlock()
		mu2.Lock()
		defer mu2.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu1.Lock()
		defer mu1.Unlock()
		mu2.Lock()
		defer mu2.Unlock()
	}()
	wg.Wait()
}

func GoodConsistentLockOrderWithGo() {
	var (
		wg  sync.WaitGroup
		mu1 sync.Mutex
		mu2 sync.Mutex
	)
	wg.Go(func() {
		mu1.Lock()
		defer mu1.Unlock()
		mu2.Lock()
		defer mu2.Unlock()
	})
	wg.Go(func() {
		mu1.Lock()
		defer mu1.Unlock()
		mu2.Lock()
		defer mu2.Unlock()
	})
	wg.Wait()
}
