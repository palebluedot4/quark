package channel

type Queue[T any] struct {
	ch chan T
}

func NewQueue[T any](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("channel.NewQueue: capacity must be positive")
	}
	return &Queue[T]{ch: make(chan T, capacity)}
}

func (q *Queue[T]) Enqueue(v T) {
	q.ch <- v
}

func (q *Queue[T]) TryEnqueue(v T) bool {
	select {
	case q.ch <- v:
		return true
	default:
		return false
	}
}

func (q *Queue[T]) Dequeue() T {
	return <-q.ch
}

func (q *Queue[T]) TryDequeue() (T, bool) {
	select {
	case v := <-q.ch:
		return v, true
	default:
		var zero T
		return zero, false
	}
}
