package main

import (
	"Driver-go/elevio"
	"fsm"
)

type InputChanReadStruct struct {
	drv_floors chan int
	drv_obstr  chan bool
}

type InputChanWriteStruct struct {
	drv_floors chan int
	drv_obstr  chan bool
}

func inputServer(inputServerRead InputChanReadStruct, inputServerWrite InputChanWriteStruct) {
	var ELEV_floor int = -1
	var ELEV_obstr bool = false

	select {
	case newFloor := <-inputServerWrite.drv_floors:
		{
			ELEV_floor = newFloor
		}
	case newObstructionVal := <-inputServerWrite.drv_obstr:
		{
			ELEV_obstr = newObstructionVal
		}
	case inputServerRead.drv_floors <- ELEV_floor:
	case inputServerRead.drv_obstr <- ELEV_obstr:
	}
}

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var inputServerWrite InputChanWriteStruct
	var inputServerRead InputChanReadStruct

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_stop := make(chan bool)
	FSM_orderExecuted := make(chan int)
	FSM_destination := make(chan int)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(inputServerWrite.drv_floors)
	go elevio.PollObstructionSwitch(inputServerWrite.drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go inputServer(inputServerRead, inputServerWrite)

	go fsm.FSM(FSM_destination, inputServerRead.drv_floors, FSM_orderExecuted)

	for {

	}
}
