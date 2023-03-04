package orderStateHandler

import (
	"Driver-go/elevio"
)

var localID string = "ID1"

// Datatypes
type requestState int

type singleNode_requestStates [][2]requestState

type RequestStateMatrix map[string]singleNode_requestStates

type RequestStateMatrix_with_ID struct {
	IpAdress      string             `json:"ipAdress"`
	RequestMatrix RequestStateMatrix `json:"requestMatrix"`
}

// States for hall requests
const (
	STATE_new       requestState = 0
	STATE_confirmed requestState = 1
	STATE_none      requestState = 2
)

// input channels
var (
	ReqStateMatrix_fromP2P = make(chan RequestStateMatrix_with_ID)
	HallBtnPress           = make(chan elevio.ButtonEvent)
	orderExecuted          = make(chan []elevio.ButtonEvent)
)

// output channels
var (
	HallOrderArray       = make(chan [][2]bool)
	ReqStateMatrix_toP2P = make(chan RequestStateMatrix)
)

func orderStateHandler(
	ReqStateMatrix_fromP2P <-chan RequestStateMatrix,
	HallBtnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan []elevio.ButtonEvent,
	HallOrderArray chan<- [][2]bool,
	ReqStateMatrix_toP2P chan<- RequestStateMatrix,

) {
	//init local RequestStateMatrix
	Local_ReqStatMatrix := make(RequestStateMatrix)
	Local_ReqStatMatrix["ID1"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	Local_ReqStatMatrix["ID2"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	Local_ReqStatMatrix["ID3"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}

	for {
		select {
		case matrix_fromP2P := <-ReqStateMatrix_fromP2P:

			// List of node IDs we are connected to
			nodeIDs := []string{"ID1", "ID2", "ID3"}

			// Iterate through the list of node IDs
			for _, nodeID := range nodeIDs {
				// Skip the local node
				if nodeID == localID {
					continue
				}
				// Compare the requestStates from the other nodes with the Local requestStates
				for floor := range matrix_fromP2P[nodeID] {
					for btn_UpDown, other_state := range matrix_fromP2P[nodeID][floor] {

						local_state := Local_ReqStatMatrix[localID][floor][btn_UpDown]

						//cyclic change of states
						switch other_state {
						case STATE_none:
							if local_state == STATE_confirmed {
								Local_ReqStatMatrix[localID][floor][btn_UpDown] = STATE_none
							}
						case STATE_new:
							if local_state == STATE_none {
								Local_ReqStatMatrix[localID][floor][btn_UpDown] = STATE_new
							}
						case STATE_confirmed:
							if local_state == STATE_new {
								Local_ReqStatMatrix[localID][floor][btn_UpDown] = STATE_confirmed
							}
						}
					}
				}
			}

			//Check if Order can be confirmed
			// If all orders across IDs is State_new, order is confirmed
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
						//sett light
					}
				}
			}

			ReqStateMatrix_toP2P <- Local_ReqStatMatrix
			HallOrderArray <- ConfirmedOrdersToHallOrder(Local_ReqStatMatrix, localID)

		case BtnPress := <-HallBtnPress:
			Local_ReqStatMatrix[localID][BtnPress.Floor][BtnPress.Button] = STATE_new

		case executedArray := <-orderExecuted:
			for _, btn := range executedArray {
				if btn.Button == elevio.BT_Cab {
					continue
				}

				local_State := Local_ReqStatMatrix[localID][btn.Floor][btn.Button]
				if local_State == STATE_confirmed {
					Local_ReqStatMatrix[localID][btn.Floor][btn.Button] = STATE_none
					//turn off light
				}

			}

		}
	}
}

func ConfirmedOrdersToHallOrder(requests RequestStateMatrix, localID string) [][2]bool {

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
