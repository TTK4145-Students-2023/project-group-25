package fsm

import (
	"Driver-go/elevio"
)

var getFloor = make(chan Floor)

type ElevatorState int

const (
	STATE_Init           ElevatorState = 0
	STATE_AwaitingOrder                = 1
	STATE_ExecutingOrder               = 2
	STATE_DoorOpen                     = 3
)

type Floor struct {
	current     int
	destination int
}

func floorServer(setFloor <-chan Floor, getFloor chan<- Floor) {
	var floor Floor
	for {
		select {
		case newFloor := <-setFloor:
			floor = newFloor
		case getFloor <- floor:
		}
	}
}

func calculateMovingDirection(currentFloor, destinationFloor int) elevio.MotorDirection {
	floorDifference := destinationFloor - currentFloor
	if floorDifference > 0 {
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

func transitionToState(state ElevatorState) {
	switch state {
	case STATE_Init:
	case STATE_AwaitingOrder:
	case STATE_ExecutingOrder:
		{
			elevio.SetMotorDirection(calculateMovingDirection((<-getFloor).current, (<-getFloor).destination))
		}

	case STATE_DoorOpen:
	}
}

func FSM(setFloor <-chan Floor, getOrderExecuted chan<- int) {

	go floorServer(setFloor, getFloor)

	var state ElevatorState = STATE_Init
	elevio.SetMotorDirection(elevio.MD_Down)

	for {
		switch state {
		case STATE_Init:
			if atDefinedFloor((<-getFloor).current) {
				transitionToState(STATE_AwaitingOrder)
			}
		case STATE_AwaitingOrder:

		case STATE_ExecutingOrder:

		case STATE_DoorOpen:
		}
	}
}
