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
	cabRequestsToElevCh chan<- [dt.N_FLOORS]bool,
) {
	localNodesInfo := map[string]dt.ElevData{}
	costFuncInputSlice := dt.CostFuncInputSlice{}
	peerList := peers.PeerUpdate{}

	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()
	allElevDataTimer := time.NewTimer(1)
	allElevDataTimer.Stop()

	initTimer := time.NewTimer(1 * time.Second)

	//Send cabcalls to fsm if cabcall allready exist on network
	initCabCalls := [dt.N_FLOORS]bool{}
	for initTimeOut := false; !initTimeOut; {
		select {
		case NTWData := <-nodesInfoFromNTWCh:
			fmt.Printf("NodesInfoFromNTW:\n %+v\n\n", NTWData)
			senderIP := NTWData.SenderIP
			senderNodesInfo := dt.NodeInfoSliceToMap(NTWData.AllNodeInfo)

			localNodesInfo[senderIP] = senderNodesInfo[senderIP]
			allElevDataTimer.Reset(1)

			if _, valInMap := senderNodesInfo[localIP]; valInMap {
				for floor, order := range senderNodesInfo[localIP].CabRequests {
					if order {
						initCabCalls[floor] = order
					}
				}
				//send to fsm
				fmt.Printf("\ncabcalls from NTW: %+v\n", initCabCalls)
			}
		case <-initTimer.C:
			initTimeOut = true
			cabRequestsToElevCh <- initCabCalls
		}
	}

	fmt.Println("enter normal DD")
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
			if localElevData, valInMap := localNodesInfo[localIP]; valInMap {
				aliveNodesInfo := map[string]dt.ElevData{localIP: localElevData}
				aliveNodesInfo[localIP] = localNodesInfo[localIP]
				for _, nodeIP := range peerList.Peers {
					if localNodesInfo[nodeIP] != (dt.ElevData{}) {
						aliveNodesInfo[nodeIP] = localNodesInfo[nodeIP]
					}
				}
				costFuncInputSlice = dt.CostFuncInputSlice{
					HallRequests: orders,
					States:       dt.NodeInfoMapToSlice(aliveNodesInfo),
				}
				worldViewTimer.Reset(1)
				allElevDataTimer.Reset(1)
			}

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
