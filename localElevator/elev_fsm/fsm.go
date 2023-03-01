package elevfsm

import (
	elevio "Module-go/localElevator/elev_driver"
	elevtimer "Module-go/localElevator/elev_timer"
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
	floor        int
	dirn         elevio.MotorDirection
	cabRequests  [N_FLOORS]bool
	hallRequests [N_FLOORS][2]bool
	behaviour    ElevatorBehaviour
	config       Config
}

type ElevatorData struct {
	behaviour   ElevatorBehaviour
	floor       int
	dirn        elevio.MotorDirection
	cabRequests [N_FLOORS]bool
}

type ElevatorBehaviour int

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_Moving   ElevatorBehaviour = 1
	EB_DoorOpen ElevatorBehaviour = 2
)

func getElevatorData(e Elevator) ElevatorData {
	return ElevatorData{behaviour: e.behaviour, floor: e.floor, dirn: e.dirn, cabRequests: e.cabRequests}
}

func uninitializedElevator() Elevator {
	return Elevator{
		floor:        -1,
		dirn:         elevio.MD_Stop,
		cabRequests:  [N_FLOORS]bool{},
		hallRequests: [N_FLOORS][2]bool{},
		behaviour:    EB_Idle,
		config:       Config{clearRequestVariant: CV_InDirn, doorOpenDuration_s: 3},
	}
}

func FSM(
	floor_hallRequests <-chan [N_FLOORS][2]bool,
	floor_cabButtonEvent <-chan elevio.ButtonEvent,
	drv_floors <-chan int,
	drv_obstr <-chan bool,
	elev_data chan<- ElevatorData,
	handler_hallOrdersExecuted chan<- []elevio.ButtonEvent) {

	e := uninitializedElevator()
	obstr := false
	timeout := make(chan bool)

	go elevtimer.TimerMain(timeout)

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
	elev_data <- getElevatorData(e)
	for {
		select {
		case e.hallRequests = <-floor_hallRequests:
			switch e.behaviour {
			case EB_DoorOpen:
			case EB_Moving:
			case EB_Idle:

				dirnBehaviourPair := requests_chooseDirection(e)
				e.dirn = dirnBehaviourPair.dirn
				e.behaviour = dirnBehaviourPair.behaviour
				elev_data <- getElevatorData(e)
				switch e.behaviour {
				case EB_Idle:
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					elevtimer.TimerStart(e.config.doorOpenDuration_s)

				case EB_Moving:
					elevio.SetMotorDirection(e.dirn)
				}
			}

		case cabButtonEvent := <-floor_cabButtonEvent:
			e.cabRequests[cabButtonEvent.Floor] = true
			elevio.SetButtonLamp(elevio.BT_Cab, cabButtonEvent.Floor, true)
			switch e.behaviour {
			case EB_DoorOpen:
			case EB_Moving:
			case EB_Idle:

				dirnBehaviourPair := requests_chooseDirection(e)
				e.dirn = dirnBehaviourPair.dirn
				e.behaviour = dirnBehaviourPair.behaviour
				elev_data <- getElevatorData(e)

				switch e.behaviour {
				case EB_Idle:
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					elevtimer.TimerStart(e.config.doorOpenDuration_s)

				case EB_Moving:
					elevio.SetMotorDirection(e.dirn)
				}
			}

		case e.floor = <-drv_floors:
			elevio.SetFloorIndicator(e.floor)
			switch e.behaviour {
			case EB_Idle:
			case EB_DoorOpen:
			case EB_Moving:
				if requests_shouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					e.behaviour = EB_DoorOpen
					elev_data <- getElevatorData(e)
					if !obstr {
						elevtimer.TimerStart(e.config.doorOpenDuration_s)
					}
				}
			}

		case <-timeout:
			switch e.behaviour {
			case EB_Idle:
			case EB_Moving:
			case EB_DoorOpen:
				e.cabRequests[e.floor] = false
				elevio.SetButtonLamp(elevio.BT_Cab, e.floor, false)

				hallOrdersExecuted := requests_getHallOrdersExecuted(e)
				e = requests_clearLocalHallRequest(e, hallOrdersExecuted)
				handler_hallOrdersExecuted <- hallOrdersExecuted

				dirnBehaviourPair := requests_chooseDirection(e)
				e.dirn = dirnBehaviourPair.dirn
				e.behaviour = dirnBehaviourPair.behaviour
				elev_data <- getElevatorData(e)

				switch e.behaviour {
				case EB_DoorOpen:
					elevtimer.TimerStart(e.config.doorOpenDuration_s)
				case EB_Moving:
					fallthrough
				case EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(e.dirn)
				}
			}

		case obstr = <-drv_obstr:
			if obstr {
				elevtimer.TimerKill()
			}
			switch e.behaviour {
			case EB_Idle:
			case EB_Moving:
			case EB_DoorOpen:
				if !obstr {
					elevtimer.TimerStart(e.config.doorOpenDuration_s)
				}
			}
		}
	}
}
