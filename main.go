package main

// Import statements in alphabetic order
import (
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

var (
	masterSlaveRoleCh = make(chan dt.MasterSlaveRole)

	costFuncInputCh              = make(chan dt.CostFuncInputSlice)
	ordersFromMasterCh           = make(chan []dt.SlaveOrders)
	ordersToSlavesCh             = make(chan []dt.SlaveOrders)
	hallRequestsCh               = make(chan [dt.N_FLOORS][2]bool)
	executedHallOrderCh          = make(chan elevio.ButtonEvent)
	peerUpdate_MSCh              = make(chan peers.PeerUpdate)
	peerUpdate_DataDistributorCh = make(chan peers.PeerUpdate)
	peerUpdate_OrderHandlerCh    = make(chan peers.PeerUpdate)
	peerTxEnableCh               = make(chan bool)

	buttonEventCh     = make(chan elevio.ButtonEvent)
	cabButtonEventCh  = make(chan elevio.ButtonEvent)
	hallButtonEventCh = make(chan elevio.ButtonEvent)

	floorCh    = make(chan int)
	obstrCh    = make(chan bool)
	elevDataCh = make(chan dt.ElevData)

	hallOrderArrayCh = make(chan [dt.N_FLOORS][2]bool)

	initCabRequestsCh = make(chan [dt.N_FLOORS]bool)
)

func main() {
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)
	elevio.ClearAllLights()

	go elevio.PollFloorSensor(floorCh)
	go elevio.PollButtons(buttonEventCh)
	go elevio.PollObstructionSwitch(obstrCh)

	go peers.PeerListHandler(localIP,
		peerTxEnableCh,
		peerUpdate_MSCh,
		peerUpdate_DataDistributorCh,
		peerUpdate_OrderHandlerCh)

	go masterSlaveNTW.MasterSlaveNTW(localIP,
		peerUpdate_MSCh,
		ordersToSlavesCh,
		ordersFromMasterCh,
		masterSlaveRoleCh)

	go oassign.OrderAssigner(localIP,
		masterSlaveRoleCh,
		costFuncInputCh,
		ordersFromMasterCh,
		ordersToSlavesCh,
		hallRequestsCh)

	go elevDataDistributor.DataDistributor(localIP,
		elevDataCh,
		hallOrderArrayCh,
		costFuncInputCh,
		peerUpdate_DataDistributorCh,
		initCabRequestsCh)

	go orderStateHandler.OrderStateHandler(localIP,
		hallButtonEventCh,
		executedHallOrderCh,
		hallOrderArrayCh,
		peerUpdate_OrderHandlerCh)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(hallRequestsCh,
		cabButtonEventCh,
		floorCh,
		obstrCh,
		elevDataCh,
		executedHallOrderCh,
		initCabRequestsCh,
		peerTxEnableCh)

	go btnassign.ButtonHandler(buttonEventCh, hallButtonEventCh, cabButtonEventCh)

	select {}
}
