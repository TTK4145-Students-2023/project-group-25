package fsm

import (
	"Driver-go/elevio"
	server "Driver-go/servers"
	"Driver-go/timer"
)

type ElevatorState int

const (
	STATE_Init           ElevatorState = 0
	STATE_AwaitingOrder                = 1
	STATE_ExecutingOrder               = 2
	STATE_DoorOpen                     = 3
)

func calculateMovingDirection(currentFloor, destinationFloor int) elevio.MotorDirection {
	if floorDifference := destinationFloor - currentFloor; floorDifference > 0 {
		return elevio.MD_Up
	} else if floorDifference < 0 {
		return elevio.MD_Down
	} else {
		return elevio.MD_Stop
	}
}

func atDefinedFloor(currentFloor int) bool {
	return currentFloor != -1
}

func FSM(FSM_initCompleteChan, FSM_floorVisitedChan chan int) {
	var state ElevatorState = STATE_Init
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)

	for {
		currentFloor, destinationFloor := server.GetCurrentFloor(), server.GetDestinationFloor()

		switch state {
		case STATE_Init:
			if atDefinedFloor(currentFloor) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				state = STATE_AwaitingOrder
				FSM_initCompleteChan <- 1
			}

		case STATE_AwaitingOrder:
			if server.DestinationHasChanged() {
				server.DestinationChangeIsRecieved()
				newMovingDirection := calculateMovingDirection(currentFloor, destinationFloor)
				elevio.SetMotorDirection(newMovingDirection)
				state = STATE_ExecutingOrder
			}

		case STATE_ExecutingOrder:
			if currentFloor == destinationFloor {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				timer.SetTimer(3)
				state = STATE_DoorOpen
			}

		case STATE_DoorOpen:
			if !timer.TimeLeft() && !server.GetObstrVal() {
				elevio.SetDoorOpenLamp(false)
				state = STATE_AwaitingOrder
				FSM_floorVisitedChan <- currentFloor
			}
		}
	}
}
