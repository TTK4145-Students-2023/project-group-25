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

 func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
    switch(e.config.clearRequestVariant){
    case CV_All:
        return e.floor == btn_floor
    case CV_InDirn:
        return e.floor == btn_floor && 
                (e.dirn == elevio.MD_Up   && btn_type == elevio.BT_HallUp)    ||
                (e.dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown)  ||
                (e.dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab)
    default:
        return false
    }
}

func requests_clearAtCurrentFloor(e Elevator) Elevator {
        
    switch(e.config.clearRequestVariant){
    case CV_All:
        for btn := 0; btn < N_BUTTONS; btn++{
            e.requests[e.floor][btn] = false
        }
        
    case CV_InDirn:
        e.requests[e.floor][elevio.BT_Cab] = false
        switch(e.dirn){
        case elevio.MD_Up:
            if(!requests_above(e) && !e.requests[e.floor][elevio.BT_HallUp]){
                e.requests[e.floor][elevio.BT_HallDown] = false
            }
            e.requests[e.floor][elevio.BT_HallUp] = false
            
        case elevio.MD_Down:
            if(!requests_below(e) && !e.requests[e.floor][elevio.BT_HallDown]){
                e.requests[e.floor][elevio.BT_HallUp] = false
            }
            e.requests[e.floor][elevio.BT_HallDown] = false
            
        case elevio.MD_Stop:
        default:
            e.requests[e.floor][elevio.BT_HallUp] = false
            e.requests[e.floor][elevio.BT_HallDown] = false

        }
    default:
    }
    return e
}

