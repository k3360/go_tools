package _time

import (
	"log"
	"time"
)

// 当前时间
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

// 当前时间字符串
func NowString() string {
	return Now().Format(time.DateTime)
}

// 当前时间 + 指定时间
func NowAdd(d time.Duration) time.Time {
	return Now().Add(d)
}

// 当前时间 + 指定时间 的字符串
func NowAddString(d time.Duration) string {
	return NowAdd(d).Format(time.DateTime)
}

// 时间戳转字符串
func TimestampToString(millisecond int64) string {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}
	localTime := time.Unix(0, millisecond*int64(time.Millisecond)).In(location)
	return localTime.Format(time.DateTime)
}

// 字符串转时间戳
func StringToTimestamp(datetime string) (int64, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	// 使用 time.ParseInLocation 解析该字符串
	t, err := time.ParseInLocation(time.DateTime, datetime, loc)
	if err != nil {
		return 0, err
	}
	// 获取毫秒时间戳
	milliseconds := t.UnixNano() / int64(time.Millisecond)
	return milliseconds, nil
}
