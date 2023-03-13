package intestNTW

import (
	"project/Network/Utilities/bcast"
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

	ordersFromDistributor      = make(chan dt.CostFuncInput) // Input from order distributor
	ordersFromMaster           = make(chan []byte)           // Input read from Master-Slave network module
	ordersToSlaves             = make(chan []byte)           // Input written to Master-Slave network module
	ordersLocal                = make(chan [][2]bool)
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
	ReqStateMatrix_fromP2P = make(chan dt.RequestStateMatrix)
	HallOrderArray         = make(chan [][2]bool)
	ReqStateMatrix_toP2P   = make(chan dt.RequestStateMatrix)

	// Data distributor channels
	allElevData_fromP2P = make(chan dt.AllElevDataJSON_withID)
	allElevData_toP2P   = make(chan dt.AllElevDataJSON_withID)
)

func RunNetworkWithAllTest() {
	elevio.Init("localhost:15657", elevfsm.N_FLOORS)

	//elvio
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

	// Peerlist handler
	go peers.PeerListHandler(peerUpdate_MS,
		peerUpdate_DataDistributor,
		peerUpdate_OrderHandler)

	// Receive from NTW
	go bcast.Receiver(15647, allElevData_fromP2P)
	go bcast.Receiver(15648, ReqStateMatrix_fromP2P)

	// Send to NTW
	go bcast.Transmitter(15647, allElevData_toP2P)
	go bcast.Transmitter(15648, ReqStateMatrix_toP2P)

	go masterSlaveNTW.MasterSlaveNTW(peerUpdate_MS,
		ordersToSlaves,
		ordersFromMaster,
		masterSlaveRoleChan,
	)

	// order assigner
	go oassign.OrderAssigner(masterSlaveRoleChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	// DistributingSystem
	go elevDataDistributor.DataDistributor(allElevData_fromP2P,
		elev_data,
		HallOrderArray,
		allElevData_toP2P,
		ordersFromDistributor,
		peerUpdate_DataDistributor)
	go orderStateHandler.OrderStateHandler(ReqStateMatrix_fromP2P,
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
		time.Sleep(time.Second * 10)
	}
}
