package orderStateHandler

import (
	"fmt"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
)

// States for hall requests
const (
	STATE_none      dt.RequestState = 0
	STATE_new       dt.RequestState = 1
	STATE_confirmed dt.RequestState = 2
)

func OrderStateHandler(localIP string,
	ReqStateMatrix_fromP2P <-chan dt.RequestStateMatrix_with_ID,
	HallBtnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan []elevio.ButtonEvent,
	HallOrderArray chan<- [dt.N_FLOORS][2]bool,
	ReqStateMatrix_toP2P chan<- dt.RequestStateMatrix,
	peerUpdateChan <-chan peers.PeerUpdate,
) {

	Local_ReqStatMatrix := make(dt.RequestStateMatrix)
	peerList := peers.PeerUpdate{}

	for {
		select {
		case peerList = <-peerUpdateChan:
			for _, nodeID := range peerList.Peers {
				if _, valInMap := Local_ReqStatMatrix[nodeID]; !valInMap {
					Local_ReqStatMatrix[nodeID] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
				}
			}
		case matrix_fromP2P := <-ReqStateMatrix_fromP2P:
			// update external states based on sender ID
			Local_ReqStatMatrix[matrix_fromP2P.IpAdress] = matrix_fromP2P.RequestMatrix[matrix_fromP2P.IpAdress]
			// Iterate through the list of node IDs
			for _, nodeID := range peerList.Peers {
				// Skip the local node
				if nodeID == localIP {
					continue
				}
				// Compare the requestStates from the other nodes with the Local requestStates
				for floor := range matrix_fromP2P.RequestMatrix[nodeID] {
					for btn_UpDown, other_state := range matrix_fromP2P.RequestMatrix[nodeID][floor] {

						localStateArray := Local_ReqStatMatrix[localIP]
						//cyclic change of states
						switch other_state {
						case STATE_none:
							if localStateArray[floor][btn_UpDown] == STATE_confirmed {
								localStateArray[floor][btn_UpDown] = STATE_none
								Local_ReqStatMatrix[localIP] = localStateArray
							}
						case STATE_new:
							if localStateArray[floor][btn_UpDown] == STATE_none {
								localStateArray[floor][btn_UpDown] = STATE_new
								Local_ReqStatMatrix[localIP] = localStateArray
							}
						case STATE_confirmed:
							if localStateArray[floor][btn_UpDown] == STATE_new {
								localStateArray[floor][btn_UpDown] = STATE_confirmed
								Local_ReqStatMatrix[localIP] = localStateArray
							}
						}
					}
				}
			}
		case BtnPress := <-HallBtnPress:
			localStateArray := Local_ReqStatMatrix[localIP]
			if localStateArray[BtnPress.Floor][BtnPress.Button] == STATE_none {
				localStateArray[BtnPress.Floor][BtnPress.Button] = STATE_new
				Local_ReqStatMatrix[localIP] = localStateArray
			}
		case executedArray := <-orderExecuted:
			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}
				localStateArray := Local_ReqStatMatrix[localIP]
				if localStateArray[btn.Floor][btn.Button] == STATE_confirmed {
					localStateArray[btn.Floor][btn.Button] = STATE_none
					Local_ReqStatMatrix[localIP] = localStateArray
				}
			}
		}
		//Check if Order can be confirmed
		//If all orders across IDs is State_new, order is confirmed and sendt to order Assigner
		for floor, floorStateArray := range Local_ReqStatMatrix[localIP] {
			for btn_UpDown := range floorStateArray {

				if floorStateArray[btn_UpDown] != STATE_new {
					continue
				}

				NewOrder_OnAll_IDs := true
				for _, nodeID := range peerList.Peers {
					if Local_ReqStatMatrix[nodeID][floor][btn_UpDown] != STATE_new {
						NewOrder_OnAll_IDs = false
						break
					}
				}

				if NewOrder_OnAll_IDs {
					localStateArray := Local_ReqStatMatrix[localIP]
					localStateArray[floor][btn_UpDown] = STATE_confirmed
					Local_ReqStatMatrix[localIP] = localStateArray
					elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true) //turn on light?
				}
			}
		}
		if reqStateMatrixUpdated {
			ReqStateMatrix_toP2P <- Local_ReqStatMatrix
			HallOrderArray <- ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localIP)
		}
	}
}

func ConfirmedOrdersToHallOrder(requests dt.RequestStateMatrix, localID string) [dt.N_FLOORS][2]bool {

	Local_HallOrderArray := [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	for floor := range requests[localID] {
		for btn_UpDown := range requests[localID][floor] {
			if requests[localID][floor][btn_UpDown] == STATE_confirmed {
				Local_HallOrderArray[floor][btn_UpDown] = true
			}
		}
	}
	return Local_HallOrderArray
}
