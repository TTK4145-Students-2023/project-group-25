package orderStateHandler

import (
	"fmt"
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
)

var localID string = "127.0.0.1"

// States for hall requests
const (
	STATE_none      dt.RequestState = 0
	STATE_new       dt.RequestState = 1
	STATE_confirmed dt.RequestState = 2
)

// // input channels
// var (
// 	ReqStateMatrix_fromP2P = make(chan dt.RequestStateMatrix_with_ID)
// 	HallBtnPress           = make(chan elevio.ButtonEvent)
// )

// // output channels
// var (
// 	HallOrderArray       = make(chan [][2]bool)
// 	ReqStateMatrix_toP2P = make(chan dt.RequestStateMatrix)
// )

func OrderStateHandler(
	ReqStateMatrix_fromP2P <-chan dt.RequestStateMatrix,
	HallBtnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan []elevio.ButtonEvent,
	HallOrderArray chan<- [][2]bool,
	ReqStateMatrix_toP2P chan<- dt.RequestStateMatrix,

) {
	//init local RequestStateMatrix
	Local_ReqStatMatrix := make(dt.RequestStateMatrix)
	Local_ReqStatMatrix[localID] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	// Local_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	// Local_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}

	// List of node IDs we are connected to
	nodeIDs := []string{localID} //, "ID2", "ID3"}

	for {
		select {
		case matrix_fromP2P := <-ReqStateMatrix_fromP2P:

			fmt.Printf("\n___ORDERSTATEHANDLER___ \n Input recived from P2P \n%+v\n", matrix_fromP2P)

			// Iterate through the list of node IDs
			for _, nodeID := range nodeIDs {
				// Skip the local node
				if nodeID == localID {
					continue
				}
				// Compare the requestStates from the other nodes with the Local requestStates
				for floor := range matrix_fromP2P[nodeID] {
					for btn_UpDown, other_state := range matrix_fromP2P[nodeID][floor] {

						local_state := &Local_ReqStatMatrix[localID][floor][btn_UpDown]

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

			ReqStateMatrix_toP2P <- Local_ReqStatMatrix

		case BtnPress := <-HallBtnPress:
			fmt.Printf("\n___ORDERSTATEHANDLER___: \n Buttnpress recieved: \n%+v\n", BtnPress)
			Local_ReqStatMatrix[localID][BtnPress.Floor][BtnPress.Button] = STATE_new

		case executedArray := <-orderExecuted:
			fmt.Printf("\n___ORDERSTATEHANDLER___: \n  ExecutedArray Received \n%+v\n", executedArray)

			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}

				local_State := Local_ReqStatMatrix[localID][btn.Floor][btn.Button]
				if local_State == STATE_confirmed {
					Local_ReqStatMatrix[localID][btn.Floor][btn.Button] = STATE_none
					elevio.SetButtonLamp(btn.Button, btn.Floor, false) //turn off light?
				}

			}

		}
		//Check if Order can be confirmed
		// If all orders across IDs is State_new, order is confirmed and sendt to order Assigner
		for floor := range Local_ReqStatMatrix[localID] {
			for btn_UpDown := range Local_ReqStatMatrix[localID][floor] {

				if Local_ReqStatMatrix[localID][floor][btn_UpDown] != STATE_new {
					continue
				}

				NewOrder_OnAll_IDs := true
				for _, nodeID := range nodeIDs {
					if Local_ReqStatMatrix[nodeID][floor][btn_UpDown] != STATE_new {
						NewOrder_OnAll_IDs = false
						break
					}
				}

				if NewOrder_OnAll_IDs {
					Local_ReqStatMatrix[localID][floor][btn_UpDown] = STATE_confirmed
					elevio.SetButtonLamp(elevio.ButtonType(btn_UpDown), floor, true) //turn on light?
					HallOrderArray <- ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localID)
					fmt.Printf("\n___ORDERSTATEHANDLER___: \n Hallorders sendt to DataDist: \n%+v\n", ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localID))

				}
			}
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
