package main

import (
	"testing"
	"time"
)

func TestOrChannel(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("single channel", func(t *testing.T) {
		start := time.Now()
		<-or(sig(1 * time.Second))
		duration := time.Since(start)

		if duration < 1*time.Second {
			t.Errorf("Expected at least 1 second, but got %v", duration)
		}
	})

	t.Run("multiple channels", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(1*time.Second),
			sig(1*time.Hour),
			sig(1*time.Minute),
		)
		duration := time.Since(start)

		if duration < 1*time.Second || duration >= 2*time.Hour {
			t.Errorf("Expected around 1 second, but got %v", duration)
		}
	})

	t.Run("no channels", func(t *testing.T) {
		if or() != nil {
			t.Errorf("Expected nil for no channels, but got non-nil")
		}
	})

	t.Run("immediate close", func(t *testing.T) {
		c := make(chan interface{})
		close(c)

		start := time.Now()
		<-or(c)
		duration := time.Since(start)

		if duration > time.Millisecond {
			t.Errorf("Expected immediate return, but got %v", duration)
		}
	})
}
