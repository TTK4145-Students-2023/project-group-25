package intest

import (
	"project/Network/Utilities/localip"
	"project/Network/Utilities/peers"
	btnassign "project/buttonAssigner"
	dt "project/commonDataTypes"
	elevDataDistributor "project/distributingSystem/elevDataDistributor"
	orderStateHandler "project/distributingSystem/orderStateHandler"
	elevio "project/localElevator/elev_driver"
	elevfsm "project/localElevator/elev_fsm"
	oassign "project/orderAssigner"
	"time"
)

func testSpecDistributor(masterSlaveRoleChan chan dt.MasterSlaveRole) {
	masterSlaveRole := dt.MS_Master

	for {
		masterSlaveRoleChan <- masterSlaveRole
	}
}

var (
	masterSlaveRoleChan = make(chan dt.MasterSlaveRole)

	ordersFromDistributor      = make(chan dt.CostFuncInput)                // Input from order distributor
	ordersFromMaster           = make(chan map[string][dt.N_FLOORS][2]bool) // Input read from Master-Slave network module
	ordersToSlaves             = make(chan map[string][dt.N_FLOORS][2]bool) // Input written to Master-Slave network module
	ordersLocal                = make(chan [dt.N_FLOORS][2]bool)
	handler_hallOrdersExecuted = make(chan []elevio.ButtonEvent)
	peerUpdate_DataDistributor = make(chan peers.PeerUpdate)
	peerUpdate_OrderHandler    = make(chan peers.PeerUpdate)

	btnEvent  = make(chan elevio.ButtonEvent)
	hallEvent = make(chan elevio.ButtonEvent)
	cabEvent  = make(chan elevio.ButtonEvent)

	drv_floors = make(chan int)
	drv_obstr  = make(chan bool)
	elev_data  = make(chan dt.ElevDataJSON)

	//orderStateHandler channels
	ReqStateMatrix_fromP2P = make(chan dt.RequestStateMatrix_with_ID)
	HallOrderArray         = make(chan [dt.N_FLOORS][2]bool)
	ReqStateMatrix_toP2P   = make(chan dt.RequestStateMatrix)

	// Data distributor channels
	allElevData_fromP2P = make(chan dt.AllElevDataJSON_withID)
	allElevData_toP2P   = make(chan dt.AllElevDataJSON_withID)
)

func RunAssignerAndDistributionIntegrationTest() {
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

	// Order assigner
	go testSpecDistributor(masterSlaveRoleChan)
	go oassign.OrderAssigner(localIP,
		masterSlaveRoleChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	// DistributingSystem
	go elevDataDistributor.DataDistributor(localIP,
		allElevData_fromP2P,
		elev_data,
		HallOrderArray,
		allElevData_toP2P,
		ordersFromDistributor,
		peerUpdate_DataDistributor)
	go orderStateHandler.OrderStateHandler(localIP,
		ReqStateMatrix_fromP2P,
		hallEvent,
		handler_hallOrdersExecuted,
		HallOrderArray,
		ReqStateMatrix_toP2P,
		peerUpdate_OrderHandler)

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
