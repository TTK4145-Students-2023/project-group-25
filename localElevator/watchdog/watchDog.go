package watchdog

import (
	dt "project/commonDataTypes"
	"time"
)

const (
	watchDogTime = 10 // Time [s] for watchdog
)

type ElevState string

func WatchDog(
	elevData chan dt.ElevDataJSON) {
	deadOrAliveCh := make(chan string)
	watchDogTimer := time.NewTimer(watchDogTime * time.Second)
	for {
		select {
		case <-elevData:
			deadOrAlive := <-deadOrAliveCh
			switch deadOrAlive {
			case "ALIVE":
				watchDogTimer.Reset(watchDogTime * time.Second)
			case "DEAD":
				// Connect to NTW
				watchDogTimer.Reset(watchDogTime * time.Second)
				deadOrAliveCh <- "ALIVE"
			}
		case <-watchDogTimer.C:
		}
	}
}
