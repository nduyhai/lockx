package lockx

import "sync"

type Mutex[T any] struct {
	_  noCopy
	mu sync.Mutex
	v  T
}

var _ Locker[int] = (*Mutex[int])(nil)

func NewMutex[T any](initial T) *Mutex[T] {
	return &Mutex[T]{v: initial}
}

func (x *Mutex[T]) With(fn func(v *T)) {
	x.mu.Lock()
	defer x.mu.Unlock()

	fn(&x.v)
}

func (x *Mutex[T]) WithErr(fn func(v *T) error) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	return fn(&x.v)
}

func (x *Mutex[T]) TryWith(fn func(v *T)) bool {
	if !x.mu.TryLock() {
		return false
	}
	defer x.mu.Unlock()

	fn(&x.v)
	return true
}

func (x *Mutex[T]) TryWithErr(fn func(v *T) error) (bool, error) {
	if !x.mu.TryLock() {
		return false, nil
	}
	defer x.mu.Unlock()

	return true, fn(&x.v)
}

func (x *Mutex[T]) Load() T {
	x.mu.Lock()
	defer x.mu.Unlock()

	return x.v
}

func (x *Mutex[T]) Store(v T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = v
}

func (x *Mutex[T]) Swap(v T) T {
	x.mu.Lock()
	defer x.mu.Unlock()

	old := x.v
	x.v = v
	return old
}

func (x *Mutex[T]) Update(fn func(old T) T) {
	x.mu.Lock()
	defer x.mu.Unlock()

	x.v = fn(x.v)
}
