package time

import (
	"time"
)

func NowString() (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	locTime := time.Now().In(loc)
	currentTime := locTime.Format("2006-01-02 15:04:05")
	return currentTime, nil
}
