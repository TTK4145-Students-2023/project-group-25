package main

import (
	elevio "Module-go/localElevator/elev_driver"
	elevfsm "Module-go/localElevator/elev_fsm"
	"time"
)

var (
	floor_hallRequest      = make(chan [elevfsm.N_FLOORS][2]bool)
	floor_cabButtonEvent   = make(chan elevio.ButtonEvent)
	handler_ordersExecuted = make(chan []elevio.ButtonEvent)
	drv_buttons            = make(chan elevio.ButtonEvent)
	drv_floors             = make(chan int)
	drv_obstr              = make(chan bool)
	elev_data              = make(chan elevfsm.ElevatorData)
)

func orderHandler(buttonPress chan elevio.ButtonEvent, orderClear chan []elevio.ButtonEvent, hallOrder chan [elevfsm.N_FLOORS][2]bool, floor_cabButtonEvent chan elevio.ButtonEvent) {
	elev_hallOrder := [elevfsm.N_FLOORS][2]bool{}

	for {
		select {
		case buttonEvent := <-buttonPress:
			if buttonEvent.Button == elevio.BT_Cab {
				floor_cabButtonEvent <- buttonEvent
			} else {
				elev_hallOrder[buttonEvent.Floor][buttonEvent.Button] = true
				elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
				hallOrder <- elev_hallOrder
			}

		case clearEvent := <-orderClear:
			for i := 0; i < len(clearEvent); i++ {
				if clearEvent[i].Button != elevio.BT_Cab {
				} else {
					elev_hallOrder[clearEvent[i].Floor][clearEvent[i].Button] = false
					elevio.SetButtonLamp(clearEvent[i].Button, clearEvent[i].Floor, false)
				}
			}
			hallOrder <- elev_hallOrder
		}
	}
}

func main() {
	elevio.Init("localhost:15657", elevfsm.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollObstructionSwitch(drv_obstr)
	go orderHandler(drv_buttons, handler_ordersExecuted, floor_hallRequest, floor_cabButtonEvent)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(floor_hallRequest, floor_cabButtonEvent, drv_floors, drv_obstr, elev_data, handler_ordersExecuted)

	for {
		<-elev_data
	}
}
