package stack_test

import (
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/container/stack"
)

func TestStack(t *testing.T) {
	t.Parallel()
	t.Run("lifo order", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		for _, v := range []int{1, 2, 3} {
			s.Push(v)
		}
		var got []int
		for range 3 {
			v, ok := s.Pop()
			if !ok {
				t.Fatalf("Pop() = _, %v, want _, %v", ok, true)
			}
			got = append(got, v)
		}
		want := []int{3, 2, 1}
		if !slices.Equal(got, want) {
			t.Errorf("Pop() sequence = %v, want %v", got, want)
		}
	})

	t.Run("pop on empty", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		if got, ok := s.Pop(); ok || got != 0 {
			t.Errorf("Pop() = %v, %v, want %v, %v", got, ok, 0, false)
		}
	})

	t.Run("peek on empty", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		if got, ok := s.Peek(); ok || got != 0 {
			t.Errorf("Peek() = %v, %v, want %v, %v", got, ok, 0, false)
		}
	})

	t.Run("peek keeps the top", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		s.Push(1)
		s.Push(2)
		for range 2 {
			if got, ok := s.Peek(); !ok || got != 2 {
				t.Errorf("Peek() = %v, %v, want %v, %v", got, ok, 2, true)
			}
		}
		if got := s.Len(); got != 2 {
			t.Errorf("Len() = %v, want %v", got, 2)
		}
	})

	t.Run("len tracks pushes and pops", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %v, want %v", got, 0)
		}
		for i := range 3 {
			s.Push(i)
			if got := s.Len(); got != i+1 {
				t.Errorf("Len() = %v, want %v", got, i+1)
			}
		}
		for want := 2; want >= 0; want-- {
			s.Pop()
			if got := s.Len(); got != want {
				t.Errorf("Len() = %v, want %v", got, want)
			}
		}
	})

	t.Run("reuse after draining", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		s.Push(1)
		s.Pop()
		s.Push(2)
		if got, ok := s.Pop(); !ok || got != 2 {
			t.Errorf("Pop() = %v, %v, want %v, %v", got, ok, 2, true)
		}
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %v, want %v", got, 0)
		}
	})

	t.Run("all yields top to bottom", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		for _, v := range []int{1, 2, 3} {
			s.Push(v)
		}
		got, want := slices.Collect(s.All()), []int{3, 2, 1}
		if !slices.Equal(got, want) {
			t.Errorf("slices.Collect(All()) = %v, want %v", got, want)
		}
		if got := s.Len(); got != 3 {
			t.Errorf("Len() = %v, want %v", got, 3)
		}
	})

	t.Run("all is reusable", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		for _, v := range []int{1, 2, 3} {
			s.Push(v)
		}
		seq, want := s.All(), []int{3, 2, 1}
		for range 2 {
			if got := slices.Collect(seq); !slices.Equal(got, want) {
				t.Errorf("slices.Collect(All()) = %v, want %v", got, want)
			}
		}
	})

	t.Run("all stops early", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		for _, v := range []int{1, 2, 3} {
			s.Push(v)
		}
		var got []int
		for v := range s.All() {
			got = append(got, v)
			if len(got) == 2 {
				break
			}
		}
		want := []int{3, 2}
		if !slices.Equal(got, want) {
			t.Errorf("All() before break = %v, want %v", got, want)
		}
	})

	t.Run("all on empty", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		if got := slices.Collect(s.All()); len(got) != 0 {
			t.Errorf("slices.Collect(All()) = %v, want no elements", got)
		}
	})

	t.Run("all observes later pushes", func(t *testing.T) {
		t.Parallel()
		s := stack.NewStack[int]()
		s.Push(1)
		seq := s.All()
		s.Push(2)
		got, want := slices.Collect(seq), []int{2, 1}
		if !slices.Equal(got, want) {
			t.Errorf("slices.Collect(All()) = %v, want %v", got, want)
		}
	})

	t.Run("zero value is usable", func(t *testing.T) {
		t.Parallel()
		var s stack.Stack[string]
		if got := s.Len(); got != 0 {
			t.Errorf("Len() = %v, want %v", got, 0)
		}
		if got, ok := s.Peek(); ok || got != "" {
			t.Errorf("Peek() = %q, %v, want %q, %v", got, ok, "", false)
		}
		if got := slices.Collect(s.All()); len(got) != 0 {
			t.Errorf("slices.Collect(All()) = %q, want no elements", got)
		}
		s.Push("a")
		if got, ok := s.Pop(); !ok || got != "a" {
			t.Errorf("Pop() = %q, %v, want %q, %v", got, ok, "a", true)
		}
	})
}

//nolint:paralleltest
func TestStackPopReleasesElement(t *testing.T) {
	// This test must not run in parallel with other tests as it asserts on when
	// the popped element becomes unreachable, and their allocations delay the
	// cleanup queue it drains with runtime.GC.
	s := stack.NewStack[*byte]()
	var collected atomic.Bool
	buf := make([]byte, 1024)
	runtime.AddCleanup(&buf[0], func(b *atomic.Bool) { b.Store(true) }, &collected)
	s.Push(&buf[0])
	runtime.GC()
	if collected.Load() {
		t.Fatal("Push() released the pushed element before Pop(), want retained")
	}
	if _, ok := s.Pop(); !ok {
		t.Fatalf("Pop() = _, %v, want _, %v", ok, true)
	}
	if !collectedAfterGC(&collected) {
		t.Error("Pop() retained the popped element, want released")
	}
	runtime.KeepAlive(s)
}

func collectedAfterGC(b *atomic.Bool) bool {
	for range 100 {
		runtime.GC()
		if b.Load() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

func BenchmarkStackPush(b *testing.B) {
	const items = 1000
	for b.Loop() {
		s := stack.NewStack[int]()
		for i := range items {
			s.Push(i)
		}
	}
}

func BenchmarkStackPushPop(b *testing.B) {
	s := stack.NewStack[int]()
	for b.Loop() {
		s.Push(1)
		s.Pop()
	}
}
