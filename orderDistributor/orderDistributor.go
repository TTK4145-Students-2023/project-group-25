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
	STATE_updateData         DistributorState = 0
	STATE_distributeBTNPress DistributorState = 1
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

	distributor_state := STATE_updateData
	worldView := []int
	for {
		select {
		case DataFromP2P := <-allElevData_fromP2P:

			worldView = update_worldView(DataFromP2P, localElevData)

			switch distributor_state {
			case STATE_updateData:

				allElevData_toAssigner <- worldView
				allElevData_toP2P <- worldView

			case STATE_distributeBTNPress:

				if orderDistributed(DataFromP2P, localElevData) == true {
					allElevData_toAssigner <- worldView
					allElevData_toP2P <- worldView
					//lights on
					distributor_state = STATE_updateData

				} else {
					allElevData_toP2P <- worldView
				}

			}

		case executedOrder := <-orderExecuted:

			deleteOrder(executedOrder)

			switch distributor_state {
			case STATE_updateData:
				allElevData_toP2P <- worldView

				//allElevData_toAssigner <- worldView   //not neccesary?
				//lights out??

			case STATE_distributeBTNPress:
				allElevData_toP2P <- worldView

				//allElevData_toAssigner <- worldView   //not neccesary?
				//lights out??
			}

		case pressedBtn := <-btnPress:

			addOrder(pressedBtn)

			switch distributor_state {
			case STATE_updateData:

				allElevData_toP2P <- worldView
				distributor_state = STATE_distributeBTNPress

			case STATE_distributeBTNPress:
				// dont accept New buttonpresses when distrubiting

			}
		}
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func update_worldView(x ...interface{}) int

func orderDistributed(x ...interface{}) bool

func deleteOrder(x ...interface{})

func addOrder(x ...interface{})
