package lockx

import "sync"

type Value[T any] struct {
	mu sync.Mutex
	v  T
}

func New[T any](initial T) *Value[T] {
	return &Value[T]{v: initial}
}

func (x *Value[T]) With(fn func(v *T)) {
	x.mu.Lock()
	defer x.mu.Unlock()

	fn(&x.v)
}

func (x *Value[T]) Load() T {
	x.mu.Lock()
	defer x.mu.Unlock()

	return x.v
}

func (x *Value[T]) Store(v T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = v
}

func (x *Value[T]) Swap(v T) T {
	x.mu.Lock()
	defer x.mu.Unlock()

	old := x.v
	x.v = v

	return old
}

func (x *Value[T]) Update(fn func(old T) T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = fn(x.v)
}
