# lockx

[![Go](https://img.shields.io/badge/go-1.26+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/lockx)](LICENSE)

`lockx` is a simple, thread-safe wrapper for values in Go using generics. It provides a clean API for concurrent access and modification of shared state.

## Features

- ✅ Generic support for any type
- ✅ Thread-safe operations (Load, Store, Swap, Update)
- ✅ `With` method for direct, protected access to the underlying value
- ✅ Simple and idiomatic API

## Installation

```bash
go get github.com/nduyhai/lockx
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/nduyhai/lockx"
)

func main() {
	// Create a new thread-safe value
	counter := lockx.New(0)

	// Update the value
	counter.Update(func(old int) int {
		return old + 1
	})

	// Load the value
	fmt.Println("Counter:", counter.Load()) // Output: Counter: 1

	// Direct access using With
	counter.With(func(v *int) {
		*v += 10
	})

	fmt.Println("Counter after With:", counter.Load()) // Output: Counter after With: 11
}
```

### Methods

- `New[T](initial T) *Value[T]`: Creates a new thread-safe value.
- `Load() T`: Atomically loads and returns the value.
- `Store(v T)`: Atomically stores a new value.
- `Swap(v T) T`: Atomically swaps the value and returns the old one.
- `Update(fn func(old T) T)`: Atomically updates the value using the provided function.
- `With(fn func(v *T))`: Executes a function with a pointer to the underlying value while holding the lock.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

