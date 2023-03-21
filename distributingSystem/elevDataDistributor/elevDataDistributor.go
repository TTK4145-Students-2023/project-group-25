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
	HallRequestsCh <-chan [dt.N_FLOORS][2]bool,
	costFuncInputCh chan<- dt.CostFuncInputSlice,
	peerUpdateCh <-chan peers.PeerUpdate,
	cabRequestsToElevCh chan<- [dt.N_FLOORS]bool,
) {
	const dataDistributorPort = 15667
	var (
		transmittCh = make(chan dt.AllNodeInfoWithSenderIP)
		receiveCh   = make(chan dt.AllNodeInfoWithSenderIP)

		broadCastTimer = time.NewTimer(time.Hour)
		costFuncTimer  = time.NewTimer(time.Hour)
		initTimer      = time.NewTimer(2 * time.Second)

		costFuncInputSlice = dt.CostFuncInputSlice{}
		allElevData        = map[string]dt.ElevData{}
		peerList           = peers.PeerUpdate{}
	)

	broadCastTimer.Stop()
	costFuncTimer.Stop()

	go bcast.Transmitter(dataDistributorPort, transmittCh)
	go bcast.Receiver(dataDistributorPort, receiveCh)

initialization:
	for cabRequests := [dt.N_FLOORS]bool{}; ; {
		select {
		case receivedData := <-receiveCh:
			senderIP := receivedData.SenderIP
			allNodesInfo := receivedData.AllNodeInfo

			for _, nodeInfo := range allNodesInfo {
				if elevData := nodeInfo.Data; nodeInfo.IP == localIP {
					for floor, order := range elevData.CabRequests {
						cabRequests[floor] = cabRequests[floor] || order
					}
				} else if nodeInfo.IP == senderIP {
					allElevData[senderIP] = elevData
				}
			}
		case <-initTimer.C:
			cabRequestsToElevCh <- cabRequests
			allElevData[localIP] = <-localElevDataCh
			broadCastTimer.Reset(1)
			break initialization
		}
	}
	for {
		select {
		case peerList = <-peerUpdateCh:
		case allElevData[localIP] = <-localElevDataCh:
		case hallRequests := <-HallRequestsCh:
			aliveNodesElevData := []dt.NodeInfo{{IP: localIP, Data: allElevData[localIP]}}
			for _, nodeIP := range peerList.Peers {
				if nodeElevData, nodeElevDataSaved := allElevData[nodeIP]; nodeElevDataSaved {
					aliveNodesElevData = append(aliveNodesElevData, dt.NodeInfo{IP: nodeIP, Data: nodeElevData})
				}
			}
			costFuncInputSlice = dt.CostFuncInputSlice{
				HallRequests: hallRequests,
				States:       aliveNodesElevData,
			}
			costFuncTimer.Reset(1)
		case receivedData := <-receiveCh:
			senderIP := receivedData.SenderIP
			allNodesInfo := receivedData.AllNodeInfo

			if senderIP == localIP {
				break
			}
			for _, nodeInfo := range allNodesInfo {
				if nodeInfo.IP == senderIP && !reflect.DeepEqual(allElevData[senderIP], nodeInfo.Data) {
					allElevData[senderIP] = nodeInfo.Data
				}
			}
		case <-broadCastTimer.C:
			transmittCh <- dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(allElevData)}
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
