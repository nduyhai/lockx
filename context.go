package lockx

import "context"

type Context[T any] struct {
	_    noCopy
	lock chan struct{}
	v    T
}

var _ ContextLocker[int] = (*Context[int])(nil)

func NewContext[T any](initial T) *Context[T] {
	x := &Context[T]{
		lock: make(chan struct{}, 1),
		v:    initial,
	}

	x.lock <- struct{}{}

	return x
}

func (x *Context[T]) acquire() {
	<-x.lock
}

func (x *Context[T]) release() {
	x.lock <- struct{}{}
}

func (x *Context[T]) acquireContext(ctx context.Context) error {
	select {
	case <-x.lock:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (x *Context[T]) With(fn func(v *T)) {
	x.acquire()
	defer x.release()

	fn(&x.v)
}

func (x *Context[T]) WithErr(fn func(v *T) error) error {
	x.acquire()
	defer x.release()

	return fn(&x.v)
}

func (x *Context[T]) WithContext(ctx context.Context, fn func(v *T)) error {
	if err := x.acquireContext(ctx); err != nil {
		return err
	}
	defer x.release()

	fn(&x.v)
	return nil
}

func (x *Context[T]) WithContextErr(ctx context.Context, fn func(v *T) error) error {
	if err := x.acquireContext(ctx); err != nil {
		return err
	}
	defer x.release()

	return fn(&x.v)
}

func (x *Context[T]) TryWith(fn func(v *T)) bool {
	select {
	case <-x.lock:
		defer x.release()

		fn(&x.v)
		return true

	default:
		return false
	}
}

func (x *Context[T]) TryWithErr(fn func(v *T) error) (bool, error) {
	select {
	case <-x.lock:
		defer x.release()

		return true, fn(&x.v)

	default:
		return false, nil
	}
}

func (x *Context[T]) Load() T {
	x.acquire()
	defer x.release()

	return x.v
}

func (x *Context[T]) Store(v T) {
	x.acquire()
	defer x.release()

	x.v = v
}

func (x *Context[T]) Swap(v T) T {
	x.acquire()
	defer x.release()

	old := x.v
	x.v = v
	return old
}

func (x *Context[T]) Update(fn func(old T) T) {
	x.acquire()
	defer x.release()

	x.v = fn(x.v)
}
