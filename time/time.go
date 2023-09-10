package _time

import (
	"log"
	"time"
)

func Now() time.Time {
	locTime := time.Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("时间本地化失败", err)
	} else {
		locTime = locTime.In(loc)
	}
	return locTime
}

func NowString() string {
	return Now().Format("2006-01-02 15:04:05")
}

func NowAdd(d time.Duration) time.Time {
	return Now().Add(d)
}

func NowAddString(d time.Duration) string {
	return NowAdd(d).Format("2006-01-02 15:04:05")
}
