package elevDataDistributor

import (
	"fmt"
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

// Statemachine for Distributor
func DataDistributor(localIP string,
	localElevDataCh <-chan dt.ElevData,
	HallOrderArrayCh <-chan [dt.N_FLOORS][2]bool,
	costFuncInputCh chan<- dt.CostFuncInputSlice,
	peerUpdateCh <-chan peers.PeerUpdate,
	cabRequestsToElevCh chan<- [dt.N_FLOORS]bool,
) {
	localNodesInfo := map[string]dt.ElevData{}
	costFuncInputSlice := dt.CostFuncInputSlice{}
	peerList := peers.PeerUpdate{}

	var (
		transmittNodesInfo = make(chan dt.AllNodeInfoWithSenderIP)
		receiveNodesInfo   = make(chan dt.AllNodeInfoWithSenderIP)
	)
	go bcast.Receiver(15667, receiveNodesInfo)
	go bcast.Transmitter(15667, transmittNodesInfo)

	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()
	broadCastTimer := time.NewTimer(1)
	broadCastTimer.Stop()

	initTimer := time.NewTimer(1 * time.Second)

	//Send cabcalls to fsm if cabcall allready exist on network
	initCabCalls := [dt.N_FLOORS]bool{}
	for initTimeOut := false; !initTimeOut; {
		select {
		case NTWData := <-receiveNodesInfo:
			fmt.Printf("NodesInfoFromNTW:\n %+v\n\n", NTWData)
			senderIP := NTWData.SenderIP
			senderNodesInfo := dt.NodeInfoSliceToMap(NTWData.AllNodeInfo)

			localNodesInfo[senderIP] = senderNodesInfo[senderIP]

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
		case newNodesInfo := <-receiveNodesInfo:
			senderData := dt.NodeInfoSliceToMap(newNodesInfo.AllNodeInfo)
			senderIP := newNodesInfo.SenderIP
			//fmt.Printf("nodesInfo fro NTW: %+v\n", senderData)
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localNodesInfo[senderIP]) {
				localNodesInfo[senderIP] = senderData[senderIP]
			}

		case elevData := <-localElevDataCh:
			if _, valInMap := localNodesInfo[localIP]; !valInMap {
				broadCastTimer.Reset(1)
			}
			localNodesInfo[localIP] = elevData

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
			}
		case <-broadCastTimer.C:
			//fmt.Printf("nodesInfo to NTW: %+v\n", dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(localNodesInfo)})
			transmittNodesInfo <- dt.AllNodeInfoWithSenderIP{SenderIP: localIP, AllNodeInfo: dt.NodeInfoMapToSlice(localNodesInfo)}
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)

		case <-worldViewTimer.C:
			select {
			case costFuncInputCh <- costFuncInputSlice:

			default:
				worldViewTimer.Reset(1)
			}
		}
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
