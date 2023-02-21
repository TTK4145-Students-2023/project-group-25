package main

import (
	"Driver-go/elevio"
	"Driver-go/fsm"
	"time"
)

var (
	floor_request     = make(chan [fsm.N_FLOORS][fsm.N_BUTTONS]bool)
	drv_ordersExecuted = make(chan []elevio.ButtonEvent)
	drv_buttons       = make(chan elevio.ButtonEvent)
	drv_floors        = make(chan int)
	drv_obstr         = make(chan bool)
)

func orderHandler(buttonPress chan elevio.ButtonEvent, orderClear chan []elevio.ButtonEvent, order chan [fsm.N_FLOORS][fsm.N_BUTTONS]bool) {
	elevOrder := [fsm.N_FLOORS][fsm.N_BUTTONS]bool{}
	for {
		select {
		case buttonEvent := <-buttonPress:
			elevOrder[buttonEvent.Floor][buttonEvent.Button] = true
			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
			order <- elevOrder

		case clearEvent := <-orderClear:
			for i := 0; i<len(clearEvent); i++{
				elevOrder[clearEvent[i].Floor][clearEvent[i].Button] = false
				elevio.SetButtonLamp(clearEvent[i].Button, clearEvent[i].Floor, false)
			}
			order <- elevOrder
		}
	}
}

func main() {
	elevio.Init("localhost:15657", fsm.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollObstructionSwitch(drv_obstr)
	go orderHandler(drv_buttons, drv_ordersExecuted, floor_request)

	time.Sleep(time.Millisecond * 40)

	go fsm.FSM(floor_request, drv_floors, drv_obstr, drv_ordersExecuted)

	for {
		time.Sleep(time.Millisecond * 40)
	}
}
