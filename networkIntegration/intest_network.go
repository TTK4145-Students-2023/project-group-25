package intestNTW

import (
	P2P "project/Network/P2Pntw"
	"project/Network/Utilities/localip"
	"project/Network/Utilities/peers"
	masterSlaveNTW "project/Network/masterSlaveNTW"
	btnassign "project/buttonAssigner"
	dt "project/commonDataTypes"
	elevDataDistributor "project/distributingSystem/elevDataDistributor"
	orderStateHandler "project/distributingSystem/orderStateHandler"
	elevio "project/localElevator/elev_driver"
	elevfsm "project/localElevator/elev_fsm"
	oassign "project/orderAssigner"
	"time"
)

var (
	masterSlaveRoleChan = make(chan dt.MasterSlaveRole)

	ordersFromDistributor      = make(chan dt.CostFuncInput)                // Input from order distributor
	ordersFromMaster           = make(chan map[string][dt.N_FLOORS][2]bool) // Input read from Master-Slave network module
	ordersToSlaves             = make(chan map[string][dt.N_FLOORS][2]bool) // Input written to Master-Slave network module
	ordersLocal                = make(chan [dt.N_FLOORS][2]bool)
	handler_hallOrdersExecuted = make(chan []elevio.ButtonEvent)
	peerUpdate_MS              = make(chan peers.PeerUpdate)
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

func RunNetworkWithAllTest() {
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)

	for floor := 0; floor < dt.N_FLOORS; floor++ {
		for button := 0; button < dt.N_BUTTONS; button++ {
			elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
		}
	}

	//elvio
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

	// Peerlist handler
	go peers.PeerListHandler(localIP,
		peerUpdate_MS,
		peerUpdate_DataDistributor,
		peerUpdate_OrderHandler)

	go masterSlaveNTW.MasterSlaveNTW(localIP,
		peerUpdate_MS,
		ordersToSlaves,
		ordersFromMaster,
		masterSlaveRoleChan,
	)
	go P2P.P2Pntw(localIP,
		allElevData_toP2P,
		ReqStateMatrix_toP2P,
		allElevData_fromP2P,
		ReqStateMatrix_fromP2P)
	// order assigner
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

	// FSM
	go elevfsm.FSM(ordersLocal,
		cabEvent,
		drv_floors,
		drv_obstr,
		elev_data,
		handler_hallOrdersExecuted)

	for {
		event := <-btnEvent
		if event.Button == elevio.BT_Cab {
			cabEvent <- event
		} else {
			hallEvent <- event
		}
	}
}
