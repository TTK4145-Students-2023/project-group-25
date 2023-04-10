package elevfsm

import (
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	"time"
)

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
	Floor            int
	Dirn             elevio.MotorDirection
	Behaviour        ElevatorBehaviour
	CabRequests      [dt.N_FLOORS]bool
	HallRequests     [dt.N_FLOORS][2]bool
	doorOpenDuration time.Duration
}
type WD_Role string

const (
	WD_ALIVE     WD_Role = "alive"
	WD_DEAD      WD_Role = "dead"
	watchDogTime         = 5 * time.Second
)

func FSM(
	hallRequestsCh <-chan [dt.N_FLOORS][2]bool,
	cabButtonEventCh <-chan elevio.ButtonEvent,
	floorCh <-chan int,
	obstrCh <-chan bool,
	elevDataCh chan<- dt.ElevData,
	executedHallOrdersCh chan<- []elevio.ButtonEvent,
	initCabRequestsCh <-chan [dt.N_FLOORS]bool,
	peerTxEnableCh chan<- bool) {

	var (
		watchDogTimer  = time.NewTimer(time.Hour)
		watchDogStatus = WD_ALIVE

		executedHallOrders      = []elevio.ButtonEvent{}
		executedHallOrdersTimer = time.NewTimer(time.Hour)

		obstr         = false
		doorOpenTimer = time.NewTimer(time.Hour)
		elevDataTimer = time.NewTimer(time.Hour)

		e = Elevator{
			Floor:       -1,
			Dirn:        elevio.MD_Stop,
			Behaviour:   EB_IDLE,
			CabRequests: [dt.N_FLOORS]bool{false, false, false, false},
			HallRequests: [dt.N_FLOORS][2]bool{
				{false, false},
				{false, false},
				{false, false},
				{false, false}},
			doorOpenDuration: 3 * time.Second,
		}
	)
	doorOpenTimer.Stop()
	elevDataTimer.Stop()
	watchDogTimer.Stop()
	executedHallOrdersTimer.Stop()

initialization:
	for {
		select {
		case e.Floor = <-floorCh:
			elevio.SetMotorDirection(elevio.MD_Stop)

			e.CabRequests = <-initCabRequestsCh
			for floor, order := range e.CabRequests {
				elevio.SetButtonLamp(elevio.BT_Cab, floor, order)
			}

			dirnBehaviourPair := requests_chooseDirection(e)
			e.Behaviour = dirnBehaviourPair.Behaviour
			e.Dirn = dirnBehaviourPair.Dirn
			elevDataTimer.Reset(1)

			switch e.Behaviour {
			case EB_IDLE:
			case EB_DOOR_OPEN:
				elevio.SetDoorOpenLamp(true)
				doorOpenTimer.Reset(e.doorOpenDuration)

			case EB_MOVING:
				elevio.SetMotorDirection(e.Dirn)
			}
			break initialization
		default:
			elevio.SetMotorDirection(elevio.MD_Down)
		}
	}

	for {
		select {
		case e.HallRequests = <-hallRequestsCh:
			switch e.Behaviour {
			case EB_DOOR_OPEN:
			case EB_MOVING:
			case EB_IDLE:
				dirnBehaviourPair := requests_chooseDirection(e)
				e.Behaviour = dirnBehaviourPair.Behaviour
				e.Dirn = dirnBehaviourPair.Dirn
				elevDataTimer.Reset(1)

				switch e.Behaviour {
				case EB_IDLE:
				case EB_DOOR_OPEN:
					elevio.SetDoorOpenLamp(true)
					doorOpenTimer.Reset(e.doorOpenDuration)

				case EB_MOVING:
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case cabButtonEvent := <-cabButtonEventCh:
			e.CabRequests[cabButtonEvent.Floor] = true
			elevio.SetButtonLamp(elevio.BT_Cab, cabButtonEvent.Floor, true)
			switch e.Behaviour {
			case EB_DOOR_OPEN:
			case EB_MOVING:
			case EB_IDLE:
				dirnBehaviourPair := requests_chooseDirection(e)
				e.Behaviour = dirnBehaviourPair.Behaviour
				e.Dirn = dirnBehaviourPair.Dirn
				elevDataTimer.Reset(1)

				switch e.Behaviour {
				case EB_IDLE:
				case EB_DOOR_OPEN:
					elevio.SetDoorOpenLamp(true)
					doorOpenTimer.Reset(e.doorOpenDuration)

				case EB_MOVING:
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case e.Floor = <-floorCh:
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
						doorOpenTimer.Reset(e.doorOpenDuration)
					}
				}
			}
		case <-doorOpenTimer.C:
			switch e.Behaviour {
			case EB_IDLE:
			case EB_MOVING:
			case EB_DOOR_OPEN:
				e.CabRequests[e.Floor] = false
				elevio.SetButtonLamp(elevio.BT_Cab, e.Floor, false)
				executedHallOrders = requests_getExecutedHallOrders(e)
				if len(executedHallOrders) > 0 {
					e = requests_clearLocalHallRequest(e, executedHallOrders)
					executedHallOrdersTimer.Reset(1)
				}
				dirnBehaviourPair := requests_chooseDirection(e)
				if e.Dirn != dirnBehaviourPair.Dirn || e.Behaviour != dirnBehaviourPair.Behaviour {
					e.Behaviour = dirnBehaviourPair.Behaviour
					e.Dirn = dirnBehaviourPair.Dirn
					elevDataTimer.Reset(1)
				}
				switch e.Behaviour {
				case EB_DOOR_OPEN:
					doorOpenTimer.Reset(e.doorOpenDuration)
				case EB_MOVING, EB_IDLE:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(e.Dirn)
				}
			}
		case obstr = <-obstrCh:
			if obstr {
				doorOpenTimer.Stop()
			}
			switch e.Behaviour {
			case EB_IDLE:
			case EB_MOVING:
			case EB_DOOR_OPEN:
				if !obstr {
					doorOpenTimer.Reset(e.doorOpenDuration)
				}
			}
		case <-elevDataTimer.C:
			select {
			case elevDataCh <- getElevatorData(e):
				watchDogTimer.Reset(watchDogTime)
				switch watchDogStatus {
				case WD_ALIVE:
				case WD_DEAD:
					peerTxEnableCh <- true
					watchDogStatus = WD_ALIVE
				}
			default:
				elevDataTimer.Reset(1)
			}
		case <-watchDogTimer.C:
			switch e.Behaviour {
			case EB_IDLE:
			case EB_MOVING, EB_DOOR_OPEN:
				peerTxEnableCh <- false
				watchDogStatus = WD_DEAD
			}
		case <-executedHallOrdersTimer.C:
			executedHallOrdersCopy := make([]elevio.ButtonEvent, len(executedHallOrders))
			copy(executedHallOrdersCopy, executedHallOrders)
			select {
			case executedHallOrdersCh <- executedHallOrdersCopy:
			default:
				executedHallOrdersTimer.Reset(1)
			}
		}
	}
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
