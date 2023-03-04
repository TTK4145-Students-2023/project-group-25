package orderStateHandler

import (
	"Driver-go/elevio"
	"reflect"
	"testing"
)

// testOrderStatHanlder tests the statehandler function
func TestOrderStateHandler(t *testing.T) {
	//Test Case 1: New ReqMatrix from P2P
	t.Run("New hall button order is received", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan RequestStateMatrix)

		// Start the orderStateHandler as a goroutine
		go orderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
		)

		//Mocking request matrix input from P2P
		input_ReqStatMatrix := make(RequestStateMatrix)
		input_ReqStatMatrix["ID1"] = singleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID2"] = singleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		input_ReqStatMatrix["ID3"] = singleNode_requestStates{{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none},
			{STATE_none, STATE_none}}

		//Send input on channel
		ReqStateMatrix_fromP2P <- input_ReqStatMatrix

		//read output on channel
		output_mat := <-ReqStateMatrix_toP2P
		output_arr := <-HallOrderArray

		//assert
		if !reflect.DeepEqual(1, 0) {
			t.Errorf("mat: %v , arr: %v", output_arr, output_mat)
		}
	})

	//Test case 2: New hall button press
	t.Run("New hall button order is received", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan RequestStateMatrix)

		// Start the orderStateHandler as a goroutine
		go orderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
		)

		//Mocking Hall btn press
		Hallorder := elevio.ButtonEvent{
			Floor:  2,
			Button: elevio.BT_HallUp,
		}

		//Sending input on channel
		HallBtnPress <- Hallorder

		//Reading output on channel

		//assert
		if !reflect.DeepEqual(1, 0) {
			t.Errorf("assert error msg ")
		}

	})

	//Test case 3: Input: executed order array
	t.Run("New Executed order array is received", func(t *testing.T) {
		//input and output channels
		ReqStateMatrix_fromP2P := make(chan RequestStateMatrix)
		HallBtnPress := make(chan elevio.ButtonEvent)
		orderExecuted := make(chan []elevio.ButtonEvent)
		HallOrderArray := make(chan [][2]bool)
		ReqStateMatrix_toP2P := make(chan RequestStateMatrix)

		// Start the orderStateHandler as a goroutine
		go orderStateHandler(
			ReqStateMatrix_fromP2P,
			HallBtnPress,
			orderExecuted,
			HallOrderArray,
			ReqStateMatrix_toP2P,
		)

		//Mock executed orders
		exec_btn_1 := elevio.ButtonEvent{
			Floor:  2,
			Button: elevio.BT_HallUp,
		}
		exec_btn_2 := elevio.ButtonEvent{
			Floor:  2,
			Button: elevio.BT_HallDown,
		}
		executedOrders := []elevio.ButtonEvent{exec_btn_1, exec_btn_2}

		//Send inputs on channel
		orderExecuted <- executedOrders

		//Read outout on channel

		//assert
		if !reflect.DeepEqual(1, 0) {
			t.Errorf("assert error msg ")
		}
	})

}
