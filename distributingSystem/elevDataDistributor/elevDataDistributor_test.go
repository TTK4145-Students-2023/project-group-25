package elevDataDistributor

import (
	dt "project/commonDataTypes"
	"testing"
)

//Tests for elevDataDist

func TestDataDistributor(t *testing.T) {

	//Test 1:
	t.Run("Check output for external and local state input", func(t *testing.T) {
		//input and output channels
		allElevData_fromP2P := make(chan dt.AllElevDataJSON_withID)
		localElevData := make(chan dt.ElevDataJSON)
		HallOrderArray := make(chan [][2]bool)
		allElevData_toP2P := make(chan dt.AllElevDataJSON_withID)
		WorlView_toAssigner := make(chan dt.CostFuncInput)

		//start the distributor as a goroutine
		go DataDistributor(
			allElevData_fromP2P,
			localElevData,
			HallOrderArray,
			allElevData_toP2P,
			WorlView_toAssigner,
		)

		//mocking inputs
		input_localData := dt.ElevDataJSON{
			Behavior:    "Moving",
			Floor:       3,
			Direction:   "up",
			CabRequests: []bool{false, false, false, false},
		}

		//input data from P2P
		DataMatrix := make(dt.AllElevDataJSON)
		DataMatrix["ID1"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID2"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID3"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		input_datamatrix_withID := dt.AllElevDataJSON_withID{
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
		allElevData_fromP2P := make(chan dt.AllElevDataJSON_withID)
		localElevData := make(chan dt.ElevDataJSON)
		HallOrderArray := make(chan [][2]bool)
		allElevData_toP2P := make(chan dt.AllElevDataJSON_withID)
		WorlView_toAssigner := make(chan dt.CostFuncInput)

		//start the distributor as a goroutine
		go DataDistributor(
			allElevData_fromP2P,
			localElevData,
			HallOrderArray,
			allElevData_toP2P,
			WorlView_toAssigner,
		)

		//mocking inputs
		input_HallOrders := [][2]bool{{true, false}, {true, false}, {true, false}, {true, false}}

		//input data from P2P
		DataMatrix := make(dt.AllElevDataJSON)
		DataMatrix["ID1"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID2"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		DataMatrix["ID3"] = dt.ElevDataJSON{
			Behavior:    "Idle",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{true, false, true, false},
		}

		input_datamatrix_withID := dt.AllElevDataJSON_withID{
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
		t.Logf("\ninput datamtrix: \n%v\n input orderArray \n%v\n output datamatrix \n%v\n output WorlView: \n%v\n",
			input_datamatrix_withID, input_HallOrders, output_Datamatrix, output_WorldView)

	})

}
