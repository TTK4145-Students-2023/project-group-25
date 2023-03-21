package main

// Import statements in alphabetic order
import (
	P2P "project/Network/P2Pntw"
	localip "project/Network/Utilities/localip"
	peers "project/Network/Utilities/peers"
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

// Constant declerations (should we have a "constant" package where we defines all constants?)
// const (
// 	N_FLOORS  = 4
// 	N_BUTTONS = 3
// )
// const (
// 	MS_Master dt.MasterSlaveRole = "master"
// 	MS_Slave  dt.MasterSlaveRole = "slave"
// )

// Variable declerations

// Channels
var (
	masterSlaveRoleChan = make(chan dt.MasterSlaveRole)

	ordersFromDistributor      = make(chan dt.CostFuncInput)
	ordersFromMaster           = make(chan map[string][dt.N_FLOORS][2]bool)
	ordersToSlaves             = make(chan map[string][dt.N_FLOORS][2]bool)
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
	allElevData_toP2P   = make(chan dt.AllElevDataJSON)
)

func main() {
	// Initialization phase
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)
	elevio.ResetLights()

	// Main program
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

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

	go oassign.OrderAssigner(localIP,
		masterSlaveRoleChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

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

	time.Sleep(time.Millisecond * 40) // why this?

	go elevfsm.FSM(ordersLocal,
		cabEvent,
		drv_floors,
		drv_obstr,
		elev_data,
		handler_hallOrdersExecuted)

	// Should make another solution to keep the program running
	for {
		event := <-btnEvent
		if event.Button == elevio.BT_Cab {
			cabEvent <- event
		} else {
			hallEvent <- event
		}
	}
}
