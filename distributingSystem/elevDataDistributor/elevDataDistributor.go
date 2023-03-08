package elevDataDistributor

import "fmt"

var localID string = "ID1"

// Datatypes
type ElevData struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type AllElevData map[string]ElevData

type AllElevData_withID struct {
	ID      string      `json:"id"`
	AllData AllElevData `json:"allData"`
}

type WorldView struct {
	HallRequests [][2]bool   `json:"hallRequests"`
	AllData      AllElevData `json:"elevStates"`
}

// input channels
var (
	allElevData_fromP2P = make(chan AllElevData_withID)
	localElevData       = make(chan ElevData)
	HallOrderArray      = make(chan [][2]bool)
)

// output channels
var (
	allElevData_toP2P   = make(chan AllElevData)
	WorlView_toAssigner = make(chan WorldView)
)

// Statemachine for Distributor
func dataDistributor(
	allElevData_fromP2P <-chan AllElevData_withID,
	localElevData <-chan ElevData,
	HallOrderArray <-chan [][2]bool,
	allElevData_toP2P chan<- AllElevData_withID,
	WorldView_toAssigner chan<- WorldView,

) {

	//init local Data Matrix with local ID
	Local_DataMatrix := make(AllElevData)
	Local_DataMatrix["ID1"] = ElevData{}
	Local_DataMatrix["ID2"] = ElevData{}
	Local_DataMatrix["ID3"] = ElevData{}

	Local_withID := AllElevData_withID{
		ID:      "ID1",
		AllData: Local_DataMatrix,
	}

	// List of node IDs we are connected to
	nodeIDs := []string{"ID1", "ID2", "ID3"}

	for {
		select {
		case DataFromP2P := <-allElevData_fromP2P:
			recivedID := DataFromP2P.ID
			recivedData := DataFromP2P.AllData[recivedID]

			//fmt.Printf(" recieved ID: \n%v\n recieved DATA: \n%v\n", recivedID, recivedData)

			Local_withID.AllData[recivedID] = recivedData

			fmt.Printf("Local data after update: \n%v\n", Local_withID)

			allElevData_toP2P <- Local_withID

		case localData := <-localElevData:

			Local_DataMatrix[localID] = localData

		case orders := <-HallOrderArray:

			data_aliveNodes := make(AllElevData)
			for _, nodeID := range nodeIDs {
				data_aliveNodes[nodeID] = Local_withID.AllData[nodeID]
			}

			currentWorldView := WorldView{
				HallRequests: orders,
				AllData:      data_aliveNodes,
			}

			WorldView_toAssigner <- currentWorldView

		}
	}
}
