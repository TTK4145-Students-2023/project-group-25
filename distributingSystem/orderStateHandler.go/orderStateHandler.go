package orderStateHandler

import (
	"Driver-go/elevio"
)

var localIP string = "Local IP adress"

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
	orderExecuted          = make(chan [][2]int)
)

// output channels
var (
	HallOrderArray       = make(chan [][2]bool)
	ReqStateMatrix_toP2P = make(chan RequestStateMatrix)
)

func orderStateHandler(
	ReqStateMatrix_fromP2P <-chan RequestStateMatrix,
	HallBtnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan [][2]int,
	HallOrderArray chan<- [][2]bool,
	ReqStateMatrix_toP2P chan<- RequestStateMatrix,

) {
	//init local RequestStateMatrix
	Local_ReqStatMatrix := make(RequestStateMatrix)
	Local_ReqStatMatrix["ID1"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	Local_ReqStatMatrix["ID2"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}
	Local_ReqStatMatrix["ID3"] = singleNode_requestStates{{STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}, {STATE_none, STATE_none}}

	//init local hallOrderArray
	Local_HallOrderArray := [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}}
	UNUSED(Local_HallOrderArray)

	for {
		select {
		case matrix_fromP2P := <-ReqStateMatrix_fromP2P:
			UNUSED(matrix_fromP2P)

		case BtnPress := <-HallBtnPress:
			UNUSED(BtnPress)

		case executedArray := <-orderExecuted:
			UNUSED(executedArray)

		}
	}

}

func UNUSED(x ...interface{}) {}
