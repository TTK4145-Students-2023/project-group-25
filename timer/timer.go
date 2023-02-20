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
		alarm := time.After(time.Second * s)
		select {
		case <-alarm:
			timeout <- true
		case <-killTimer:
		}
	}
}

func TimerStart(seconds time.Duration) {
	startTimer <- seconds
}

func TimerKill() {
	killTimer <- true
}
