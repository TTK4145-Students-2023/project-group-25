package orderStateHandler

import (
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	"time"
)

// States for hall requests
const (
	STATE_NONE      dt.OrderState = "none"
	STATE_NEW       dt.OrderState = "new"
	STATE_CONFIRMED dt.OrderState = "confirmed"
)

func OrderStateHandler(localIP string,
	allNOSfromNTWCh <-chan dt.AllNOS_WithSenderIP,
	hallBtnPressCh <-chan elevio.ButtonEvent,
	hallOrdersExecutedCh <-chan []elevio.ButtonEvent,
	hallOrderArrayCh chan<- [dt.N_FLOORS][2]bool,
	NOStoNTWCh chan<- []dt.NodeOrderStates,
	peerUpdateCh <-chan peers.PeerUpdate,
) {
	localNodeOrderStates := map[string][dt.N_FLOORS][2]dt.OrderState{}
	peerList := peers.PeerUpdate{}

	reqStateMatrixTimer := time.NewTimer(1)
	reqStateMatrixTimer.Stop()
	hallOrderArrayTimer := time.NewTimer(1)
	hallOrderArrayTimer.Stop()
	for {
		select {
		case peerList = <-peerUpdateCh:
			for _, nodeID := range peerList.Peers {
				if _, valInMap := localNodeOrderStates[nodeID]; !valInMap {
					localNodeOrderStates[nodeID] = [dt.N_FLOORS][2]dt.OrderState{{STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}, {STATE_NONE, STATE_NONE}}
					reqStateMatrixTimer.Reset(1)
					hallOrderArrayTimer.Reset(1)
				}
			}
		case allNOSfromP2P := <-allNOSfromNTWCh:
			senderIP := allNOSfromP2P.SenderIP
			senderNOS := dt.NOSSliceToMap(allNOSfromP2P.AllNOS)
			localNodeOrderStates[senderIP] = senderNOS[senderIP]
			// Iterate through the list of node IDs
			for _, nodeIP := range peerList.Peers {
				// Skip the local node
				if nodeIP == localIP {
					continue
				}
				// Compare the requestStates from the other nodes with the Local requestStates
				for floor := range senderNOS[nodeIP] {
					for btn_UpDown, other_state := range senderNOS[nodeIP][floor] {

						localOrderStates := localNodeOrderStates[localIP]
						//cyclic change of states
						switch other_state {
						case STATE_NONE:
							if localOrderStates[floor][btn_UpDown] == STATE_CONFIRMED {
								localOrderStates[floor][btn_UpDown] = STATE_NONE
								localNodeOrderStates[localIP] = localOrderStates
								elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, false)
								reqStateMatrixTimer.Reset(1)
								hallOrderArrayTimer.Reset(1)
							}
						case STATE_NEW:
							if localOrderStates[floor][btn_UpDown] == STATE_NONE {
								localOrderStates[floor][btn_UpDown] = STATE_NEW
								localNodeOrderStates[localIP] = localOrderStates
								reqStateMatrixTimer.Reset(1)
								hallOrderArrayTimer.Reset(1)
							}
						case STATE_CONFIRMED:
							if localOrderStates[floor][btn_UpDown] == STATE_NEW {
								localOrderStates[floor][btn_UpDown] = STATE_CONFIRMED
								localNodeOrderStates[localIP] = localOrderStates
								elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true)
								reqStateMatrixTimer.Reset(1)
								hallOrderArrayTimer.Reset(1)
							}
						}
					}
				}
			}
		case BtnPress := <-hallBtnPressCh:
			localStateArray := localNodeOrderStates[localIP]
			if localStateArray[BtnPress.Floor][BtnPress.Button] == STATE_NONE {
				localStateArray[BtnPress.Floor][BtnPress.Button] = STATE_NEW
				localNodeOrderStates[localIP] = localStateArray
				reqStateMatrixTimer.Reset(1)
				hallOrderArrayTimer.Reset(1)
			}
		case executedArray := <-hallOrdersExecutedCh:
			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}
				localStateArray := localNodeOrderStates[localIP]
				// Kva skjer her dersom order blir executed før state_confirmed?
				if localStateArray[btn.Floor][btn.Button] == STATE_CONFIRMED {
					localStateArray[btn.Floor][btn.Button] = STATE_NONE
					localNodeOrderStates[localIP] = localStateArray
					reqStateMatrixTimer.Reset(1)
					hallOrderArrayTimer.Reset(1)
					elevio.SetButtonLamp(btn.Button, btn.Floor, false)
				}
			}
		case <-hallOrderArrayTimer.C:
			select {
			case hallOrderArrayCh <- ConfirmedOrdersToHallOrder(localNodeOrderStates, localIP):
			default:
				hallOrderArrayTimer.Reset(1)
			}
		case <-reqStateMatrixTimer.C:
			select {
			case NOStoNTWCh <- dt.NOSMapToSlice(localNodeOrderStates):
			default:
				reqStateMatrixTimer.Reset(1)
			}
		}
		//Check if Order can be confirmed
		//If all orders across IDs is State_new, order is confirmed and sendt to order Assigner
		for floor, floorStateArray := range localNodeOrderStates[localIP] {
			for btn_UpDown := range floorStateArray {

				if floorStateArray[btn_UpDown] != STATE_NEW {
					continue
				}

				NewOrder_OnAll_IDs := true
				for _, nodeID := range peerList.Peers {
					if localNodeOrderStates[nodeID][floor][btn_UpDown] != STATE_NEW {
						NewOrder_OnAll_IDs = false
						break
					}
				}

				if NewOrder_OnAll_IDs {
					localStateArray := localNodeOrderStates[localIP]
					localStateArray[floor][btn_UpDown] = STATE_CONFIRMED
					localNodeOrderStates[localIP] = localStateArray
					reqStateMatrixTimer.Reset(1)
					hallOrderArrayTimer.Reset(1)
					elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true) //turn on light?
				}
			}
		}
	}
}

func ConfirmedOrdersToHallOrder(allNOSMap map[string][dt.N_FLOORS][2]dt.OrderState, localIP string) [dt.N_FLOORS][2]bool {

	Local_HallOrderArray := [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	for floor := range allNOSMap[localIP] {
		for btn_UpDown := range allNOSMap[localIP][floor] {
			if allNOSMap[localIP][floor][btn_UpDown] == STATE_CONFIRMED {
				Local_HallOrderArray[floor][btn_UpDown] = true
			}
		}
	}
	return Local_HallOrderArray
}
