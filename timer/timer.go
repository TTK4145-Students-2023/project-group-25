package timer

import (
	"time"
)

var (
	startTimer = make(chan time.Duration)
	killTimer  = make(chan bool)
)

func TimerMain(timeout chan<- bool) {

	for {
		s := <-startTimer
		select {
		case <-time.After(time.Second * s):
			timeout <- true
		case <-killTimer:
		}
	}
}

func TimerStart(seconds time.Duration) {
	startTimer <- seconds
}

func TimerKill() {
	select {
	case killTimer <- true:
	default:
	}
}
