package _time

import (
	"sync"
	"time"
)

type TimerFunc func()

func NewTimer(d time.Duration, timerFunc TimerFunc) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			timerFunc()
		}
	}()
	wg.Wait()
}
