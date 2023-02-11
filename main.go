package main

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	server "Driver-go/servers"
	"math/rand"
	"time"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	FSM_orderExecuted := make(chan int)

	go elevio.PollFloorSensor()
	go elevio.PollObstructionSwitch()
	go elevio.PollStopButton()

	go server.InputServer()
	go server.DestinationServer()

	go fsm.FSM(FSM_orderExecuted)

	for {
		server.SetDestinationFloor(rand.Intn(3-0) + 0)
		time.Sleep(10 * time.Second)
	}
}
