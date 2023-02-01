package fsm

import (
	"Driver-go/elevio"
)

type ElevatorState int

const (
	STATE_Init           ElevatorState = 0
	STATE_AwaitingOrder                = 1
	STATE_ExecutingOrder               = 2
	STATE_DoorOpen                     = 3
)

func stateServer(stateRead chan<- ElevatorState, stateWrite <-chan ElevatorState, stateChanged chan<- int) {
	var state ElevatorState = STATE_Init
	for {
		select {
		case newState := <-stateWrite:
			if newState != state {
				state = newState
				stateChanged <- 1
			}
		case stateRead <- state:
		}
	}
}

func destinationServer(destinationRead chan<- int, destinationWrite <-chan int) {
	var destinationFloor int = -1
	for {
		select {
		case newDestinationFloor := <-destinationWrite:
			destinationFloor = newDestinationFloor
		case destinationRead <- destinationFloor:
		}
	}
}

func calculateMovingDirection(currentFloor, destinationFloor <-chan int) elevio.MotorDirection {
	floorDifference := <-destinationFloor - <-currentFloor
	if floorDifference > 0 {
		return elevio.MD_Up
	} else if floorDifference < 0 {
		return elevio.MD_Down
	} else {
		return elevio.MD_Stop
	}
}

func calculateState(stateRead <-chan ElevatorState) {
	state := <-stateRead
	switch state {
	case STATE_Init:
	case STATE_AwaitingOrder:
	case STATE_ExecutingOrder:
	case STATE_DoorOpen:
	}
}

func FSM(destinationWrite, currentFloor <-chan int, orderExecuted chan<- int) {
	var stateRead = make(chan ElevatorState)
	var stateWrite = make(chan ElevatorState)
	var stateChanged = make(chan int)
	var destinationRead = make(chan int)

	go stateServer(stateRead, stateWrite, stateChanged)
	go destinationServer(destinationRead, destinationWrite)
	go calculateState(stateRead)

	for {
		<-stateChanged
		state := <-stateRead
		switch state {
		case STATE_Init:
		case STATE_AwaitingOrder:
		case STATE_ExecutingOrder:
		case STATE_DoorOpen:
		}
	}
}
