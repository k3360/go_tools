package _time

import (
	"time"
)

type TimerFunc func()

func NewTimer(d time.Duration, timerFunc TimerFunc) {
	go func() {
		ticker := time.NewTicker(d)
		for range ticker.C {
			timerFunc()
		}
	}()
}
