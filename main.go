package main

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	server "Driver-go/servers"
	"Driver-go/timer"
	"math/rand"
	"time"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	FSM_floorVisitedChan := make(chan int)
	FSM_initCompleteChan := make(chan int)

	go elevio.PollFloorSensor(server.SetCurrentFloorChan)
	go elevio.PollObstructionSwitch(server.SetObstrValChan)
	go elevio.PollStopButton(server.SetStopValChan)

	go timer.Timer()

	go server.InputServer()
	go server.DestinationServer()

	time.Sleep(time.Millisecond * 40)

	go fsm.FSM(FSM_initCompleteChan, FSM_floorVisitedChan)
	<-FSM_initCompleteChan
	for {
		server.SetDestinationFloor(rand.Intn(4))
	}
}
