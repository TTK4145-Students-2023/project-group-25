package orderStateHandler

import (
	"Driver-go/elevio"
	"reflect"
	"testing"
)

// testOrderStatHanlder tests the statehandler function
func TestOrderStateHandler(t *testing.T) {
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

		//setting timeout to ensure completion of test
		//timeout := time.After(1 * time.Second)

		//Possible inputs

		// Hallorder := elevio.ButtonEvent{
		// 	Floor:  2,
		// 	Button: elevio.BT_HallUp,
		// }

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

		// exec_btn_1 := elevio.ButtonEvent{
		// 	Floor:  2,
		// 	Button: elevio.BT_HallUp,
		// }

		// exec_btn_2 := elevio.ButtonEvent{
		// 	Floor:  2,
		// 	Button: elevio.BT_HallUp,
		// }

		// executedOrders := []elevio.ButtonEvent{exec_btn_1, exec_btn_2}

		// Send inputs on channels
		// HallBtnPress <- Hallorder

		ReqStateMatrix_fromP2P <- input_ReqStatMatrix

		// orderExecuted <- executedOrders

		output_mat := <-ReqStateMatrix_toP2P
		output_arr := <-HallOrderArray

		if !reflect.DeepEqual(1, 0) {
			t.Errorf("mat: %v , arr: %v", output_arr, output_mat)
		}
	})

}
