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
	HallOrderArray chan<- [][2]bool,
	ReqStateMatrix_toP2P chan<- dt.RequestStateMatrix,
	peerUpdateChan <-chan peers.PeerUpdate,
) {

	Local_ReqStatMatrix := make(dt.RequestStateMatrix)
	peerList := peers.PeerUpdate{}

	for {
		reqStateMatrixUpdated := false
		select {
		case peerList = <-peerUpdateChan:
			for _, nodeID := range peerList.Peers {
				if _, valInMap := Local_ReqStatMatrix[nodeID]; !valInMap {
					Local_ReqStatMatrix[nodeID] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
					reqStateMatrixUpdated = true
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

						local_state := &Local_ReqStatMatrix[localIP][floor][btn_UpDown]
						//cyclic change of states
						switch other_state {
						case STATE_none:
							if *local_state == STATE_confirmed {
								*local_state = STATE_none
								elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, false)
								reqStateMatrixUpdated = true
							}
						case STATE_new:

							if *local_state == STATE_none {
								*local_state = STATE_new
								reqStateMatrixUpdated = true
							}

						case STATE_confirmed:

							if *local_state == STATE_new {
								*local_state = STATE_confirmed
								elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true)
								reqStateMatrixUpdated = true
							}
						}
					}
				}
			}

			fmt.Printf("______RSM sent to P2P__________\n")
			fmt.Printf("Sender ID: %v\n", localIP)
			fmt.Printf("Data: %v\n", Local_ReqStatMatrix)
			fmt.Printf("_________________________\n")

		case BtnPress := <-HallBtnPress:
			//fmt.Printf("\n___ORDERSTATEHANDLER___: \n Buttnpress recieved: \n%+v\n", BtnPress)
			if Local_ReqStatMatrix[localIP][BtnPress.Floor][BtnPress.Button] == STATE_none {
				Local_ReqStatMatrix[localIP][BtnPress.Floor][BtnPress.Button] = STATE_new
				reqStateMatrixUpdated = true
			}

		case executedArray := <-orderExecuted:
			//fmt.Printf("\n___ORDERSTATEHANDLER___: \n  ExecutedArray Received \n%+v\n", executedArray)
			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}
				local_State := &Local_ReqStatMatrix[localIP][btn.Floor][btn.Button]
				// Kva skjer her dersom order blir executed før state_confirmed?
				if *local_State == STATE_confirmed {
					*local_State = STATE_none
					reqStateMatrixUpdated = true
					elevio.SetButtonLamp(btn.Button, btn.Floor, false)
				}
			}
		}

		//Check if Order can be confirmed
		//If all orders across IDs is State_new, order is confirmed and sendt to order Assigner
		for floor := range Local_ReqStatMatrix[localIP] {
			for btn_UpDown := range Local_ReqStatMatrix[localIP][floor] {

				if Local_ReqStatMatrix[localIP][floor][btn_UpDown] != STATE_new {
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
					Local_ReqStatMatrix[localIP][floor][btn_UpDown] = STATE_confirmed
					reqStateMatrixUpdated = true
					elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true) //turn on light?

					//fmt.Printf("\n___ORDERSTATEHANDLER___: \n Hallorders sendt to DataDist: \n%+v\n", ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localIpAdress))

				}
			}
		}
		if reqStateMatrixUpdated {
			ReqStateMatrix_toP2P <- Local_ReqStatMatrix
			HallOrderArray <- ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localIP)
		}
	}
}

func ConfirmedOrdersToHallOrder(requests dt.RequestStateMatrix, localID string) [][2]bool {

	Local_HallOrderArray := [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}

	for floor := range requests[localID] {
		for btn_UpDown := range requests[localID][floor] {
			if requests[localID][floor][btn_UpDown] == STATE_confirmed {
				Local_HallOrderArray[floor][btn_UpDown] = true
			}
		}
	}
	return Local_HallOrderArray
}
