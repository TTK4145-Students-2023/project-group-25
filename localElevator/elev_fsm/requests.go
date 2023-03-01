package elevfsm

import (
	elevio "Module-go/localElevator/elev_driver"
)

type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection
	behaviour ElevatorBehaviour
}

func requests_mergeHallAndCab(hallRequests [N_FLOORS][2]bool, cabRequests [N_FLOORS]bool) [N_FLOORS][N_BUTTONS]bool {
	var requests [N_FLOORS][N_BUTTONS]bool
	for i := 0; i < N_FLOORS; i++ {
		requests[i] = [N_BUTTONS]bool{hallRequests[i][0], hallRequests[i][1], cabRequests[i]}
	}
	return requests
}

func requests_above(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.hallRequests, e.cabRequests)
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_below(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.hallRequests, e.cabRequests)
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_here(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.hallRequests, e.cabRequests)
	for btn := 0; btn < N_BUTTONS; btn++ {
		if requests[e.floor][btn] {
			return true
		}
	}
	return false
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case elevio.MD_Up:
		if requests_above(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_DoorOpen}
		}
		if requests_below(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}

	case elevio.MD_Down:
		if requests_below(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_DoorOpen}
		}
		if requests_above(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}

	case elevio.MD_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
		if requests_here(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_DoorOpen}
		}
		if requests_above(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_Moving}
		}
		if requests_below(e) {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}

	default:
		return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}
	}
}

func requests_shouldStop(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.hallRequests, e.cabRequests)
	switch e.dirn {
	case elevio.MD_Down:
		return requests[e.floor][elevio.BT_HallDown] ||
			requests[e.floor][elevio.BT_Cab] ||
			!requests_below(e)
	case elevio.MD_Up:
		return requests[e.floor][elevio.BT_HallUp] ||
			requests[e.floor][elevio.BT_Cab] ||
			!requests_above(e)
	case elevio.MD_Stop:
		return true
	default:
		return true
	}
}

func requests_getHallOrdersExecuted(e Elevator) []elevio.ButtonEvent {
	if e.config.clearRequestVariant == CV_All {
		return []elevio.ButtonEvent{
			{Floor: e.floor, Button: elevio.BT_HallDown},
			{Floor: e.floor, Button: elevio.BT_HallUp}}
	}

	requests := requests_mergeHallAndCab(e.hallRequests, e.cabRequests)
	orders := []elevio.ButtonEvent{}
	switch e.dirn {
	case elevio.MD_Stop:
		if requests[e.floor][elevio.BT_HallDown] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallDown})
		}
		if requests[e.floor][elevio.BT_HallUp] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallUp})
		}
	case elevio.MD_Up:
		if requests[e.floor][elevio.BT_HallDown] && !requests_above(e) {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallDown})
		}
		if requests[e.floor][elevio.BT_HallUp] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallUp})
		}
	case elevio.MD_Down:
		if requests[e.floor][elevio.BT_HallUp] && !requests_below(e) {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallUp})
		}
		if requests[e.floor][elevio.BT_HallDown] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallDown})
		}
	}
	return orders
}

func requests_clearLocalHallRequest(e Elevator, clearEvent []elevio.ButtonEvent) Elevator {
	for i := 0; i < len(clearEvent); i++ {
		e.hallRequests[clearEvent[i].Floor][clearEvent[i].Button] = false
	}
	return e
}
