package channel_test

import (
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/channel"
)

//nolint:paralleltest
func TestBlockForever(t *testing.T) {
	// This test must not run in parallel with other tests as it infers that f
	// still blocks from a 50ms timeout, and their load makes that window
	// unreliable.
	tests := []struct {
		name string
		f    func()
	}{
		{name: "BlockOnEmptySelect", f: channel.BlockOnEmptySelect},
		{name: "BlockOnNilChannel", f: channel.BlockOnNilChannel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				tt.f()
				close(done)
			}()
			select {
			case <-done:
				t.Errorf("%s() returned, want it to block forever", tt.name)
			case <-time.After(50 * time.Millisecond):
			}
		})
	}
}
