package elevDataDistributor

import (
	"fmt"
	"project/Network/Utilities/localip"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
)

// Statemachine for Distributor
func DataDistributor(
	allElevData_fromP2P <-chan dt.AllElevDataJSON_withID,
	localElevData <-chan dt.ElevDataJSON,
	HallOrderArray <-chan [][2]bool,
	allElevData_toP2P chan<- dt.AllElevDataJSON_withID,
	WorldView_toAssigner chan<- dt.CostFuncInput,
	peerUpdateChan <-chan peers.PeerUpdate,
) {
	localIpAdress, _ := localip.LocalIP()
	//init local Data Matrix with local ID
	Local_DataMatrix := make(dt.AllElevDataJSON)

	Local_withID := dt.AllElevDataJSON_withID{
		ID:      localIpAdress,
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
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n AllElevdata from P2P recieved: \n%+v\n", DataFromP2P)

			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			Local_withID.AllData[recivedID] = recivedData

			allElevData_toP2P <- Local_withID

		case localData := <-localElevData:
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n Local Elevdata recieved: \n%+v\n", localData)
			Local_DataMatrix[localIpAdress] = localData

		case orders := <-HallOrderArray:
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n HallOrderArray recieved: \n%+v\n", orders)
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
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
