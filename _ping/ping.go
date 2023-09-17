package _ping

import (
	"github.com/go-ping/ping"
	"time"
)

// 测试Ping的平均时间
func NewAverage(ip string, len int) (time.Duration, error) {
	var delayed time.Duration = 0
	var rTime time.Duration = 0
	var rNum int8
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return 0, err
	}
	pinger.Count = len
	pinger.Timeout = time.Second * time.Duration(len+3)
	pinger.OnRecv = func(pkt *ping.Packet) {
		rNum++
		rTime += pkt.Rtt
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		if rNum == 0 {
			return
		}
		delayed = rTime / time.Duration(rNum)
	}
	err = pinger.Run()
	if err != nil {
		return 0, err
	}
	return delayed, nil
}
