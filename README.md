# lockx

[![Go](https://img.shields.io/badge/go-1.18+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/lockx)](LICENSE)

`lockx` is a simple, thread-safe wrapper for values in Go using generics. It provides a clean API for concurrent access and modification of shared state.

## Features

- ✅ Generic support for any type
- ✅ Thread-safe operations (Load, Store, Swap, Update)
- ✅ `With` method for direct, protected access to the underlying value
- ✅ Context-aware locking
- ✅ RWLock support with fast `View` method
- ✅ Simple and idiomatic API

## Installation

```bash
go get github.com/nduyhai/lockx
```

## Usage

### Mutex

```go
counter := lockx.NewMutex(0)

counter.With(func(v *int) {
	*v++
})
```

### Context-aware

```go
counter := lockx.NewContext(0)

ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

err := counter.WithContext(ctx, func(v *int) {
	*v++
})
```

### RWMutex

```go
state := lockx.NewRW(MyStruct{})

// Read-only view
state.View(func(v *MyStruct) {
    fmt.Println(v.Name)
})

// Update
state.With(func(v *MyStruct) {
    v.Name = "new name"
})
```

#### Important warning for RW.View

`View(fn func(v *T))` is fast because it avoids copying large structs.

But users **must not mutate** inside `View`, because it only holds a read lock.

Example safe usage:

```go
state.View(func(v *State) {
	fmt.Println(v.Name)
})
```

Bad usage:

```go
state.View(func(v *State) {
	v.Name = "changed" // don't do this
})
```

## go.mod

Use at least Go 1.18 because `TryLock` / `TryRLock` require it.

```go
module github.com/yourname/lockx

go 1.22
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

