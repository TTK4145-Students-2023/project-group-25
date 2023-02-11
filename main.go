package main

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"Driver-go/orderhandler"
	"math/rand"
)

func inputServer(inputServerRead, inputServerWrite orderhandler.InputServerChan) {
	inputData := orderhandler.InputServerData{
		DRV_floors: -1,
		DRV_obstr:  false,
		DRV_stop:   false}

	select {
	case newFloor := <-inputServerWrite.DRV_floors:
		{
			inputData.DRV_floors = newFloor
		}
	case newObstructionVal := <-inputServerWrite.DRV_obstr:
		{
			inputData.DRV_obstr = newObstructionVal
		}
	case newStopVal := <-inputServerWrite.DRV_stop:
		{
			inputData.DRV_stop = newStopVal
		}
	case inputServerRead.DRV_floors <- inputData.DRV_floors:
	case inputServerRead.DRV_obstr <- inputData.DRV_obstr:
	case inputServerRead.DRV_stop <- inputData.DRV_stop:
	}
}

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var inputServerRead, inputServerWrite orderhandler.InputServerChan
	var FSM_setDestination, FSM_getDestination, FSM_orderExecuted chan int

	go elevio.PollFloorSensor(inputServerWrite.DRV_floors)
	go elevio.PollObstructionSwitch(inputServerWrite.DRV_obstr)
	go elevio.PollStopButton(inputServerWrite.DRV_stop)

	go inputServer(inputServerRead, inputServerWrite)

	go fsm.FSM(FSM_setDestination, FSM_getDestination, FSM_orderExecuted, inputServerRead)

	for {
		FSM_setDestination <- (rand.Intn(3-0) + 0)
		<-FSM_orderExecuted
	}
}
