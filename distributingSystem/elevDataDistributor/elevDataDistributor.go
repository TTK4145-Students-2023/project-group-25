package elevDataDistributor

import (
	"fmt"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"time"
)

// Statemachine for Distributor
func DataDistributor(localIP string,
	nodesInfoFromNTWCh <-chan dt.AllNodeInfoWithSenderIP,
	localElevDataCh <-chan dt.ElevData,
	HallOrderArrayCh <-chan [dt.N_FLOORS][2]bool,
	nodeInfoToNTWCh chan<- []dt.NodeInfo,
	costFuncInputCh chan<- dt.CostFuncInputSlice,
	peerUpdateCh <-chan peers.PeerUpdate,
) {
	localNodesInfo := map[string]dt.ElevData{}
	costFuncInputSlice := dt.CostFuncInputSlice{}
	peerList := peers.PeerUpdate{}

	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()
	allElevDataTimer := time.NewTimer(1)
	allElevDataTimer.Stop()

	//Send cabcalls to fsm if cabcall allready exist on network
	select {
	case NTWData := <-nodesInfoFromNTWCh:
		senderNodesInfo := dt.NodeInfoSliceToMap(NTWData.AllNodeInfo)
		if _, valInMap := senderNodesInfo[localIP]; valInMap {
			cabcalls := senderNodesInfo[localIP].CabRequests
			//send to fsm
			fmt.Printf("\ncabcalls from NTW: %+v\n", cabcalls)
		}
	}

	for {
		select {
		case peerList = <-peerUpdateCh:
		case NTWData := <-nodesInfoFromNTWCh:

			senderIP := NTWData.SenderIP
			senderNodesInfo := dt.NodeInfoSliceToMap(NTWData.AllNodeInfo)

			localNodesInfo[senderIP] = senderNodesInfo[senderIP]
			allElevDataTimer.Reset(1)

		case elevData := <-localElevDataCh:
			localNodesInfo[localIP] = elevData
			allElevDataTimer.Reset(1)

		case orders := <-HallOrderArrayCh:
			aliveNodesInfo := map[string]dt.ElevData{}
			for _, nodeIP := range peerList.Peers {
				if localNodesInfo[nodeIP] != (dt.ElevData{}) {
					aliveNodesInfo[nodeIP] = localNodesInfo[nodeIP]
				}
			}

			if len(aliveNodesInfo) > 0 {
				costFuncInputSlice = dt.CostFuncInputSlice{
					HallRequests: orders,
					States:       dt.NodeInfoMapToSlice(aliveNodesInfo),
				}
				worldViewTimer.Reset(1)
			}
			allElevDataTimer.Reset(1)
		case <-worldViewTimer.C:
			select {
			case costFuncInputCh <- costFuncInputSlice:

			default:
				worldViewTimer.Reset(1)
			}
		case <-allElevDataTimer.C:
			select {
			case nodeInfoToNTWCh <- dt.NodeInfoMapToSlice(localNodesInfo):
			default:
				allElevDataTimer.Reset(1)
			}
		}

	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
