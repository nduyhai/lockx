package lockx

import "sync"

type RW[T any] struct {
	_  noCopy
	mu sync.RWMutex
	v  T
}

var _ RWLocker[int] = (*RW[int])(nil)

func NewRW[T any](initial T) *RW[T] {
	return &RW[T]{v: initial}
}

func (x *RW[T]) With(fn func(v *T)) {
	x.mu.Lock()
	defer x.mu.Unlock()

	fn(&x.v)
}

func (x *RW[T]) WithErr(fn func(v *T) error) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	return fn(&x.v)
}

func (x *RW[T]) TryWith(fn func(v *T)) bool {
	if !x.mu.TryLock() {
		return false
	}
	defer x.mu.Unlock()

	fn(&x.v)
	return true
}

func (x *RW[T]) TryWithErr(fn func(v *T) error) (bool, error) {
	if !x.mu.TryLock() {
		return false, nil
	}
	defer x.mu.Unlock()

	return true, fn(&x.v)
}

func (x *RW[T]) View(fn func(v *T)) {
	x.mu.RLock()
	defer x.mu.RUnlock()

	fn(&x.v)
}

func (x *RW[T]) ViewErr(fn func(v *T) error) error {
	x.mu.RLock()
	defer x.mu.RUnlock()

	return fn(&x.v)
}

func (x *RW[T]) TryView(fn func(v *T)) bool {
	if !x.mu.TryRLock() {
		return false
	}
	defer x.mu.RUnlock()

	fn(&x.v)
	return true
}

func (x *RW[T]) TryViewErr(fn func(v *T) error) (bool, error) {
	if !x.mu.TryRLock() {
		return false, nil
	}
	defer x.mu.RUnlock()

	return true, fn(&x.v)
}

func (x *RW[T]) Load() T {
	x.mu.RLock()
	defer x.mu.RUnlock()

	return x.v
}

func (x *RW[T]) Store(v T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = v
}

func (x *RW[T]) Swap(v T) T {
	x.mu.Lock()
	defer x.mu.Unlock()

	old := x.v
	x.v = v
	return old
}

func (x *RW[T]) Update(fn func(old T) T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = fn(x.v)
}
