package elevDataDistributor

import (
	"fmt"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
)

// Statemachine for Distributor
func DataDistributor(localIP string,
	allElevData_fromP2P <-chan dt.AllElevDataJSON_withID,
	localElevData <-chan dt.ElevDataJSON,
	HallOrderArray <-chan [][2]bool,
	allElevData_toP2P chan<- dt.AllElevDataJSON_withID,
	WorldView_toAssigner chan<- dt.CostFuncInput,
	peerUpdateChan <-chan peers.PeerUpdate,
) {
	//init local Data Matrix with local ID
	Local_DataMatrix := make(dt.AllElevDataJSON)

	Local_withID := dt.AllElevDataJSON_withID{
		ID:      localIP,
		AllData: Local_DataMatrix,
	}

	// List of node IDs we are connected to
	peerList := peers.PeerUpdate{}

	for {
		select {
		case peerList = <-peerUpdateChan:
			// Initilize new nodes
			for _, nodeID := range peerList.Peers {
				if _, valInMap := Local_DataMatrix[nodeID]; !valInMap {
					Local_DataMatrix[nodeID] = dt.ElevDataJSON{}
				}
			}
		case DataFromP2P := <-allElevData_fromP2P:

			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			Local_withID.AllData[recivedID] = recivedData

		case localData := <-localElevData:
			Local_DataMatrix[localIP] = localData

		case orders := <-HallOrderArray:
			data_aliveNodes := make(dt.AllElevDataJSON)
			for _, nodeID := range peerList.Peers {
				data_aliveNodes[nodeID] = Local_withID.AllData[nodeID]
			}

			currentWorldView := dt.CostFuncInput{
				HallRequests: orders,
				States:       data_aliveNodes,
			}

			WorldView_toAssigner <- currentWorldView

		}
		allElevData_toP2P <- Local_withID
		fmt.Printf("WW sendt to P2P:\n %+v\n", Local_withID)

	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
