package elevfsm

import (
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	elevtimer "project/localElevator/elev_timer"
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

type ElevatorConfig struct {
	ClearRequestVariant ClearRequestVariant
	DoorOpenDuration_s  time.Duration
}

type ElevatorBehaviour string

const (
	EB_DoorOpen ElevatorBehaviour = "doorOpen"
	EB_Moving   ElevatorBehaviour = "moving"
	EB_Idle     ElevatorBehaviour = "idle"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

type Elevator struct {
	Floor        int
	Dirn         elevio.MotorDirection
	CabRequests  []bool
	HallRequests [][2]bool
	Behaviour    ElevatorBehaviour
	Config       ElevatorConfig
}

func getElevatorData(e Elevator) dt.ElevDataJSON {
	dirnToString := map[elevio.MotorDirection]string{
		elevio.MD_Down: "down",
		elevio.MD_Up:   "up",
		elevio.MD_Stop: "stop"}

	return dt.ElevDataJSON{
		Behavior:    string(e.Behaviour),
		Floor:       e.Floor,
		Direction:   dirnToString[e.Dirn],
		CabRequests: e.CabRequests}
}

func uninitializedElevator() Elevator {
	return Elevator{
		Floor:        -1,
		Dirn:         elevio.MD_Stop,
		CabRequests:  []bool{false, false, false, false},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		Behaviour:    EB_Idle,
		Config:       ElevatorConfig{ClearRequestVariant: CV_InDirn, DoorOpenDuration_s: 3},
	}
}

func FSM(
	floor_hallRequests <-chan [][2]bool,
	floor_cabButtonEvent <-chan elevio.ButtonEvent,
	drv_floors <-chan int,
	drv_obstr <-chan bool,
	elev_data chan<- dt.ElevDataJSON,
	handler_hallOrdersExecuted chan<- []elevio.ButtonEvent) {

	e := uninitializedElevator()
	obstr := false
	timeout := make(chan bool)

	go elevtimer.TimerMain(timeout)

	select {
	case e.Floor = <-drv_floors:
	default:
	}

	if e.Floor == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
		e.Dirn = elevio.MD_Down
		e.Behaviour = EB_Moving

		e.Floor = <-drv_floors

		elevio.SetMotorDirection(elevio.MD_Stop)
		e.Dirn = elevio.MD_Stop
		e.Behaviour = EB_Idle
	}
	elev_data <- getElevatorData(e)
	for {
		select {
		case e.HallRequests = <-floor_hallRequests:
			switch e.Behaviour {
			case EB_DoorOpen:
			case EB_Moving:
			case EB_Idle:

				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elev_data <- getElevatorData(e)
				}
				switch e.Behaviour {
				case EB_Idle:
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					elevtimer.TimerStart(e.Config.DoorOpenDuration_s)

				case EB_Moving:
					elevio.SetMotorDirection(e.Dirn)
				}
			}

		case cabButtonEvent := <-floor_cabButtonEvent:
			e.CabRequests[cabButtonEvent.Floor] = true
			elevio.SetButtonLamp(elevio.BT_Cab, cabButtonEvent.Floor, true)
			switch e.Behaviour {
			case EB_DoorOpen:
			case EB_Moving:
			case EB_Idle:

				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elev_data <- getElevatorData(e)
				}

				switch e.Behaviour {
				case EB_Idle:
				case EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					elevtimer.TimerStart(e.Config.DoorOpenDuration_s)

				case EB_Moving:
					elevio.SetMotorDirection(e.Dirn)
				}
			}

		case e.Floor = <-drv_floors:
			elevio.SetFloorIndicator(e.Floor)
			switch e.Behaviour {
			case EB_Idle:
			case EB_DoorOpen:
			case EB_Moving:
				if requests_shouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					e.Behaviour = EB_DoorOpen
					elev_data <- getElevatorData(e)
					if !obstr {
						elevtimer.TimerStart(e.Config.DoorOpenDuration_s)
					}
				}
			}

		case <-timeout:
			switch e.Behaviour {
			case EB_Idle:
			case EB_Moving:
			case EB_DoorOpen:
				e.CabRequests[e.Floor] = false
				elevio.SetButtonLamp(elevio.BT_Cab, e.Floor, false)

				hallOrdersExecuted := requests_getHallOrdersExecuted(e)
				e = requests_clearLocalHallRequest(e, hallOrdersExecuted)
				handler_hallOrdersExecuted <- hallOrdersExecuted

				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elev_data <- getElevatorData(e)
				}

				switch e.Behaviour {
				case EB_DoorOpen:
					elevtimer.TimerStart(e.Config.DoorOpenDuration_s)
				case EB_Moving:
					fallthrough
				case EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(e.Dirn)
				}
			}

		case obstr = <-drv_obstr:
			if obstr {
				elevtimer.TimerKill()
			}
			switch e.Behaviour {
			case EB_Idle:
			case EB_Moving:
			case EB_DoorOpen:
				if !obstr {
					elevtimer.TimerStart(e.Config.DoorOpenDuration_s)
				}
			}
		}
	}
}
