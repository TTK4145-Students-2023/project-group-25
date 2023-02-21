package fsm

import (
	"Driver-go/elevio"
	"Driver-go/timer"
	"time"
)

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ClearRequestVariant int

const (
	// Assume everyone waiting for the elevator gets on the elevator, even if
	// they will be traveling in the "wrong" direction for a while
	CV_All ClearRequestVariant = 0

	// Assume that only those that want to travel in the current direction
	// enter the elevator, and keep waiting outside otherwise
	CV_InDirn ClearRequestVariant = 1
)

type Config struct {
	clearRequestVariant ClearRequestVariant
	doorOpenDuration_s  time.Duration
}

type Elevator struct {
	floor     int
	dirn      elevio.MotorDirection
	requests  [N_FLOORS][N_BUTTONS]bool
	behaviour ElevatorBehaviour
	config    Config
}

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_Moving   ElevatorBehaviour = 1
	EB_DoorOpen ElevatorBehaviour = 2
)

func uninitializedElevator() Elevator {
	return Elevator{
		floor:     -1,
		dirn:      elevio.MD_Stop,
		requests:  [N_FLOORS][N_BUTTONS]bool{},
		behaviour: EB_Idle,
		config:    Config{clearRequestVariant: CV_InDirn, doorOpenDuration_s: 3},
	}
}

func FSM(
	floor_requests <-chan [N_FLOORS][N_BUTTONS]bool,
	drv_floors <-chan int,
	drv_obstr <-chan bool,
	drv_orderExecuted chan<- []elevio.ButtonEvent) {

	e := uninitializedElevator()
	obstr := false
	timeout := make(chan bool)

	go timer.TimerMain(timeout)

	select {
	case e.floor = <-drv_floors:
	default:
	}

	if e.floor == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
		e.dirn = elevio.MD_Down
		e.behaviour = EB_Moving

		e.floor = <-drv_floors

		elevio.SetMotorDirection(elevio.MD_Stop)
		e.dirn = elevio.MD_Stop
		e.behaviour = EB_Idle
	}

	for {
		select {
		case e.floor = <-drv_floors:
			elevio.SetFloorIndicator(e.floor)
			switch e.behaviour {
			case EB_Idle:
			case EB_DoorOpen:
			case EB_Moving:
				if requests_shouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					if !obstr { timer.TimerStart(e.config.doorOpenDuration_s) }
					e.behaviour = EB_DoorOpen
				}
			}

		case <-timeout:
			switch e.behaviour {
			case EB_Moving:
			case EB_Idle:
			case EB_DoorOpen:
				ordersExecuted := requests_calculateOrdersToBeCleared(e)
				e = requests_clearLocalRequest(e, ordersExecuted)
				drv_orderExecuted <- ordersExecuted

				dirnBehaviourPair := requests_chooseDirection(e)
				e.dirn = dirnBehaviourPair.dirn
				e.behaviour = dirnBehaviourPair.behaviour

				switch e.behaviour {
				case EB_DoorOpen:
					timer.TimerStart(e.config.doorOpenDuration_s)

				case EB_Moving:
					fallthrough
				case EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(e.dirn)
				}
			}
		case obstr = <-drv_obstr:
			if obstr {
				timer.TimerKill()
			} else if !obstr && e.behaviour == EB_DoorOpen {
				timer.TimerStart(e.config.doorOpenDuration_s)
			}
		case e.requests = <-floor_requests:
			switch e.behaviour {
			case EB_DoorOpen:
			case EB_Moving:
			case EB_Idle:
				dirnBehaviourPair := requests_chooseDirection(e)
				e.dirn = dirnBehaviourPair.dirn
				e.behaviour = dirnBehaviourPair.behaviour

				switch e.behaviour {
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					timer.TimerStart(e.config.doorOpenDuration_s)

				case EB_Moving:
					elevio.SetMotorDirection(e.dirn)

				case EB_Idle:
				}
			}
		}
	}
}
