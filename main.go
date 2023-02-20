package main

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"time"
)

var (
	floor_request     = make(chan [fsm.N_FLOORS][fsm.N_BUTTONS]bool)
	drv_orderExecuted = make(chan elevio.ButtonEvent)
	drv_buttons       = make(chan elevio.ButtonEvent)
	drv_floors        = make(chan int)
	drv_obstr         = make(chan bool)
)

func orderHandler(buttonPress chan elevio.ButtonEvent, orderClear chan elevio.ButtonEvent, order chan [fsm.N_FLOORS][fsm.N_BUTTONS]bool) {
	elevOrder :=  [fsm.N_FLOORS][fsm.N_BUTTONS]bool{}

	select{
	case buttonEvent := <- buttonPress:
		elevOrder[buttonEvent.Floor][buttonEvent.Button] = true
		elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
		order <- elevOrder
	
	case clearEvent := <- orderClear:
		elevOrder[clearEvent.Floor][elevio.BT_Cab] = false
		elevOrder[clearEvent.Floor][elevio.BT_HallDown] = false
		elevOrder[clearEvent.Floor][elevio.BT_HallUp] = false

		elevio.SetButtonLamp(elevio.BT_Cab, clearEvent.Floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, clearEvent.Floor, false)
		elevio.SetButtonLamp(elevio.BT_HallUp, clearEvent.Floor, false)
		order <- elevOrder
	}
}


func main() {
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollObstructionSwitch(drv_obstr)
	go orderHandler(drv_buttons, drv_orderExecuted, floor_request)

	time.Sleep(time.Millisecond * 40)

	go fsm.FSM(floor_request, drv_floors, drv_obstr, drv_orderExecuted)

	for {
		time.Sleep(time.Millisecond * 40)
	}
}
