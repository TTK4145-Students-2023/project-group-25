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
	HallOrderArray <-chan [dt.N_FLOORS][2]bool,
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
		case DataFromP2P := <-allElevData_fromP2P:
			fmt.Printf("______WW recived from P2P__________\n")
			fmt.Printf("Sender ID: %v\n", DataFromP2P.ID)
			fmt.Printf("Data: %v\n", DataFromP2P.AllData)
			fmt.Printf("_________________________\n")

			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			Local_withID.AllData[recivedID] = recivedData

		case localData := <-localElevData:
			Local_DataMatrix[localIP] = localData

		case orders := <-HallOrderArray:

			data_aliveNodes := make(dt.AllElevDataJSON)
			for _, nodeID := range peerList.Peers {
				if Local_withID.AllData[nodeID] != (dt.ElevDataJSON{}) {
					data_aliveNodes[nodeID] = Local_withID.AllData[nodeID]
				}
			}
			if len(data_aliveNodes) != 0 {
				currentWorldView := dt.CostFuncInput{
					HallRequests: orders,
					States:       data_aliveNodes,
				}
				WorldView_toAssigner <- currentWorldView
			}
		}
		allElevData_toP2P <- Local_withID
		fmt.Printf("______WW sendt to P2P__________\n")
		fmt.Printf("Sender ID: %v\n", Local_withID.ID)
		fmt.Printf("Data: %v\n", Local_withID.AllData)
		fmt.Printf("_________________________\n")

	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
