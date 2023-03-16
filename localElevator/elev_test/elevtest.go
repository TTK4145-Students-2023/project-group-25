package elevtest

// import (
// 	"fmt"
// 	"project/Network/Utilities/localip"
// 	btnassign "project/buttonAssigner"
// 	dt "project/commonDataTypes"
// 	elevio "project/localElevator/elev_driver"
// 	elevfsm "project/localElevator/elev_fsm"
// 	oassign "project/orderAssigner"
// 	"time"
// )

// func intermediateOrderDistributor(
// 	hallEvent <-chan elevio.ButtonEvent,
// 	handler_hallOrdersExecuted <-chan []elevio.ButtonEvent,
// 	elev_data <-chan dt.ElevDataJSON,
// 	ordersFromDistributor chan<- dt.CostFuncInput) {

// 	orderOverview := dt.CostFuncInput{
// 		HallRequests: [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
// 		States: map[string]dt.ElevDataJSON{
// 			"127.0.0.1": {
// 				Behavior:    "idle",
// 				Floor:       2,
// 				Direction:   "stop",
// 				CabRequests: [dt.N_FLOORS]bool{false, false, false, false},
// 			},
// 		},
// 	}

// 	for {
// 		select {
// 		case buttonEvent := <-hallEvent:
// 			fmt.Printf("HallEvent.\n")
// 			orderOverview.HallRequests[buttonEvent.Floor][buttonEvent.Button] = true
// 			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)

// 		case hallOrdersExecuted := <-handler_hallOrdersExecuted:
// 			fmt.Printf("OrderExecute.\n")
// 			for i := range hallOrdersExecuted {
// 				if hallOrdersExecuted[i].Button != elevio.BT_Cab {
// 					orderOverview.HallRequests[hallOrdersExecuted[i].Floor][hallOrdersExecuted[i].Button] = false
// 					elevio.SetButtonLamp(hallOrdersExecuted[i].Button, hallOrdersExecuted[i].Floor, false)
// 				}
// 			}

// 		case elevatorData := <-elev_data:
// 			fmt.Printf("newElevData.\n")
// 			fmt.Printf("ElevatorData:\n Beh: %s\n Dir: %s\n Floor: %d\n\n", elevatorData.Behavior, elevatorData.Direction, elevatorData.Floor)
// 			orderOverview.States["127.0.0.1"] = elevatorData
// 		}
// 		ordersFromDistributor <- orderOverview
// 		fmt.Printf("Data from distro sent!\n")
// 	}
// }

// func testSpecDistributor(OrderAssignerBehaviourChan chan dt.MasterSlaveRole) {
// 	orderAssignerBehaviour := dt.MS_Master

// 	for {
// 		OrderAssignerBehaviourChan <- orderAssignerBehaviour
// 	}
// }

// var (
// 	OrderAssignerBehaviourChan = make(chan dt.MasterSlaveRole)

// 	ordersFromDistributor      = make(chan dt.CostFuncInput)                // Input from order distributor
// 	ordersFromMaster           = make(chan map[string][dt.N_FLOORS][2]bool) // Input read from Master-Slave network module
// 	ordersToSlaves             = make(chan map[string][dt.N_FLOORS][2]bool) // Input written to Master-Slave network module
// 	ordersLocal                = make(chan [dt.N_FLOORS][2]bool)
// 	handler_hallOrdersExecuted = make(chan []elevio.ButtonEvent)

// 	btnEvent  = make(chan elevio.ButtonEvent)
// 	hallEvent = make(chan elevio.ButtonEvent)
// 	cabEvent  = make(chan elevio.ButtonEvent)

// 	drv_floors = make(chan int)
// 	drv_obstr  = make(chan bool)
// 	elev_data  = make(chan dt.ElevDataJSON)
// )

// func RunSingleElevTest() {
// 	localIP, _ := localip.LocalIP()
// 	elevio.Init("localhost:15657", dt.N_FLOORS)
// 	go elevio.PollFloorSensor(drv_floors)
// 	go elevio.PollButtons(btnEvent)
// 	go elevio.PollObstructionSwitch(drv_obstr)

// 	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

// 	go intermediateOrderDistributor(hallEvent,
// 		handler_hallOrdersExecuted,
// 		elev_data,
// 		ordersFromDistributor)

// 	go testSpecDistributor(OrderAssignerBehaviourChan)
// 	go oassign.OrderAssigner(localIP,
// 		OrderAssignerBehaviourChan,
// 		ordersFromDistributor,
// 		ordersFromMaster,
// 		ordersToSlaves,
// 		ordersLocal)

// 	time.Sleep(time.Millisecond * 40)

// 	go elevfsm.FSM(ordersLocal,
// 		cabEvent,
// 		drv_floors,
// 		drv_obstr,
// 		elev_data,
// 		handler_hallOrdersExecuted)

// 	for {
// 		// Kill orders that are yet to be handeled!
// 		<-ordersToSlaves
// 	}
// }
