package atomicx

import (
	"sync"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	t.Run("reports parse failure", func(t *testing.T) {
		cfg.Store(nil)
		t.Setenv("APP_TIMEOUT", "half an hour")
		got, err := Load()
		if err == nil {
			t.Errorf("Load() = %v, want error", got)
		}
	})

	t.Run("parses once under concurrency", func(t *testing.T) {
		const (
			goroutines = 100
			addr       = "127.0.0.1:9000"
			timeout    = 15 * time.Second
		)
		cfg.Store(nil)
		t.Setenv("APP_ADDR", addr)
		t.Setenv("APP_TIMEOUT", timeout.String())
		cfgs := make([]*Config, goroutines)
		var wg sync.WaitGroup
		for i := range cfgs {
			wg.Go(func() {
				c, err := Load()
				if err != nil {
					t.Errorf("Load() error = %v, want nil", err)
					return
				}
				cfgs[i] = c
			})
		}
		wg.Wait()
		if t.Failed() {
			t.FailNow()
		}
		for _, c := range cfgs {
			if c != cfgs[0] {
				t.Fatalf("Load() returned distinct instances %p and %p, want one", c, cfgs[0])
			}
		}
		want := Config{Addr: addr, Timeout: timeout}
		if *cfgs[0] != want {
			t.Errorf("Load() = %+v, want %+v", *cfgs[0], want)
		}
	})

	t.Run("ignores environment once parsed", func(t *testing.T) {
		cfg.Store(nil)
		want, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		t.Setenv("APP_ADDR", "0.0.0.0:1")
		got, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		if got != want {
			t.Errorf("Load() = %p, want %p", got, want)
		}
	})

	t.Run("replaces config on reload", func(t *testing.T) {
		const (
			addr    = "10.0.0.1:7000"
			timeout = time.Minute
		)
		cfg.Store(nil)
		prev, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		t.Setenv("APP_ADDR", addr)
		t.Setenv("APP_TIMEOUT", timeout.String())
		if err = Reload(); err != nil {
			t.Fatalf("Reload() error = %v, want nil", err)
		}
		got, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		if got == prev {
			t.Fatalf("Load() = %p, want a new instance", got)
		}
		want := Config{Addr: addr, Timeout: timeout}
		if *got != want {
			t.Errorf("Load() = %+v, want %+v", *got, want)
		}
	})

	t.Run("keeps config when reload fails", func(t *testing.T) {
		cfg.Store(nil)
		want, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		t.Setenv("APP_TIMEOUT", "half an hour")

		if err = Reload(); err == nil {
			t.Fatal("Reload() error = nil, want error")
		}
		got, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v, want nil", err)
		}
		if got != want {
			t.Errorf("Load() = %p, want %p", got, want)
		}
	})

	t.Run("reads complete snapshots during reload", func(t *testing.T) {
		const (
			readers = 8
			reads   = 1000
			reloads = 100
		)
		want := [2]Config{
			{Addr: "192.0.2.1:8000", Timeout: time.Second},
			{Addr: "192.0.2.2:9000", Timeout: 2 * time.Second},
		}
		cfg.Store(nil)
		t.Setenv("APP_ADDR", want[0].Addr)
		t.Setenv("APP_TIMEOUT", want[0].Timeout.String())
		if err := Reload(); err != nil {
			t.Fatalf("Reload() error = %v, want nil", err)
		}
		var wg sync.WaitGroup
		for range readers {
			wg.Go(func() {
				for range reads {
					c, err := Load()
					if err != nil {
						t.Errorf("Load() error = %v, want nil", err)
						return
					}
					if *c != want[0] && *c != want[1] {
						t.Errorf("Load() = %+v, want %+v or %+v", *c, want[0], want[1])
						return
					}
				}
			})
		}
		for i := range reloads {
			t.Setenv("APP_ADDR", want[i%2].Addr)
			t.Setenv("APP_TIMEOUT", want[i%2].Timeout.String())
			if err := Reload(); err != nil {
				t.Errorf("Reload() error = %v, want nil", err)
				break
			}
		}
		wg.Wait()
	})
}
