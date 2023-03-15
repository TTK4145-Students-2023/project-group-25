package printing

import (
	"fmt"
	dt "project/commonDataTypes"
	"testing"
)

func Test_WW_toString(t *testing.T) {
	// Set up some test data.

	data := make(dt.AllElevDataJSON)
	data["ID1"] = dt.ElevDataJSON{
		Behavior:    "Idle",
		Floor:       2,
		Direction:   "up",
		CabRequests: [dt.N_FLOORS]bool{true, false, true, false},
	}
	data["ID2"] = dt.ElevDataJSON{
		Behavior:    "Moving",
		Floor:       3,
		Direction:   "down",
		CabRequests: [dt.N_FLOORS]bool{false, true, false, true},
	}

	// WW_t := dt.CostFuncInput{
	// 	HallRequests: [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
	// 	States:       data,
	// }

	fmt.Println(WW_toString(data))

}

func Test_RSM_toString(t *testing.T) {

	//mocking RSM
	input_ReqStatMatrix := make(dt.RequestStateMatrix)
	input_ReqStatMatrix["127.098.34"] = dt.SingleNode_requestStates{{STATE_none, STATE_none},
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

	fmt.Printf(RSM_toString(input_ReqStatMatrix))
}

func Test_OrdersToString(t *testing.T) {

	role := dt.MS_Master
	sentOrders := make(map[string][dt.N_FLOORS][2]bool)
	sentOrders["127.0943.32"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	sentOrders["ID2"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	receivedOrders := make(map[string][dt.N_FLOORS][2]bool)
	receivedOrders["ID1"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	receivedOrders["ID2"] = [dt.N_FLOORS][2]bool{{false, false}, {true, false}, {true, false}, {false, false}}

	ordersString := OrdersToString(role, sentOrders, receivedOrders)
	fmt.Println(ordersString)

}
