package oassign

import (
	"fmt"
	elevfsm "project/localElevator/elev_fsm"
)

func OrderAssignerTestFunc() {
	var (
		OrderAssignerBehaviourChan = make(chan OrderAssignerBehaviour)
		localIpAdressChan          = make(chan string)   // Chanel where local IP-adress is fetched
		ordersFromDistributor      = make(chan HRAInput) // Input from order distributor
		ordersFromMaster           = make(chan []byte)   // Input read from Master-Slave network module
		ordersToSlaves             = make(chan []byte)   // Input written to Master-Slave network module
		localOrders                = make(chan [elevfsm.N_FLOORS][2]bool)
	)
	go OrderAssigner(OrderAssignerBehaviourChan, localIpAdressChan, ordersFromDistributor, ordersFromMaster, ordersToSlaves, localOrders)

	for {
		OrderAssignerBehaviourChan <- MS_Master
		localIpAdressChan <- "127.0.0.1"

		input := HRAInput{
			HallRequests: [elevfsm.N_FLOORS][2]bool{{true, true}, {true, true}, {true, true}, {false, true}},
			States: map[string]HRAElevState{
				"127.0.0.1": {
					Behavior:    "idle",
					Floor:       2,
					Direction:   "up",
					CabRequests: [elevfsm.N_FLOORS]bool{true, false, true, true},
				},
				"127.0.0.2": {
					Behavior:    "idle",
					Floor:       0,
					Direction:   "stop",
					CabRequests: [elevfsm.N_FLOORS]bool{false, true, false, false},
				},
				"127.0.0.3": {
					Behavior:    "idle",
					Floor:       1,
					Direction:   "stop",
					CabRequests: [elevfsm.N_FLOORS]bool{true, false, false, false},
				},
			},
		}

		ordersFromDistributor <- input

		local := <-localOrders
		for i := 0; i < elevfsm.N_FLOORS; i++ {
			fmt.Printf("Local orders at floor %d: {%t, %t} \n", i, local[i][0], local[i][1])
		}
	}
}
