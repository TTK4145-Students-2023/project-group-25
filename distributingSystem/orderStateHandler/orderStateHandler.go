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

	//init local RequestStateMatrix
	Local_ReqStatMatrix := make(dt.RequestStateMatrix)
	// Local_ReqStatMatrix[localIpAdress] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}

	// List of node IDs we are connected to
	peerList := peers.PeerUpdate{} // peerList := []string{localIpAdress}

	for {
		select {
		case peerList = <-peerUpdateChan:
			// Initilize new nodes
			for _, nodeID := range peerList.Peers {
				if _, valInMap := Local_ReqStatMatrix[nodeID]; !valInMap {
					Local_ReqStatMatrix[nodeID] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
					fmt.Printf("We are here in PeerList update init\n")
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
							}
						case STATE_new:

							if *local_state == STATE_none {
								*local_state = STATE_new
							}

						case STATE_confirmed:

							if *local_state == STATE_new {
								*local_state = STATE_confirmed
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
			Local_ReqStatMatrix[localIP][BtnPress.Floor][BtnPress.Button] = STATE_new

		case executedArray := <-orderExecuted:
			//fmt.Printf("\n___ORDERSTATEHANDLER___: \n  ExecutedArray Received \n%+v\n", executedArray)

			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}

				local_State := Local_ReqStatMatrix[localIP][btn.Floor][btn.Button]
				if local_State == STATE_confirmed {
					Local_ReqStatMatrix[localIP][btn.Floor][btn.Button] = STATE_none
					elevio.SetButtonLamp(btn.Button, btn.Floor, false) //turn off light?
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
					elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true) //turn on light?
					HallOrderArray <- ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localIP)
					//fmt.Printf("\n___ORDERSTATEHANDLER___: \n Hallorders sendt to DataDist: \n%+v\n", ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localIpAdress))

				}
			}
		}

		// Send updated Reqmatrix to P2P
		ReqStateMatrix_toP2P <- Local_ReqStatMatrix
		// fmt.Printf("______RSM sent to P2P__________\n")
		// fmt.Printf("Sender ID: %v\n", localIP)
		// fmt.Printf("Data: %v\n", Local_ReqStatMatrix)
		// fmt.Printf("_________________________\n")

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
