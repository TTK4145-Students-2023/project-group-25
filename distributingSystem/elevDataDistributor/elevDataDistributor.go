package elevDataDistributor

import (
	"fmt"
	dt "project/commonDataTypes"
)

var localID string = "127.0.0.1"

// Statemachine for Distributor
func DataDistributor(
	allElevData_fromP2P <-chan dt.AllElevDataJSON_withID,
	localElevData <-chan dt.ElevDataJSON,
	HallOrderArray <-chan [][2]bool,
	allElevData_toP2P chan<- dt.AllElevDataJSON_withID,
	WorldView_toAssigner chan<- dt.CostFuncInput,
) {
	//init local Data Matrix with local ID
	Local_DataMatrix := make(dt.AllElevDataJSON)
	Local_DataMatrix[localID] = dt.ElevDataJSON{}
	//Local_DataMatrix["ID2"] = dt.ElevDataJSON{}
	//Local_DataMatrix["ID3"] = dt.ElevDataJSON{}

	Local_withID := dt.AllElevDataJSON_withID{
		ID:      localID,
		AllData: Local_DataMatrix,
	}

	// List of node IDs we are connected to
	nodeIDs := []string{localID} //, "ID2", "ID3"}

	for {
		select {
		case DataFromP2P := <-allElevData_fromP2P:
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n AllElevdata from P2P recieved: \n%+v\n", DataFromP2P)

			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			Local_withID.AllData[recivedID] = recivedData

			allElevData_toP2P <- Local_withID

		case localData := <-localElevData:
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n Local Elevdata recieved: \n%+v\n", localData)
			Local_DataMatrix[localID] = localData

		case orders := <-HallOrderArray:
			fmt.Printf("\n___DATA_DISTRIBUTOR___: \n HallOrderArray recieved: \n%+v\n", orders)
			data_aliveNodes := make(dt.AllElevDataJSON)
			for _, nodeID := range nodeIDs {
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
