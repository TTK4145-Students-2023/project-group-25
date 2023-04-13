package elevDataDistributor

import (
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

// Statemachine for Distributor
func DataDistributor(localIP string,
	localElevDataCh <-chan dt.ElevData,
	confirmedOrdersCh <-chan [dt.N_FLOORS][2]bool,
	costFuncInputCh chan<- dt.CostFuncInputSlice,
	peerUpdateCh <-chan peers.PeerUpdate,
	initCabRequestsCh chan<- [dt.N_FLOORS]bool,
) {

	var (
		receiveCh   = make(chan dt.AllNodeInfoWithSenderIP)
		transmittCh = make(chan dt.AllNodeInfoWithSenderIP)

		initTimer      = time.NewTimer(time.Hour)
		costFuncTimer  = time.NewTimer(time.Hour)
		broadCastTimer = time.NewTimer(time.Hour)

		peerList           = peers.PeerUpdate{}
		allElevData        = map[string]dt.ElevData{}
		costFuncInputSlice = dt.CostFuncInputSlice{}
	)

	initTimer.Stop()
	costFuncTimer.Stop()
	broadCastTimer.Stop()

	go bcast.Receiver(dt.DATA_DISTRIBUTOR_PORT, receiveCh)
	go bcast.Transmitter(dt.DATA_DISTRIBUTOR_PORT, transmittCh)

	initTimer.Reset(time.Second * 3)
	initCabRequests := [dt.N_FLOORS]bool{}
initialization:
	for {
		select {
		case receivedData := <-receiveCh:
			senderIP := receivedData.SenderIP
			allNodesInfo := receivedData.AllNodeInfo

			for nodeIP, receivedElevData := range allNodesInfo {
				if nodeIP == localIP {
					for floor, receivedOrder := range receivedElevData.CabRequests {
						initCabRequests[floor] = initCabRequests[floor] || receivedOrder
					}
				} else if nodeIP == senderIP {
					allElevData[senderIP] = receivedElevData
				}
			}
		case <-initTimer.C:
			initCabRequestsCh <- initCabRequests
			allElevData[localIP] = <-localElevDataCh
			broadCastTimer.Reset(1)
			break initialization
		}
	}
	for {
		select {
		case peerList = <-peerUpdateCh:
		case allElevData[localIP] = <-localElevDataCh:
		case hallRequests := <-confirmedOrdersCh:
			aliveNodesData := []dt.NodeInfo{{IP: localIP, Data: allElevData[localIP]}}
			for _, nodeIP := range peerList.Peers {
				nodeData, nodeDataExists := allElevData[nodeIP]
				if nodeDataExists {
					aliveNodesData = append(aliveNodesData, dt.NodeInfo{IP: nodeIP, Data: nodeData})
				}
			}
			costFuncInputSlice = dt.CostFuncInputSlice{
				HallRequests: hallRequests,
				States:       aliveNodesData,
			}
			costFuncTimer.Reset(1)
		case receivedData := <-receiveCh:
			senderIP := receivedData.SenderIP
			receivedNodesInfo := receivedData.AllNodeInfo

			if senderIP == localIP {
				break
			}
			for NodeIP, receivedElevData := range receivedNodesInfo {
				if NodeIP == senderIP && !reflect.DeepEqual(allElevData[senderIP], receivedElevData) {
					allElevData[senderIP] = receivedElevData
				}
			}
		case <-broadCastTimer.C:
			transmittCh <- dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: allElevData}
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)

		case <-costFuncTimer.C:
			select {
			case costFuncInputCh <- costFuncInputSlice:
			default:
				costFuncTimer.Reset(1)
			}
		}
	}
}
