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
	peerList := peers.PeerUpdate{}

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
				currentWorldView := dt.CostFuncInput{
					HallRequests: orders,
					States:       data_aliveNodes,
				}
				fmt.Printf("DATADIST, deadlock 1! ")
				WorldView_toAssigner <- currentWorldView
				fmt.Printf("... kidding, no DATADIST deadlock 1...\n ")
			}
		}
		fmt.Printf("DATADIST, deadlock 2! ")
		allElevData_toP2P <- dt.AllElevDataJSON_withID{ID: localIP, AllData: Local_DataMatrix}
		fmt.Printf("... kidding, no DATADIST deadlock 2...\n ")
	}
}

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
