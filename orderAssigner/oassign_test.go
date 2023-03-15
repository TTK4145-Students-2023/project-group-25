package oassign

import (
	"fmt"
	"project/Network/Utilities/localip"
	dt "project/commonDataTypes"
)

func OrderAssignerTestFunc() {
	var (
		masterSlaveRoleChan = make(chan dt.MasterSlaveRole)
		localIpAdressChan   = make(chan string) // Input with local IP-adress

		ordersFromDistributor = make(chan dt.CostFuncInput)                // Input from order distributor
		ordersFromMaster      = make(chan map[string][dt.N_FLOORS][2]bool) // Input from Master-Slave network module
		ordersToSlaves        = make(chan map[string][dt.N_FLOORS][2]bool) // Output to Master-Slave network module
		ordersLocal           = make(chan [dt.N_FLOORS][2]bool)            // Output to local elevator
	)
	localIP, _ := localip.LocalIP()
	go OrderAssigner(localIP,
		masterSlaveRoleChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	for {
		masterSlaveRoleChan <- dt.MS_Master
		localIpAdressChan <- "127.0.0.1"

		input := dt.CostFuncInput{
			HallRequests: [dt.N_FLOORS][2]bool{{true, true}, {true, true}, {true, true}, {false, true}},
			States: map[string]dt.ElevDataJSON{
				"127.0.0.1": {
					Behavior:    "idle",
					Floor:       2,
					Direction:   "up",
					CabRequests: [dt.N_FLOORS]bool{true, false, true, true},
				},
				"127.0.0.2": {
					Behavior:    "idle",
					Floor:       0,
					Direction:   "stop",
					CabRequests: [dt.N_FLOORS]bool{false, true, false, false},
				},
				"127.0.0.3": {
					Behavior:    "idle",
					Floor:       1,
					Direction:   "stop",
					CabRequests: [dt.N_FLOORS]bool{true, false, false, false},
				},
			},
		}

		ordersFromDistributor <- input

		local := <-ordersLocal
		for i := range local {
			fmt.Printf("Local orders at floor %d: {%t, %t} \n", i, local[i][0], local[i][1])
		}
	}
}
