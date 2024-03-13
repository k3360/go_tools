package _time

import (
	"sync"
	"time"
)

type TimerFunc func()

func NewTimer(d time.Duration, timerFunc TimerFunc) {
	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			timerFunc()
		}
	}()
	wg.Wait()
}
