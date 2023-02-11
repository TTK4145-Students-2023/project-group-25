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

func destinationServer(setDestination <-chan int, getDestination chan<- int) {
	var destination int
	for {
		select {
		case newDestination := <-setDestination:
			destination = newDestination
		case getDestination <- destination:
		}
	}
}

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

func transitionToState(state ElevatorState) {
	switch state {
	case STATE_Init:
		elevio.SetMotorDirection(elevio.MD_Down)
	case STATE_AwaitingOrder:
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(false)
	case STATE_ExecutingOrder:
		elevio.SetMotorDirection(calculateMovingDirection(server.GetCurrentFloor(), server.GetDestinationFloor()))
	case STATE_DoorOpen:
		elevio.SetDoorOpenLamp(true)
		elevio.SetMotorDirection(elevio.MD_Stop)
		timer.SetTimer(3)
	}
}

func FSM(FSM_orderExecuted chan int) {

	go timer.TimerServer()

	var state ElevatorState = STATE_Init
	transitionToState(STATE_Init)

	for {
		currentFloor, destinationFloor := server.GetCurrentFloor(), server.GetDestinationFloor()

		switch state {
		case STATE_Init:
			if atDefinedFloor(currentFloor) {
				transitionToState(STATE_AwaitingOrder)
			}

		case STATE_AwaitingOrder:
			if atDefinedFloor(currentFloor) && (currentFloor != destinationFloor) {
				transitionToState(STATE_ExecutingOrder)
			}

		case STATE_ExecutingOrder:
			if currentFloor == destinationFloor {
				transitionToState(STATE_DoorOpen)
			}

		case STATE_DoorOpen:
			if !timer.TimeLeft() {
				transitionToState(STATE_AwaitingOrder)
				FSM_orderExecuted <- currentFloor
			}
		}
	}
}
