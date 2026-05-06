package lockx

import "context"

type Locker[T any] interface {
	With(fn func(v *T))
	WithErr(fn func(v *T) error) error

	TryWith(fn func(v *T)) bool
	TryWithErr(fn func(v *T) error) (bool, error)

	Load() T
	Store(v T)
	Swap(v T) T
	Update(fn func(old T) T)
}

type RWLocker[T any] interface {
	Locker[T]

	View(fn func(v *T))
	ViewErr(fn func(v *T) error) error

	TryView(fn func(v *T)) bool
	TryViewErr(fn func(v *T) error) (bool, error)
}

// ContextLocker Real context-aware lock.
type ContextLocker[T any] interface {
	Locker[T]

	WithContext(ctx context.Context, fn func(v *T)) error
	WithContextErr(ctx context.Context, fn func(v *T) error) error
}
