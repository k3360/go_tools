package _time

import (
	"time"
)

// 当前时间字符串
func String(format string) string {
	return Now().Format(time.DateOnly)
}
