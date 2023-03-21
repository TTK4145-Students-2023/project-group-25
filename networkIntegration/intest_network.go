package intestNTW

import (
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
	masterSlaveRoleCh = make(chan dt.MasterSlaveRole)

	costFuncInputCh              = make(chan dt.CostFuncInputSlice) // Input from order distributor
	ordersFromMasterCh           = make(chan []dt.SlaveOrders)      // Input read from Master-Slave network module
	ordersToSlavesCh             = make(chan []dt.SlaveOrders)      // Input written to Master-Slave network module
	ordersElevCh                 = make(chan [dt.N_FLOORS][2]bool)
	hallOrdersExecutedCh         = make(chan []elevio.ButtonEvent)
	peerUpdate_MSCh              = make(chan peers.PeerUpdate)
	peerUpdate_DataDistributorCh = make(chan peers.PeerUpdate)
	peerUpdate_OrderHandlerCh    = make(chan peers.PeerUpdate)
	peerTxEnableCh               = make(chan bool)

	btnPressCh     = make(chan elevio.ButtonEvent)
	hallBtnPressCh = make(chan elevio.ButtonEvent)
	cabBtnPressCh  = make(chan elevio.ButtonEvent)

	drv_floors      = make(chan int)
	drv_obstr       = make(chan bool)
	localElevDataCh = make(chan dt.ElevData)

	//orderStateHandler channels
	hallOrderArrayCh = make(chan [dt.N_FLOORS][2]bool)

	// Data distributor channels
	cabRequestsToElevCh = make(chan [dt.N_FLOORS]bool)
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
	go elevio.PollButtons(btnPressCh)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnPressCh, hallBtnPressCh, cabBtnPressCh)

	// Peerlist handler
	go peers.PeerListHandler(localIP,
		peerTxEnableCh,
		peerUpdate_MSCh,
		peerUpdate_DataDistributorCh,
		peerUpdate_OrderHandlerCh)

	go masterSlaveNTW.MasterSlaveNTW(localIP,
		peerUpdate_MSCh,
		ordersToSlavesCh,
		ordersFromMasterCh,
		masterSlaveRoleCh,
	)

	// order assigner
	go oassign.OrderAssigner(localIP,
		masterSlaveRoleCh,
		costFuncInputCh,
		ordersFromMasterCh,
		ordersToSlavesCh,
		ordersElevCh)

	// DistributingSystem
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

	// FSM
	go elevfsm.FSM(ordersElevCh,
		cabBtnPressCh,
		drv_floors,
		drv_obstr,
		localElevDataCh,
		hallOrdersExecutedCh,
		cabRequestsToElevCh,
		peerTxEnableCh)

	for {
		event := <-btnPressCh
		if event.Button == elevio.BT_Cab {
			cabBtnPressCh <- event
		} else {
			hallBtnPressCh <- event
		}
	}
}
