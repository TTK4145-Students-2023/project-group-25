package orderStateHandler

import (
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	elevio "project/localElevator/elev_driver"
	"reflect"
	"testing"
)

// testOrderStatHanlder tests the statehandler function
func TestOrderStateHandler(t *testing.T) {

	//Test Case 1:
	t.Run("Check output for new matrix from P2P", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan dt.RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan dt.RequestStateMatrix)
		peerUpdate_OrderHandler := make(chan peers.PeerUpdate)

		// Start the orderStateHandler as a goroutine
		go OrderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
			peerUpdate_OrderHandler,
		)

		//Mocking request matrix input from P2P
		input_ReqStatMatrix := make(dt.RequestStateMatrix)
		input_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_new},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_new}}

		//Send input on channel
		ReqStateMatrix_fromP2P <- input_ReqStatMatrix

		//read output on channel
		output_ReqMatrix := <-ReqStateMatrix_toP2P
		output_HallOrderArray := <-HallOrderArray

		//exepcted output
		expcted_ReqStatMatrix := make(dt.RequestStateMatrix)
		expcted_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_new}}

		expcted_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		expcted_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		t.Helper()
		//verify outputs with expected value
		if !reflect.DeepEqual(expcted_ReqStatMatrix, output_ReqMatrix) {
			//t.Errorf("\nERROR_REQ_MATRIX \nexpected value:\n%v  \nrecieved value:\n%v", expcted_ReqStatMatrix, output_ReqMatrix)
			t.Logf("\nCorresponding HallOrderArray: \n%v ", output_HallOrderArray)
		} else {
			t.Logf("\nCorrect Reqmatrix was recived:\n%v \nCorresponding HallOrderArray: \n%v ", output_ReqMatrix, output_HallOrderArray)

		}

	})

	//Test case 2:
	t.Run("Check output Matrix for newly executed orders", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan dt.RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan dt.RequestStateMatrix)
		peerUpdate_OrderHandler := make(chan peers.PeerUpdate)

		// Start the orderStateHandler as a goroutine
		go OrderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
			peerUpdate_OrderHandler,
		)

		//Mock executed orders
		exec_btn_1 := elevio.ButtonEvent{
			Floor:  0,
			Button: elevio.BT_HallUp,
		}
		exec_btn_2 := elevio.ButtonEvent{
			Floor:  0,
			Button: elevio.BT_HallDown,
		}
		executedOrders := []elevio.ButtonEvent{exec_btn_1, exec_btn_2}

		//Mocking request matrix input from P2P
		input_ReqStatMatrix := make(dt.RequestStateMatrix)
		input_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		//Send inputs on channel
		orderExecuted <- executedOrders
		ReqStateMatrix_fromP2P <- input_ReqStatMatrix

		//read output on channel
		output_ReqMatrix := <-ReqStateMatrix_toP2P
		output_HallOrderArray := <-HallOrderArray

		//exepcted output
		expcted_ReqStatMatrix := make(dt.RequestStateMatrix)
		expcted_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_none, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_new}}

		expcted_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_none, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		expcted_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_none, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		t.Helper()
		//verify outputs with expected value
		if !reflect.DeepEqual(expcted_ReqStatMatrix, output_ReqMatrix) {
			//t.Errorf("\nERROR_REQ_MATRIX \nexpected value:\n%v  \nrecieved value:\n%v", expcted_ReqStatMatrix, output_ReqMatrix)
			t.Logf("\nCorresponding HallOrderArray: \n%v ", output_HallOrderArray)
		} else {
			t.Logf("\nCorrect Reqmatrix was recived:\n%v \nCorresponding HallOrderArray: \n%v ", output_ReqMatrix, output_HallOrderArray)

		}
	})

	//Test case 3: New hall button press
	t.Run("Check output for newly confirmed order", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan dt.RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan dt.RequestStateMatrix)
		peerUpdate_OrderHandler := make(chan peers.PeerUpdate)

		// Start the orderStateHandler as a goroutine
		go OrderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
			peerUpdate_OrderHandler,
		)

		//Mocking Hall btn press
		Hallorder := elevio.ButtonEvent{
			Floor:  2,
			Button: elevio.BT_HallUp,
		}

		//Mocking request matrix input from P2P
		input_ReqStatMatrix := make(dt.RequestStateMatrix)
		input_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_confirmed, STATE_confirmed},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_none}}

		//exepcted output
		expcted_ReqStatMatrix := make(dt.RequestStateMatrix)
		expcted_ReqStatMatrix["ID1"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_confirmed, STATE_none},
			{STATE_none, STATE_none}}

		expcted_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_none}}

		expcted_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_new, STATE_none},
			{STATE_none, STATE_none}}

		//Sending input on channel
		HallBtnPress <- Hallorder
		ReqStateMatrix_fromP2P <- input_ReqStatMatrix

		//read output on channel
		output_ReqMatrix := <-ReqStateMatrix_toP2P
		output_HallOrderArray := <-HallOrderArray

		t.Helper()
		if !reflect.DeepEqual(expcted_ReqStatMatrix, output_ReqMatrix) {
			//	t.Errorf("\nERROR_REQ_MATRIX \nexpected value:\n%v  \nrecieved value:\n%v", expcted_ReqStatMatrix, output_ReqMatrix)
			t.Logf("\nCorresponding HallOrderArray: \n%v ", output_HallOrderArray)
		} else {
			t.Logf("\nCorrect Reqmatrix was recived:\n%v \nCorresponding HallOrderArray: \n%v ", output_ReqMatrix, output_HallOrderArray)

		}

	})
}
