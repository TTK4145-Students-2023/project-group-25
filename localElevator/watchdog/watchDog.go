package watchdog

import (
	dt "project/commonDataTypes"
	"time"
)

type WD_Role string

const (
	watchDogTime         = 10 // Time [s] for watchdog
	WD_ALIVE     WD_Role = "ALIVE"
	WD_DEAD      WD_Role = "DEAD"
)

func WatchDog(
	elevData chan dt.ElevDataJSON,
	peerTxEnable chan bool) {
	watchDogTimer := time.NewTimer(watchDogTime * time.Second)
	watchDogRole := WD_ALIVE
	elevState := dt.ElevDataJSON{}
	for {
		select {
		case <-elevData:
			watchDogTimer.Reset(watchDogTime * time.Second)
			switch watchDogRole {
			case "ALIVE":
			case "DEAD":
				peerTxEnable <- true // Connect to network
				watchDogRole = "ALIVE"
			}
		case <-watchDogTimer.C:
			switch elevState.Behavior {
			case "idle":
				watchDogTimer.Reset(watchDogTime * time.Second)
			case "moving", "doorOpen":
				peerTxEnable <- false // Disconnect from network
				watchDogRole = "DEAD"
			}
		}
	}
}
