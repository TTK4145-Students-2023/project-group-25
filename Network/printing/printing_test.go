package printing

// import (
// 	"fmt"
// 	dt "project/commonDataTypes"
// 	"testing"
// )

// func Test_WW_toString(t *testing.T) {
// 	// Set up some test data.

// 	data := []dt.nodeInfo{}
// 	data["ID1"] = dt.ElevData{
// 		Behavior:    "Idle",
// 		Floor:       2,
// 		Direction:   "up",
// 		CabRequests: [dt.N_FLOORS]bool{true, false, true, false},
// 	}
// 	data["ID2"] = dt.ElevData{
// 		Behavior:    "Moving",
// 		Floor:       3,
// 		Direction:   "down",
// 		CabRequests: [dt.N_FLOORS]bool{false, true, false, true},
// 	}

// 	// WW_t := dt.CostFuncInput{
// 	// 	HallRequests: [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
// 	// 	States:       data,
// 	// }

// 	fmt.Printf(WW_toString(data))

// }

// func Test_RSM_toString(t *testing.T) {

// 	//mocking RSM
// 	input_ReqStatMatrix := make(dt.RequestStateMatrix)
// 	input_ReqStatMatrix["127.098.34"] =  [dt.N_FLOORS][2]dt.OrderState{{dt.STATE_NONE, STATE_NONE}, //dt.SingleNode_requestStates{{STATE_NONE, STATE_NONE},
// 		{STATE_NONE, STATE_NEW},
// 		{STATE_NONE, STATE_NONE},
// 		{STATE_NONE, STATE_NONE}}

// 	input_ReqStatMatrix["ID2"] = dt.SingleNode_requestStates{{STATE_NONE, STATE_NONE},
// 		{STATE_NONE, STATE_NONE},
// 		{STATE_NEW, STATE_NONE},
// 		{STATE_NONE, STATE_NONE}}

// 	input_ReqStatMatrix["ID3"] = dt.SingleNode_requestStates{{STATE_NONE, STATE_NONE},
// 		{STATE_NONE, STATE_CONFIRMED},
// 		{STATE_NONE, STATE_NONE},
// 		{STATE_NONE, STATE_NEW}}

// 	fmt.Printf(RSM_toString(input_ReqStatMatrix))
// }

// func Test_OrdersToString(t *testing.T) {

// 	role := dt.MS_Master
// 	sentOrders := make(map[string][dt.N_FLOORS][2]bool)
// 	sentOrders["127.0943.32"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
// 	sentOrders["ID2"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
// 	receivedOrders := make(map[string][dt.N_FLOORS][2]bool)
// 	receivedOrders["ID1"] = [dt.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
// 	receivedOrders["ID2"] = [dt.N_FLOORS][2]bool{{false, false}, {true, false}, {true, false}, {false, false}}

// 	ordersString := OrdersToString(role, sentOrders, receivedOrders)
// 	fmt.Println(ordersString)

// }
