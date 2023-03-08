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
	STATE_none      requestState = 0
	STATE_new       requestState = 1
	STATE_confirmed requestState = 2
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
