package orderdistributor

import (
	"Driver-go/elevio"
)

var (
	latest_ordermatrix []int
)

type DistributorState int

type ElevData struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

const (
	STATE_updateLocalData        DistributorState = 0
	STATE_distributeLocalChanges DistributorState = 1
)

// input channels
var (
	allElevData_fromP2P = make(chan ElevData)
	btnPress            = make(chan elevio.ButtonEvent)
	orderExecuted       = make(chan bool)
	localElevData       = make(chan ElevData)
)

// output channels
var (
	allElevData_toP2P      = make(chan ElevData)
	allElevData_toAssigner = make(chan ElevData)
)

func dataDistributor(
	allElevData_fromP2P <-chan ElevData,
	btnPress <-chan elevio.ButtonEvent,
	orderExecuted <-chan bool,
	localElevData <-chan ElevData, //not used as fsm trigger
	allElevData_toP2P chan<- ElevData,
	allElevData_toAssigner chan<- ElevData,

) {

	for {
		distributor_state := STATE_updateLocalData
		select {
		case allElevData := <-allElevData_fromP2P:
			switch distributor_state {
			case STATE_updateLocalData:

				//update local storage of data and send json to assigner

				//latestValid_ordermatrix, _ = validateData(allElevData)
				//allElevData_toAssigner <- latestValid_ordermatrix.json()

				UNUSED(allElevData)

			case STATE_distributeLocalChanges:
				// check if input==output
				// if true: switch state and turn on/off light
				//if false: send change again

				//_, succesfully_distributed := validateData(allElevData)

			}

		case executedOrder := <-orderExecuted:
			switch distributor_state {
			case STATE_updateLocalData:
				//add order change to elevdata and send to p2p
				//switch state
				UNUSED(executedOrder)

			case STATE_distributeLocalChanges:
				// dont accept orderchanges when distrubiting

			}

		case pressedBtn := <-btnPress:
			switch distributor_state {
			case STATE_updateLocalData:
				//add order change to elevdata and send to p2p
				//switch state

				UNUSED(pressedBtn)

			case STATE_distributeLocalChanges:
				// dont accept orderchanges when distrubiting

			}
		}
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
