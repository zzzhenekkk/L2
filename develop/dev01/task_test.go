package main

import (
	"testing"
	"time"

	"github.com/beevik/ntp"
)

func TestNTPTime(t *testing.T) {
	_, err := ntp.Time("pool.ntp.org")
	if err != nil {
		t.Errorf("Failed to fetch NTP time: %v", err)
	}
}

func TestCurrentTime(t *testing.T) {
	currentTime := time.Now()
	if currentTime.IsZero() {
		t.Error("Current time is zero")
	}
}
