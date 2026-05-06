package lockx

import (
	"context"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	counter := NewMutex(0)

	counter.With(func(v *int) {
		*v++
	})

	if val := counter.Load(); val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	counter.Update(func(old int) int {
		return old + 10
	})

	if val := counter.Load(); val != 11 {
		t.Errorf("expected 11, got %d", val)
	}

	if ok := counter.TryWith(func(v *int) {
		*v++
	}); !ok {
		t.Errorf("expected TryWith to succeed")
	}

	if val := counter.Load(); val != 12 {
		t.Errorf("expected 12, got %d", val)
	}
}

func TestContext(t *testing.T) {
	counter := NewContext(0)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := counter.WithContext(ctx, func(v *int) {
		*v++
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if val := counter.Load(); val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	// Test timeout
	counter.With(func(v *int) {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel2()

		err := counter.WithContext(ctx2, func(v *int) {
			*v++
		})
		if err == nil {
			t.Error("expected timeout error, got nil")
		}
	})
}

func TestRW(t *testing.T) {
	state := NewRW(0)

	state.With(func(v *int) {
		*v = 42
	})

	state.View(func(v *int) {
		if *v != 42 {
			t.Errorf("expected 42, got %d", *v)
		}
	})

	if val := state.Load(); val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	state.Store(100)
	if val := state.Load(); val != 100 {
		t.Errorf("expected 100, got %d", val)
	}
}
