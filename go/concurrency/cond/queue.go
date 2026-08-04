package cond

import "sync"

type Queue[T any] struct {
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond

	buf   []T
	head  int
	tail  int
	count int
}

func NewQueue[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("cond.NewQueue: capacity must be positive")
	}
	q := &Queue[T]{buf: make([]T, capacity)}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.count == len(q.buf) {
		q.notFull.Wait()
	}
	q.buf[q.tail] = item
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
	q.notEmpty.Signal()
}

func (q *Queue[T]) Dequeue() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.count == 0 {
		q.notEmpty.Wait()
	}
	item := q.buf[q.head]
	var zero T
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	q.notFull.Signal()
	return item
}
