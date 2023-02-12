package timer

import (
	"time"
)

var (
	setAlarmTime = make(chan time.Time)
	getAlarmTime = make(chan time.Time)
)

func Timer() {
	alarmTime := time.Now().Local()
	for {
		select {
		case alarmTime = <-setAlarmTime:
		case getAlarmTime <- alarmTime:
		}
	}
}

func SetTimer(seconds int) {
	setAlarmTime <- time.Now().Local().Add(time.Second * time.Duration(seconds))
}

func TimeLeft() bool {
	difference := (<-getAlarmTime).Sub(time.Now().Local())
	return difference >= 0
}
