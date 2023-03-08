package elevDataDistributor

import (
	"testing"
)

//Tests for elevDataDist

func TestDataDistributor(t *testing.T) {

	//Test 1:
	t.Run("Check output for external and local state input", func(t *testing.T) {
		//input and output channels
		allElevData_fromP2P := make(chan AllElevData_withID)
		localElevData := make(chan ElevData)
		HallOrderArray := make(chan [][2]bool)
		allElevData_toP2P := make(chan AllElevData_withID)
		WorlView_toAssigner := make(chan WorldView)

		//start the distributor as a goroutine
		go dataDistributor(
			allElevData_fromP2P,
			localElevData,
			HallOrderArray,
			allElevData_toP2P,
			WorlView_toAssigner,
		)

		//mocking inputs
		input_localData := ElevData{
			Behavior:    "Moving",
			Floor:       3,
			Direction:   "up",
			CabRequests: []bool{false, false, false, false},
		}

		//input data from P2P
		DataMatrix := make(AllElevData)
		DataMatrix["ID1"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID2"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID3"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		input_datamatrix_withID := AllElevData_withID{
			ID:      "ID2",
			AllData: DataMatrix,
		}

		//send inputs on channel
		localElevData <- input_localData
		allElevData_fromP2P <- input_datamatrix_withID

		//read output on channel
		output_Data := <-allElevData_toP2P

		//print output
		t.Logf("\ninput datamtrix: \n%v\n input localState \n%v\n output datamatrix \n%v\n",
			input_datamatrix_withID, input_localData, output_Data)

	})

	t.Run("Check output for hallReq and external state input", func(t *testing.T) {

		//input and output channels
		allElevData_fromP2P := make(chan AllElevData_withID)
		localElevData := make(chan ElevData)
		HallOrderArray := make(chan [][2]bool)
		allElevData_toP2P := make(chan AllElevData_withID)
		WorlView_toAssigner := make(chan WorldView)

		//start the distributor as a goroutine
		go dataDistributor(
			allElevData_fromP2P,
			localElevData,
			HallOrderArray,
			allElevData_toP2P,
			WorlView_toAssigner,
		)

		//mocking inputs
		input_HallOrders := [][2]bool{{true, false}, {true, false}, {true, false}, {true, false}}

		//input data from P2P
		DataMatrix := make(AllElevData)
		DataMatrix["ID1"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID2"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID3"] = ElevData{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		input_datamatrix_withID := AllElevData_withID{
			ID:      "ID2",
			AllData: DataMatrix,
		}

		//send inputs on channel
		allElevData_fromP2P <- input_datamatrix_withID
		//read output on channel
		output_Datamatrix := <-allElevData_toP2P

		HallOrderArray <- input_HallOrders

		output_WorldView := <-WorlView_toAssigner

		//print output
		t.Errorf("\ninput datamtrix: \n%v\n input orderArray \n%v\n output datamatrix \n%v\n output WorlView: \n%v\n",
			input_datamatrix_withID, input_HallOrders, output_Datamatrix, output_WorldView)

	})

}
