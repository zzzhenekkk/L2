package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		// Если нет каналов, возвращаем замкнутый канал
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		// Если один канал, возвращаем его напрямую
		return channels[0]
	default:
		// Объединяем два канала, и рекурсивно вызываем `or` для остальных
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or(append(channels[2:], orDone)...):
			}
		}()
		return orDone
	}
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("Done after %v\n", time.Since(start))
}
