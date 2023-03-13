package elevfsm

import (
	elevio "project/localElevator/elev_driver"
)

func requests_mergeHallAndCab(hallRequests [][2]bool, cabRequests []bool) [N_FLOORS][N_BUTTONS]bool {
	var requests [N_FLOORS][N_BUTTONS]bool
	for i := 0; i < N_FLOORS; i++ {
		requests[i] = [N_BUTTONS]bool{hallRequests[i][0], hallRequests[i][1], cabRequests[i]}
	}
	return requests
}

func requests_above(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.HallRequests, e.CabRequests)
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_below(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.HallRequests, e.CabRequests)
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_here(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.HallRequests, e.CabRequests)
	for btn := 0; btn < N_BUTTONS; btn++ {
		if requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: EB_DoorOpen}
		}
		if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: EB_Idle}

	case elevio.MD_Down:
		if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: EB_DoorOpen}
		}
		if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: EB_Idle}

	case elevio.MD_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
		if requests_here(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: EB_DoorOpen}
		}
		if requests_above(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Up, Behaviour: EB_Moving}
		}
		if requests_below(e) {
			return DirnBehaviourPair{Dirn: elevio.MD_Down, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: EB_Idle}

	default:
		return DirnBehaviourPair{Dirn: elevio.MD_Stop, Behaviour: EB_Idle}
	}
}

func requests_shouldStop(e Elevator) bool {
	requests := requests_mergeHallAndCab(e.HallRequests, e.CabRequests)
	switch e.Dirn {
	case elevio.MD_Down:
		return requests[e.Floor][elevio.BT_HallDown] ||
			requests[e.Floor][elevio.BT_Cab] ||
			!requests_below(e)
	case elevio.MD_Up:
		return requests[e.Floor][elevio.BT_HallUp] ||
			requests[e.Floor][elevio.BT_Cab] ||
			!requests_above(e)
	case elevio.MD_Stop:
		return true
	default:
		return true
	}
}

func requests_getHallOrdersExecuted(e Elevator) []elevio.ButtonEvent {
	if e.Config.ClearRequestVariant == CV_All {
		return []elevio.ButtonEvent{
			{Floor: e.Floor, Button: elevio.BT_HallDown},
			{Floor: e.Floor, Button: elevio.BT_HallUp}}
	}

	requests := requests_mergeHallAndCab(e.HallRequests, e.CabRequests)
	orders := []elevio.ButtonEvent{}
	switch e.Dirn {
	case elevio.MD_Stop:
		if requests[e.Floor][elevio.BT_HallDown] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallDown})
		}
		if requests[e.Floor][elevio.BT_HallUp] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallUp})
		}
	case elevio.MD_Up:
		if requests[e.Floor][elevio.BT_HallDown] && !requests_above(e) {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallDown})
		}
		if requests[e.Floor][elevio.BT_HallUp] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallUp})
		}
	case elevio.MD_Down:
		if requests[e.Floor][elevio.BT_HallUp] && !requests_below(e) {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallUp})
		}
		if requests[e.Floor][elevio.BT_HallDown] {
			orders = append(orders, elevio.ButtonEvent{Floor: e.Floor, Button: elevio.BT_HallDown})
		}
	}
	return orders
}

func requests_clearLocalHallRequest(e Elevator, clearEvent []elevio.ButtonEvent) Elevator {
	for i := 0; i < len(clearEvent); i++ {
		e.HallRequests[clearEvent[i].Floor][clearEvent[i].Button] = false
	}
	return e
}
