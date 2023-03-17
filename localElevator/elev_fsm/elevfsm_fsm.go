package elevfsm

import (
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	"time"
)

type ClearRequestVariant int

const (
	// Assume everyone waiting for the elevator gets on the elevator, even if
	// they will be traveling in the "wrong" direction for a while
	CV_ALL ClearRequestVariant = 0

	// Assume that only those that want to travel in the current direction
	// enter the elevator, and keep waiting outside otherwise
	CV_INDIRN ClearRequestVariant = 1
)

type ElevatorConfig struct {
	ClearRequestVariant ClearRequestVariant
	DoorOpenDuration_s  time.Duration
}

type ElevatorBehaviour string

const (
	EB_DOOR_OPEN ElevatorBehaviour = "doorOpen"
	EB_MOVING    ElevatorBehaviour = "moving"
	EB_IDLE      ElevatorBehaviour = "idle"
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

type Elevator struct {
	Floor        int
	Dirn         elevio.MotorDirection
	CabRequests  [dt.N_FLOORS]bool
	HallRequests [dt.N_FLOORS][2]bool
	Behaviour    ElevatorBehaviour
	Config       ElevatorConfig
}

func getElevatorData(e Elevator) dt.ElevData {
	dirnToString := map[elevio.MotorDirection]string{
		elevio.MD_Down: "down",
		elevio.MD_Up:   "up",
		elevio.MD_Stop: "stop"}

	return dt.ElevData{
		Behavior:    string(e.Behaviour),
		Floor:       e.Floor,
		Direction:   dirnToString[e.Dirn],
		CabRequests: e.CabRequests}
}

func FSM(
	floor_hallRequests <-chan [dt.N_FLOORS][2]bool,
	floor_cabButtonEvent <-chan elevio.ButtonEvent,
	drv_floors <-chan int,
	drv_obstr <-chan bool,
	elev_data chan<- dt.ElevData,
	handler_hallOrdersExecuted chan<- []elevio.ButtonEvent,
	cabRequestsToElevCh <-chan [dt.N_FLOORS]bool) {

	obstr := false
	hallOrdersExecuted := []elevio.ButtonEvent{}
	e := Elevator{
		Floor:        -1,
		Dirn:         elevio.MD_Stop,
		CabRequests:  [dt.N_FLOORS]bool{false, false, false, false},
		HallRequests: [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		Behaviour:    EB_IDLE,
		Config:       ElevatorConfig{ClearRequestVariant: CV_INDIRN, DoorOpenDuration_s: 3 * time.Second},
	}

	ElevTimer := time.NewTimer(1)
	ElevTimer.Stop()

	elevDataTimer := time.NewTimer(1)
	hallOrdersExecutedTimer := time.NewTimer(1)
	hallOrdersExecutedTimer.Stop()
	select {
	case e.Floor = <-drv_floors:
	default:
	}

	if e.Floor == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
		e.Dirn = elevio.MD_Down
		e.Behaviour = EB_MOVING

		e.Floor = <-drv_floors

		elevio.SetMotorDirection(elevio.MD_Stop)
		e.Dirn = elevio.MD_Stop
		e.Behaviour = EB_IDLE
		elevDataTimer.Reset(1)
	}

	e.CabRequests = <-cabRequestsToElevCh
	for floor, order := range e.CabRequests {
		elevio.SetButtonLamp(elevio.BT_Cab, floor, order)
	}
	dirnBehaviourPair := requests_chooseDirection(e)
	if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
		e.Dirn = dirnBehaviourPair.Dirn
		e.Behaviour = dirnBehaviourPair.Behaviour
		elevDataTimer.Reset(1)

		switch e.Behaviour {
		case EB_IDLE:
		case EB_DOOR_OPEN:
			elevio.SetDoorOpenLamp(true)
			ElevTimer.Reset(e.Config.DoorOpenDuration_s)

		case EB_MOVING:
			elevio.SetMotorDirection(e.Dirn)
		}
	}
	for {
		select {
		case e.HallRequests = <-floor_hallRequests:
			switch e.Behaviour {
			case EB_DOOR_OPEN:
			case EB_MOVING:
			case EB_IDLE:
				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elevDataTimer.Reset(1)
				}
				switch e.Behaviour {
				case EB_IDLE:
				case EB_DOOR_OPEN:
					elevio.SetDoorOpenLamp(true)
					ElevTimer.Reset(e.Config.DoorOpenDuration_s)

				case EB_MOVING:
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case cabButtonEvent := <-floor_cabButtonEvent:
			e.CabRequests[cabButtonEvent.Floor] = true
			elevio.SetButtonLamp(elevio.BT_Cab, cabButtonEvent.Floor, true)
			switch e.Behaviour {
			case EB_DOOR_OPEN:
			case EB_MOVING:
			case EB_IDLE:
				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elevDataTimer.Reset(1)
				}
				switch e.Behaviour {
				case EB_IDLE:
				case EB_DOOR_OPEN:
					elevio.SetDoorOpenLamp(true)
					ElevTimer.Reset(e.Config.DoorOpenDuration_s)

				case EB_MOVING:
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case e.Floor = <-drv_floors:
			elevio.SetFloorIndicator(e.Floor)
			switch e.Behaviour {
			case EB_IDLE:
			case EB_DOOR_OPEN:
			case EB_MOVING:
				elevDataTimer.Reset(1)
				if requests_shouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					e.Behaviour = EB_DOOR_OPEN
					if !obstr {
						ElevTimer.Reset(e.Config.DoorOpenDuration_s)
					}
				}
			}
		case <-ElevTimer.C:
			switch e.Behaviour {
			case EB_IDLE:
			case EB_MOVING:
			case EB_DOOR_OPEN:
				e.CabRequests[e.Floor] = false
				elevio.SetButtonLamp(elevio.BT_Cab, e.Floor, false)
				hallOrdersExecuted = requests_getHallOrdersExecuted(e)
				if len(hallOrdersExecuted) > 0 {
					e = requests_clearLocalHallRequest(e, hallOrdersExecuted)
					hallOrdersExecutedTimer.Reset(1)
				}
				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Dirn = dirnBehaviourPair.Dirn
					e.Behaviour = dirnBehaviourPair.Behaviour
					elevDataTimer.Reset(1)
				}
				switch e.Behaviour {
				case EB_DOOR_OPEN:
					ElevTimer.Reset(e.Config.DoorOpenDuration_s)
				case EB_MOVING:
					fallthrough
				case EB_IDLE:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case obstr = <-drv_obstr:
			if obstr {
				ElevTimer.Stop()
			}
			switch e.Behaviour {
			case EB_IDLE:
			case EB_MOVING:
			case EB_DOOR_OPEN:
				if !obstr {
					ElevTimer.Reset(e.Config.DoorOpenDuration_s)
				}
			}
		case <-elevDataTimer.C:
			select {
			case elev_data <- getElevatorData(e):
			default:
				elevDataTimer.Reset(1)
			}
		case <-hallOrdersExecutedTimer.C:
			executedOrdersToAssigner := make([]elevio.ButtonEvent, len(hallOrdersExecuted))
			copy(executedOrdersToAssigner, hallOrdersExecuted)
			select {
			case handler_hallOrdersExecuted <- executedOrdersToAssigner:
			default:
				hallOrdersExecutedTimer.Reset(1)
			}
		}
	}
}
