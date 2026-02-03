package waitgroup

import "sync"

func BadLateAdd() {
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
	}()
	wg.Wait()
}

func BadLateAddLoop() {
	var wg sync.WaitGroup
	for range 10 {
		go func() {
			wg.Add(1)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func GoodEarlyAddBulk() {
	var wg sync.WaitGroup
	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func GoodEarlyAddIterative() {
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func GoodSafeAddWithGo() {
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
		})
	}
	wg.Wait()
}

func BadWaitGroupReuse() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		wg.Done()
	}()
	wg.Wait()
}

func GoodWaitGroupReuse() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
	}()
	wg.Wait()
	wg.Add(1)
	go func() {
		wg.Done()
	}()
	wg.Wait()
}

func GoodWaitGroupReuseWithGo() {
	var wg sync.WaitGroup
	wg.Go(func() {
	})
	wg.Wait()
	wg.Go(func() {
	})
	wg.Wait()
}
