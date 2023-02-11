package timer

import (
	"time"
)

var setAlarmTime, getAlarmTime chan time.Time

func TimerServer() {
	alarmTime := time.Now().Local()
	for {
		select {
		case newAlarmTime := <-setAlarmTime:
			alarmTime = newAlarmTime

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
