package timer

import (
	"time"
)

var setAlarmTime, getAlarmTime chan time.Time

func TimerServer() {
	alarmTime := time.Now().Local()

	select {
	case newAlarmTime := <-setAlarmTime:
		alarmTime = newAlarmTime

	case getAlarmTime <- alarmTime:
	}
}

func SetTimer(seconds int) {
	setAlarmTime <- time.Now().Local().Add(time.Second * time.Duration(seconds))
}

func timeLeft() bool {

	difference := (<-getAlarmTime).Sub(time.Now().Local())
	if difference < 0 {
		return false
	}
	return true
}
