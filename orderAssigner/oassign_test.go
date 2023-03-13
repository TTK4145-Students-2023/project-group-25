package oassign

import (
	"fmt"
	dt "project/commonDataTypes"
	elevfsm "project/localElevator/elev_fsm"
)

func OrderAssignerTestFunc() {
	var (
		masterSlaveRoleChan = make(chan dt.MasterSlaveRole)
		localIpAdressChan   = make(chan string) // Input with local IP-adress

		ordersFromDistributor = make(chan dt.CostFuncInput) // Input from order distributor
		ordersFromMaster      = make(chan []byte)           // Input from Master-Slave network module
		ordersToSlaves        = make(chan []byte)           // Output to Master-Slave network module
		ordersLocal           = make(chan [][2]bool)        // Output to local elevator
	)
	go OrderAssigner(masterSlaveRoleChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	for {
		masterSlaveRoleChan <- dt.MS_Master
		localIpAdressChan <- "127.0.0.1"

		input := dt.CostFuncInput{
			HallRequests: [][2]bool{{true, true}, {true, true}, {true, true}, {false, true}},
			States: map[string]dt.ElevDataJSON{
				"127.0.0.1": {
					Behavior:    "idle",
					Floor:       2,
					Direction:   "up",
					CabRequests: []bool{true, false, true, true},
				},
				"127.0.0.2": {
					Behavior:    "idle",
					Floor:       0,
					Direction:   "stop",
					CabRequests: []bool{false, true, false, false},
				},
				"127.0.0.3": {
					Behavior:    "idle",
					Floor:       1,
					Direction:   "stop",
					CabRequests: []bool{true, false, false, false},
				},
			},
		}

		ordersFromDistributor <- input

		local := <-ordersLocal
		for i := 0; i < elevfsm.N_FLOORS; i++ {
			fmt.Printf("Local orders at floor %d: {%t, %t} \n", i, local[i][0], local[i][1])
		}
	}
}
