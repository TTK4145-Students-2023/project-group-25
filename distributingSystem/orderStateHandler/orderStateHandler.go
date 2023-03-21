package orderStateHandler

import (
	"fmt"
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	"reflect"
	"time"
)

const (
	STATE_NONE      dt.OrderState = "none"
	STATE_NEW       dt.OrderState = "new"
	STATE_CONFIRMED dt.OrderState = "confirmed"
)

func OrderStateHandler(localIP string,
	hallBtnPressCh <-chan elevio.ButtonEvent,
	hallOrdersExecutedCh <-chan []elevio.ButtonEvent,
	hallOrderArrayCh chan<- [dt.N_FLOORS][2]bool, //confHallOrders?
	peerUpdateCh <-chan peers.PeerUpdate,
) {
	var (
		AllNodeOrderStates = map[string][dt.N_FLOORS][2]dt.OrderState{}
		peerList           = peers.PeerUpdate{}

		hallOrderArrayTimer = time.NewTimer(1)
		broadCastTimer      = time.NewTimer(1)

		receiveCh  = make(chan dt.NodeOrderStates)
		transmitCh = make(chan dt.NodeOrderStates)
	)
	broadCastTimer.Stop()
	hallOrderArrayTimer.Stop()

	go bcast.Receiver(15668, receiveCh)
	go bcast.Transmitter(15668, transmitCh)

	for {
		select {
		case peerList = <-peerUpdateCh:
			newNodeOrderStates := removeDeadNodes(peerList, AllNodeOrderStates, localIP)
			newNodeOrderStates = addNewEmptyNodes(peerList, newNodeOrderStates)
			newNodeOrderStates = withdrawOrderConfirmations(peerList, newNodeOrderStates)
			AllNodeOrderStates = newNodeOrderStates
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)
			hallOrderArrayTimer.Reset(1)

		case newNodeOrderStates := <-receiveCh:
			newData := newNodeOrderStates.OrderStates
			senderIP := newNodeOrderStates.IP
			if senderIP == localIP || reflect.DeepEqual(newData, AllNodeOrderStates[senderIP]) {
				break
			}
			AllNodeOrderStates[senderIP] = newData

			for nodeIP := range AllNodeOrderStates {
				if nodeIP == localIP {
					continue
				}
				newOrderStates := AllNodeOrderStates[localIP]
				for floor := range newData {
					for btn, inputBtnState := range newData[floor] {
						newState := updateOrderState(inputBtnState, newOrderStates[floor][btn])
						newOrderStates[floor][btn] = newState
						switch newState {
						case STATE_NEW:
						case STATE_CONFIRMED:
							elevio.SetButtonLamp(elevio.ButtonType(btn), floor, true)
						case STATE_NONE:
							elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)
						}
					}
				}
				if AllNodeOrderStates[localIP] != newOrderStates {
					AllNodeOrderStates[localIP] = newOrderStates
					hallOrderArrayTimer.Reset(1)
				}
			}

		case BtnPress := <-hallBtnPressCh:
			newOrderStates := AllNodeOrderStates[localIP]
			if newOrderStates[BtnPress.Floor][BtnPress.Button] == STATE_NONE {
				newOrderStates[BtnPress.Floor][BtnPress.Button] = STATE_NEW
				AllNodeOrderStates[localIP] = newOrderStates
				hallOrderArrayTimer.Reset(1)
			}
		case executedOrders := <-hallOrdersExecutedCh:
			newOrderStates := AllNodeOrderStates[localIP]
			for _, order := range executedOrders {
				if newOrderStates[order.Floor][order.Button] == STATE_CONFIRMED {
					newOrderStates[order.Floor][order.Button] = STATE_NONE
					elevio.SetButtonLamp(order.Button, order.Floor, false)
				}
			}
			AllNodeOrderStates[localIP] = newOrderStates
			hallOrderArrayTimer.Reset(1)
		case <-broadCastTimer.C:
			transmitCh <- dt.NodeOrderStates{IP: localIP, OrderStates: AllNodeOrderStates[localIP]}
			fmt.Printf("NOS: %+v\n", AllNodeOrderStates)
			fmt.Printf("peerlist: %+v\n", peerList)

			broadCastTimer.Reset(dt.BROADCAST_PERIOD)

		case <-hallOrderArrayTimer.C:
			select {
			case hallOrderArrayCh <- orderStatesToBool(AllNodeOrderStates[localIP]):
			default:
				hallOrderArrayTimer.Reset(1)
			}
		}
		newOrderStates := AllNodeOrderStates[localIP]
		for floor := 0; floor < dt.N_FLOORS; floor++ {
			for btn := 0; btn < dt.N_BUTTONS-1; btn++ {
				if orderCanBeConfirmed(floor, btn, AllNodeOrderStates) {
					newOrderStates[floor][btn] = STATE_CONFIRMED
					elevio.SetButtonLamp(elevio.ButtonType(btn), floor, true)
				}
			}
		}
		if AllNodeOrderStates[localIP] != newOrderStates {
			AllNodeOrderStates[localIP] = newOrderStates
			hallOrderArrayTimer.Reset(1)
		}
	}
}

