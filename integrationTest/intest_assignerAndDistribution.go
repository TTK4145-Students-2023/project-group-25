package intest

import (
	"fmt"
	btnassign "project/buttonAssigner"
	dt "project/commonDataTypes"
	elevDataDistributor "project/distributingSystem/elevDataDistributor"
	orderStateHandler "project/distributingSystem/orderStateHandler"
	elevio "project/localElevator/elev_driver"
	elevfsm "project/localElevator/elev_fsm"
	oassign "project/orderAssigner"
	"time"
)

func intermediateOrderDistributor(hallEvent chan elevio.ButtonEvent,
	handler_hallOrdersExecuted chan []elevio.ButtonEvent,
	elev_data chan dt.ElevDataJSON,
	ordersFromDistributor chan dt.CostFuncInput) {

	orderOverview := dt.CostFuncInput{
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]dt.ElevDataJSON{
			"127.0.0.1": {
				Behavior:    "idle",
				Floor:       2,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, false},
			},
		},
	}

	for {
		select {
		case buttonEvent := <-hallEvent:
			fmt.Printf("HallEvent.\n")
			orderOverview.HallRequests[buttonEvent.Floor][buttonEvent.Button] = true
			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)

		case hallOrdersExecuted := <-handler_hallOrdersExecuted:
			fmt.Printf("OrderExecute.\n")
			for i := 0; i < len(hallOrdersExecuted); i++ {
				if hallOrdersExecuted[i].Button != elevio.BT_Cab {
					orderOverview.HallRequests[hallOrdersExecuted[i].Floor][hallOrdersExecuted[i].Button] = false
					elevio.SetButtonLamp(hallOrdersExecuted[i].Button, hallOrdersExecuted[i].Floor, false)
				}
			}

		case elevatorData := <-elev_data:
			fmt.Printf("newElevData.\n")
			orderOverview.States["127.0.0.1"] = elevatorData
		}
		ordersFromDistributor <- orderOverview
		fmt.Printf("Data from distro sent!\n")
	}
}

func testSpecDistributor(OrderAssignerBehaviourChan chan dt.OrderAssignerBehaviour,
	localIpAdressChan chan string) {
	orderAssignerBehaviour := dt.OA_Master
	localIpAdress := "127.0.0.1"

	for {
		select {
		case localIpAdressChan <- localIpAdress:
		case OrderAssignerBehaviourChan <- orderAssignerBehaviour:
		}
	}
}

var (
	OrderAssignerBehaviourChan = make(chan dt.OrderAssignerBehaviour)
	localIpAdressChan          = make(chan string) // Chanel where local IP-adress is fetched

	ordersFromDistributor      = make(chan dt.CostFuncInput) // Input from order distributor
	ordersFromMaster           = make(chan []byte)           // Input read from Master-Slave network module
	ordersToSlaves             = make(chan []byte)           // Input written to Master-Slave network module
	ordersLocal                = make(chan [][2]bool)
	handler_hallOrdersExecuted = make(chan []elevio.ButtonEvent)

	btnEvent  = make(chan elevio.ButtonEvent)
	hallEvent = make(chan elevio.ButtonEvent)
	cabEvent  = make(chan elevio.ButtonEvent)

	drv_floors = make(chan int)
	drv_obstr  = make(chan bool)
	elev_data  = make(chan dt.ElevDataJSON)

	//orderStateHandler channels
	ReqStateMatrix_fromP2P = make(chan dt.RequestStateMatrix)
	HallOrderArray         = make(chan [][2]bool)
	ReqStateMatrix_toP2P   = make(chan dt.RequestStateMatrix)

	// Data distributor channels
	allElevData_fromP2P = make(chan dt.AllElevDataJSON_withID)
	allElevData_toP2P   = make(chan dt.AllElevDataJSON_withID)
)

func RunAssignerAndDistributionIntegrationTest() {
	elevio.Init("localhost:15657", elevfsm.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

	// Order assigner
	go testSpecDistributor(OrderAssignerBehaviourChan,
		localIpAdressChan)
	go oassign.OrderAssigner(OrderAssignerBehaviourChan,
		localIpAdressChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	// DistributingSystem
	go elevDataDistributor.DataDistributor(allElevData_fromP2P,
		elev_data,
		HallOrderArray,
		allElevData_toP2P,
		ordersFromDistributor)
	go orderStateHandler.OrderStateHandler(ReqStateMatrix_fromP2P,
		hallEvent,
		handler_hallOrdersExecuted,
		HallOrderArray,
		ReqStateMatrix_toP2P)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(ordersLocal,
		cabEvent,
		drv_floors,
		drv_obstr,
		elev_data,
		handler_hallOrdersExecuted)

	for {
		// Kill orders that are yet to be handeled!
		<-ordersToSlaves
	}
}
