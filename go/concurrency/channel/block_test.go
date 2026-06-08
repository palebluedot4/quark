package channel_test

import (
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/channel"
)

func TestBlockForever(t *testing.T) {
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
