package orderStateHandler

import (
	"Driver-go/elevio"
)

//
localIP := "Local IP adress" 


//Datatypes
type requestState int

type singleNode_requestStates [][2]requestState

type requestStateMatrix map[string]singleNode_requestStates

type RequestStateMatrix_with_ID struct{
	IpAdress 		string 				`json:"ipAdress"`
	RequestMatrix 	requestStateMatrix  `json:"requestMatrix"`
}


// States for hall requests 
const (
	STATE_new       	requestState = 0
	STATE_confirmed 	requestState = 1
	STATE_none 			requestState = 2
)



// input channels
var (
	ReqStateMatrix_fromP2P 	= make(chan RequestStateMatrix_with_ID)
	HallBtnPress            = make(chan elevio.ButtonEvent)
	orderExecuted       	= make(chan [][2]int)
)

// output channels
var (
	HallOrderArray     		= make(chan [][2]bool)
	ReqStateMatrix_toP2P 	= make(chan requestStateMatrix)
)



func orderStateHandler(
	ReqStateMatrix_fromP2P <-chan requestStateMatrix,
	HallBtnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan [][2]int,
	HallOrderArray chan<- StateOfWorldView,
	ReqStateMatrix_toP2P chan<- WorldView,

) {
	//init local RequestStateMatrix
	Local_ReqStatMatrix := { 
	"ID1": singleNode_requestStates{{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none}}, 	 
	"ID2": singleNode_requestStates{{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none}},
	"ID3": singleNode_requestStates{{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none},{STATE_none, STATE_none}},
	}

	//init local hallOrderArray
	Local_HallOrderArray := [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}}


	for{
		select{
		case matrix_fromP2P := <-chan ReqStateMatrix_fromP2P:

		case BtnPress := <-chan elevio.HallBtnPress:

		case executedArray <-chan orderExecuted:

		}
	}



}