func orderCanBeConfirmed(floor, btn int, AllNodeOrderStates map[string][dt.N_FLOORS][2]dt.OrderState) bool {
	for _, nodeOrderStates := range AllNodeOrderStates {
		if nodeOrderStates[floor][btn] != STATE_NEW {
			return false
		}
	}
	return true
}

func updateOrderState(inputState, currentState dt.OrderState) dt.OrderState {
	switch inputState {
	case STATE_NONE:
		if currentState == STATE_CONFIRMED {
			return STATE_NONE
		}
	case STATE_NEW:
		if currentState == STATE_NONE {
			return STATE_NEW

		}
	case STATE_CONFIRMED:
		if currentState == STATE_NEW {
			return STATE_CONFIRMED
		}
	}
	return currentState
}

func orderStatesToBool(orderStates [dt.N_FLOORS][2]dt.OrderState) [dt.N_FLOORS][2]bool {
	hallOrders := [dt.N_FLOORS][2]bool{}
	for floor := range orderStates {
		for btn, state := range orderStates[floor] {
			if state == STATE_CONFIRMED {
				hallOrders[floor][btn] = true
			}
		}
	}
	return hallOrders
}

func addNewEmptyNodes(peerList peers.PeerUpdate, AllNodeOrderStates map[string][dt.N_FLOORS][2]dt.OrderState) map[string][dt.N_FLOORS][2]dt.OrderState {
	outputMap := make(map[string][dt.N_FLOORS][2]dt.OrderState)
	for key, value := range AllNodeOrderStates {
		outputMap[key] = value
	}
	for _, nodeIP := range peerList.Peers {
		if _, nodeOrdersSaved := outputMap[nodeIP]; !nodeOrdersSaved {
			outputMap[nodeIP] = [dt.N_FLOORS][2]dt.OrderState{{STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}}
		}
	}
	return outputMap
}

func removeDeadNodes(peerList peers.PeerUpdate, AllNodeOrderStates map[string][dt.N_FLOORS][2]dt.OrderState, localIP string) map[string][dt.N_FLOORS][2]dt.OrderState {
	outputMap := make(map[string][dt.N_FLOORS][2]dt.OrderState)
	for IP := range AllNodeOrderStates {
		if contains(peerList.Peers, IP) || IP == localIP {
			outputMap[IP] = AllNodeOrderStates[IP]
		}
	}
	return outputMap
}

func withdrawOrderConfirmations(peerList peers.PeerUpdate, AllNodeOrderStates map[string][dt.N_FLOORS][2]dt.OrderState) map[string][dt.N_FLOORS][2]dt.OrderState {
	outputMap := make(map[string][dt.N_FLOORS][2]dt.OrderState)
	for key, value := range AllNodeOrderStates {
		outputMap[key] = value
	}
	for _, nodeIP := range peerList.Peers {
		if orderStates, nodeOrdersSaved := outputMap[nodeIP]; nodeOrdersSaved {
			for floor := range orderStates {
				for btn, state := range orderStates[floor] {
					switch state {
					case STATE_NONE:
					case STATE_NEW:
					case STATE_CONFIRMED:
						orderStates[floor][btn] = STATE_NEW
						outputMap[nodeIP] = orderStates
					}
				}
			}
		}
	}
	return outputMap
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
