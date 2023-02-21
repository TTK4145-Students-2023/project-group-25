package fsm

import (
	"Driver-go/elevio"
)

type DirnBehaviourPair struct {
    dirn       elevio.MotorDirection
    behaviour  ElevatorBehaviour
}

func requests_above(e Elevator) bool{
    for f := e.floor+1; f < N_FLOORS; f++{
        for btn := 0; btn < N_BUTTONS; btn++{
            if e.requests[f][btn] {
                return true
            }
        }
    }
    return false
}

func requests_below(e Elevator) bool{
    for f := 0; f < e.floor; f++{
        for btn := 0; btn < N_BUTTONS; btn++{
            if e.requests[f][btn]{
                return true
            }
        }
    }
    return false
}

func requests_here(e Elevator) bool{
    for btn := 0; btn < N_BUTTONS; btn++{
        if e.requests[e.floor][btn] {
            return true
        }
    }
    return false
}


func requests_chooseDirection(e Elevator) DirnBehaviourPair {
    switch(e.dirn){
    case elevio.MD_Up:
        if requests_above(e){ return DirnBehaviourPair{dirn : elevio.MD_Up,   behaviour : EB_Moving} }
        if requests_here(e) { return DirnBehaviourPair{dirn : elevio.MD_Down, behaviour : EB_DoorOpen} }
        if requests_below(e){ return DirnBehaviourPair{dirn : elevio.MD_Down, behaviour : EB_Moving} }
        return DirnBehaviourPair{dirn : elevio.MD_Stop, behaviour : EB_Idle}

    case elevio.MD_Down:
        if requests_below(e){ return DirnBehaviourPair{dirn : elevio.MD_Down, behaviour : EB_Moving} }
        if requests_here(e) { return DirnBehaviourPair{dirn : elevio.MD_Up,   behaviour : EB_DoorOpen} }
        if requests_above(e){ return DirnBehaviourPair{dirn : elevio.MD_Up,   behaviour : EB_Moving} }
        return DirnBehaviourPair{dirn : elevio.MD_Stop, behaviour : EB_Idle}

    case elevio.MD_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
        if requests_here(e) { return DirnBehaviourPair{dirn : elevio.MD_Stop, behaviour : EB_DoorOpen} }
        if requests_above(e){ return DirnBehaviourPair{dirn : elevio.MD_Up,   behaviour : EB_Moving} }
        if requests_below(e){ return DirnBehaviourPair{dirn : elevio.MD_Down, behaviour : EB_Moving} }
        return DirnBehaviourPair{dirn : elevio.MD_Stop, behaviour : EB_Idle}

    default:
        return DirnBehaviourPair{dirn : elevio.MD_Stop, behaviour : EB_Idle}
    }
}

func requests_shouldStop(e Elevator) bool {
    switch e.dirn {
    case elevio.MD_Down:
        return e.requests[e.floor][elevio.BT_HallDown] ||
               e.requests[e.floor][elevio.BT_Cab]      ||
               !requests_below(e)
    case elevio.MD_Up:
        return e.requests[e.floor][elevio.BT_HallUp]   ||
               e.requests[e.floor][elevio.BT_Cab]      ||
               !requests_above(e)
    case elevio.MD_Stop:
        return true
    default:
        return true
    }
}

func requests_calculateOrdersToBeCleared(e Elevator) []elevio.ButtonEvent {
    if e.config.clearRequestVariant == CV_All {return []elevio.ButtonEvent{
                                                {Floor: e.floor, Button: elevio.BT_HallDown},
                                                {Floor: e.floor, Button: elevio.BT_HallUp},
                                                {Floor: e.floor, Button: elevio.BT_Cab}}}

    orders := []elevio.ButtonEvent{}
    switch e.dirn{
    case elevio.MD_Stop:
        if e.requests[e.floor][elevio.BT_HallDown] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallDown})}
        if e.requests[e.floor][elevio.BT_HallUp] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallUp})}
        if e.requests[e.floor][elevio.BT_Cab] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_Cab})}
    case elevio.MD_Up:
        if e.requests[e.floor][elevio.BT_HallUp] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallUp})}
        if e.requests[e.floor][elevio.BT_Cab] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_Cab})}
    case elevio.MD_Down:
        if e.requests[e.floor][elevio.BT_HallDown] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_HallDown})}
        if e.requests[e.floor][elevio.BT_Cab] {orders = append(orders , elevio.ButtonEvent{Floor: e.floor, Button: elevio.BT_Cab})}
    }
    return orders
}