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
	ordersElevCh                 = make(chan [dt.N_FLOORS][2]bool)
	hallOrdersExecutedCh         = make(chan []elevio.ButtonEvent)
	peerUpdate_MSCh              = make(chan peers.PeerUpdate)
	peerUpdate_DataDistributorCh = make(chan peers.PeerUpdate)
	peerUpdate_OrderHandlerCh    = make(chan peers.PeerUpdate)

	btnPressCh     = make(chan elevio.ButtonEvent)
	hallBtnPressCh = make(chan elevio.ButtonEvent)
	cabBtnPressCh  = make(chan elevio.ButtonEvent)

	drv_floors      = make(chan int)
	drv_obstr       = make(chan bool)
	localElevDataCh = make(chan dt.ElevData)

	hallOrderArrayCh = make(chan [dt.N_FLOORS][2]bool)

	cabRequestsToElevCh = make(chan [dt.N_FLOORS]bool)
)

func main() {
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)
	elevio.ClearAllLights()

	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnPressCh)
	go elevio.PollObstructionSwitch(drv_obstr)

	
	go peers.PeerListHandler(localIP,
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
		ordersElevCh)
		
	go elevDataDistributor.DataDistributor(localIP,
		localElevDataCh,
		hallOrderArrayCh,
		costFuncInputCh,
		peerUpdate_DataDistributorCh,
		cabRequestsToElevCh)
	
	go orderStateHandler.OrderStateHandler(localIP,
		hallBtnPressCh,
		hallOrdersExecutedCh,
		hallOrderArrayCh,
		peerUpdate_OrderHandlerCh)
	
	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(ordersElevCh,
		cabBtnPressCh,
		drv_floors,
		drv_obstr,
		localElevDataCh,
		hallOrdersExecutedCh,
		cabRequestsToElevCh)
		
	btnassign.ButtonHandler(btnPressCh, hallBtnPressCh, cabBtnPressCh)
}
