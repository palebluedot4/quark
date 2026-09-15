package stack

import (
	"iter"
	"slices"
)

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	last := len(s.items) - 1
	v := s.items[last]
	clear(s.items[last:])
	s.items = s.items[:last]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range slices.Backward(s.items) {
			if !yield(v) {
				return
			}
		}
	}
}
