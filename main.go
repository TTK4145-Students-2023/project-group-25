package main

import (
	dt "project/dataTypes"
	elevDataDistributor "project/distributingSystem/elevDataDistributor"
	orderStateHandler "project/distributingSystem/orderStateHandler"
	elevfsm "project/localElevator/FSM"
	btnEventSplitter "project/localElevator/btnEventSplitter"
	elevio "project/localElevator/driver"
	localip "project/network/localip"
	peers "project/network/peers"
	oassign "project/orderAssigner"
	"time"
)

var (
	floorCh           = make(chan int)
	obstrCh           = make(chan bool)
	elevDataCh        = make(chan dt.ElevData)
	buttonEventCh     = make(chan elevio.ButtonEvent)
	cabButtonEventCh  = make(chan elevio.ButtonEvent)
	hallButtonEventCh = make(chan elevio.ButtonEvent)

	costFuncInputCh     = make(chan dt.CostFuncInputSlice)
	executedHallOrderCh = make(chan elevio.ButtonEvent)
	assignedOrdersCh    = make(chan [dt.N_FLOORS][2]bool)
	orderStatesToBoolCh = make(chan [dt.N_FLOORS][2]bool)
	initCabRequestsCh   = make(chan [dt.N_FLOORS]bool)

	peerUpdate_OrderAssCh        = make(chan peers.PeerUpdate)
	peerUpdate_DataDistributorCh = make(chan peers.PeerUpdate)
	peerUpdate_OrderHandlerCh    = make(chan peers.PeerUpdate)
	peerTxEnableCh               = make(chan bool)
)

func main() {
	localIP, _ := localip.LocalIP()
	elevio.Init("localhost:15657", dt.N_FLOORS)
	elevio.ClearAllLights()

	go elevio.PollFloorSensor(floorCh)
	go elevio.PollButtons(buttonEventCh)
	go elevio.PollObstructionSwitch(obstrCh)
	go btnEventSplitter.BtnEventSplitter(buttonEventCh, hallButtonEventCh, cabButtonEventCh)

	go peers.PeerListHandler(localIP,
		peerTxEnableCh,
		peerUpdate_OrderAssCh,
		peerUpdate_DataDistributorCh,
		peerUpdate_OrderHandlerCh)

	go oassign.OrderAssigner(localIP,
		peerUpdate_OrderAssCh,
		costFuncInputCh,
		assignedOrdersCh)

	go elevDataDistributor.DataDistributor(localIP,
		elevDataCh,
		orderStatesToBoolCh,
		costFuncInputCh,
		peerUpdate_DataDistributorCh,
		initCabRequestsCh)

	go orderStateHandler.OrderStateHandler(localIP,
		hallButtonEventCh,
		executedHallOrderCh,
		orderStatesToBoolCh,
		peerUpdate_OrderHandlerCh)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(assignedOrdersCh,
		cabButtonEventCh,
		floorCh,
		obstrCh,
		elevDataCh,
		executedHallOrderCh,
		initCabRequestsCh,
		peerTxEnableCh)

	select {}
}
