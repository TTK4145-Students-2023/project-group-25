package orderdistributor

import (
	"Driver-go/elevio"
)

// Datatypes
type ElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type WorldView struct {
	HallRequests [][2]bool            `json:"hallRequests"`
	States       map[string]ElevState `json:"elevStates"`
}

type StateOfWorldView struct {
	CurrentWorldView   WorldView                           `json:"currentWorldView"`
	RequestStateMatrix map[string]SingleNode_RequestStates `json:"requestStateMatrix"`
}

type SingleNode_RequestStates struct {
	Requests [][2]RequestState `json:"requests"`
}

// States
type RequestState int

const (
	Req_STATE_none      RequestState = 0
	Req_STATE_new       RequestState = 1
	Req_STATE_confirmed RequestState = 2
)

type DistributorState int

const (
	disb_STATE_updateWorldView    DistributorState = 0
	disb_STATE_distributeBTNPress DistributorState = 1
)

// input channels
var (
	allElevData_fromP2P = make(chan StateOfWorldView)
	btnPress            = make(chan elevio.ButtonEvent)
	orderExecuted       = make(chan bool)
	localElevData       = make(chan ElevState)
)

// output channels
var (
	allElevData_toP2P      = make(chan StateOfWorldView)
	allElevData_toAssigner = make(chan WorldView)
)

// Statemachine for Distributor
func dataDistributor(
	allElevData_fromP2P <-chan StateOfWorldView,
	btnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan bool,
	localElevData <-chan ElevState, //not used as fsm trigger
	allElevData_toP2P chan<- StateOfWorldView,
	allElevData_toAssigner chan<- WorldView,

) {

	distributor_state := disb_STATE_updateWorldView

	current_stateOfWorldView := StateOfWorldView{
		CurrentWorldView:   WorldView{},
		RequestStateMatrix: [][2]RequestState{},
	}

	for {
		select {
		case DataFromP2P := <-allElevData_fromP2P:

			current_stateOfWorldView = update_worldView(DataFromP2P, localElevData)

			switch distributor_state {
			case disb_STATE_updateWorldView:

				allElevData_toAssigner <- current_stateOfWorldView.CurrentWorldView
				allElevData_toP2P <- current_stateOfWorldView

			case disb_STATE_distributeBTNPress:

				if orderDistributed(DataFromP2P, localElevData) == true {
					allElevData_toAssigner <- current_stateOfWorldView.CurrentWorldView
					allElevData_toP2P <- current_stateOfWorldView
					//elevio.SetButtonLamp(CORRECT LAMP)

					distributor_state = disb_STATE_updateWorldView

				} else {
					allElevData_toP2P <- current_stateOfWorldView
				}

			}

		case executedOrder := <-orderExecuted:

			current_stateOfWorldView = deleteOrder(executedOrder)

			switch distributor_state {
			case disb_STATE_updateWorldView:
				allElevData_toP2P <- current_stateOfWorldView

				//elevio.SetButtonLamp(CORRECT LAMP)

			case disb_STATE_distributeBTNPress:
				allElevData_toP2P <- current_stateOfWorldView

				//elevio.SetButtonLamp(CORRECT LAMP)
			}

		case pressedBtn := <-btnPress:

			current_stateOfWorldView = addOrder(pressedBtn)

			switch distributor_state {
			case disb_STATE_updateWorldView:

				allElevData_toP2P <- current_stateOfWorldView
				distributor_state = disb_STATE_distributeBTNPress

			case disb_STATE_distributeBTNPress:
				// dont accept New buttonpresses when distrubiting

			}
		}
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func update_worldView(x ...interface{}) StateOfWorldView

func orderDistributed(x ...interface{}) bool

func deleteOrder(x ...interface{}) StateOfWorldView

func addOrder(x ...interface{}) StateOfWorldView
