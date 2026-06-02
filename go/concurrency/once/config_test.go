package once_test

import (
	"sync"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/once"
)

func TestLoad(t *testing.T) {
	t.Setenv("APP_ADDR", "127.0.0.1:8080")
	t.Setenv("APP_TIMEOUT", "15s")
	tests := []struct {
		name string
		fn   func() (*once.Config, error)
	}{
		{name: "Load", fn: once.Load},
		{name: "LoadManual", fn: once.LoadManual},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const goroutines = 100
			cfgs := make([]*once.Config, goroutines)
			var wg sync.WaitGroup
			for i := range cfgs {
				wg.Go(func() {
					cfg, err := tt.fn()
					if err != nil {
						t.Errorf("%s() error = %v, want nil", tt.name, err)
					}
					cfgs[i] = cfg
				})
			}
			wg.Wait()
			if t.Failed() {
				t.FailNow()
			}
			for _, cfg := range cfgs {
				if cfg != cfgs[0] {
					t.Fatalf("%s() returned distinct instances %p and %p, want one", tt.name, cfg, cfgs[0])
				}
			}
			if cfgs[0].Addr != "127.0.0.1:8080" {
				t.Errorf("%s().Addr = %q, want %q", tt.name, cfgs[0].Addr, "127.0.0.1:8080")
			}
			if cfgs[0].Timeout != 15*time.Second {
				t.Errorf("%s().Timeout = %v, want %v", tt.name, cfgs[0].Timeout, 15*time.Second)
			}
		})
	}
}
