package elevDataDistributor

import (
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"time"
)

// Statemachine for Distributor
func DataDistributor(localIP string,
	allElevData_fromP2P <-chan dt.AllElevDataJSON_withID,
	localElevData <-chan dt.ElevDataJSON,
	HallOrderArray <-chan [dt.N_FLOORS][2]bool,
	allElevData_toP2P chan<- dt.AllElevDataJSON_withID,
	WorldView_toAssigner chan<- dt.CostFuncInput,
	peerUpdateChan <-chan peers.PeerUpdate,
) {
	//init local Data Matrix with local ID
	Local_DataMatrix := make(dt.AllElevDataJSON)
	currentWorldView := dt.CostFuncInput{}
	peerList := peers.PeerUpdate{}

	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()
	allElevDataTimer := time.NewTimer(1)
	allElevDataTimer.Stop()

	for {
		select {
		case peerList = <-peerUpdateChan:
		case DataFromP2P := <-allElevData_fromP2P:

			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			Local_DataMatrix[recivedID] = recivedData

		case localData := <-localElevData:
			Local_DataMatrix[localIP] = localData

		case orders := <-HallOrderArray:
			data_aliveNodes := make(dt.AllElevDataJSON)
			for _, nodeID := range peerList.Peers {
				if Local_DataMatrix[nodeID] != (dt.ElevDataJSON{}) {
					data_aliveNodes[nodeID] = Local_DataMatrix[nodeID]
				}
			}

			if len(data_aliveNodes) != 0 {
				currentWorldView = dt.CostFuncInput{
					HallRequests: orders,
					States:       data_aliveNodes,
				}
				worldViewTimer.Reset(1)
			}
		case <-worldViewTimer.C:
			select {
			case WorldView_toAssigner <- currentWorldView:
			default:
				worldViewTimer.Reset(1)
			}
		case <-allElevDataTimer.C:
			select {
			case allElevData_toP2P <- dt.AllElevDataJSON_withID{ID: localIP, AllData: Local_DataMatrix}:
			default:
			}
		}
		allElevDataTimer.Reset(1)
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
