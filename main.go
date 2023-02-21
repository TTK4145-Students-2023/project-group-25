package main

import (
	elevio "Module-go/localElevator/elev_driver"
	elevfsm "Module-go/localElevator/elev_fsm"
	"time"
)

var (
	floor_request          = make(chan [elevfsm.N_FLOORS][elevfsm.N_BUTTONS]bool)
	handler_ordersExecuted = make(chan []elevio.ButtonEvent)
	drv_buttons            = make(chan elevio.ButtonEvent)
	drv_floors             = make(chan int)
	drv_obstr              = make(chan bool)
)

func orderHandler(buttonPress chan elevio.ButtonEvent, orderClear chan []elevio.ButtonEvent, order chan [elevfsm.N_FLOORS][elevfsm.N_BUTTONS]bool) {
	elevOrder := [elevfsm.N_FLOORS][elevfsm.N_BUTTONS]bool{}
	for {
		select {
		case buttonEvent := <-buttonPress:
			elevOrder[buttonEvent.Floor][buttonEvent.Button] = true
			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
			order <- elevOrder

		case clearEvent := <-orderClear:
			for i := 0; i < len(clearEvent); i++ {
				elevOrder[clearEvent[i].Floor][clearEvent[i].Button] = false
				elevio.SetButtonLamp(clearEvent[i].Button, clearEvent[i].Floor, false)
			}
			order <- elevOrder
		}
	}
}

func main() {
	elevio.Init("localhost:15657", elevfsm.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollObstructionSwitch(drv_obstr)
	go orderHandler(drv_buttons, handler_ordersExecuted, floor_request)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(floor_request, drv_floors, drv_obstr, handler_ordersExecuted)

	for {
		time.Sleep(time.Millisecond * 40)
	}
}
