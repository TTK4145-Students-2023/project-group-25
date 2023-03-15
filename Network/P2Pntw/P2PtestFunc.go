package P2P

// import (
// 	dt "project/commonDataTypes"
// 	"project/Network/Utilities/localip"
// 	"fmt"
// 	"time"
// )

// // States for hall requests
// const (
// 	STATE_new       requestState = 0
// 	STATE_confirmed requestState = 1
// 	STATE_none      requestState = 2
// )

// // Test functions for sending different datatypes
// func MakeMatInput() RequestStateMatrix {
// 	input_ReqStatMatrix := make(RequestStateMatrix)
// 	input_ReqStatMatrix["ID1"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_new, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}

// 	input_ReqStatMatrix["ID2"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}

// 	input_ReqStatMatrix["ID3"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}
// 	return input_ReqStatMatrix
// }

// func MakeMatOutput() RequestStateMatrix {
// 	input_ReqStatMatrix := make(RequestStateMatrix)
// 	input_ReqStatMatrix["10.100.23.30"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_new, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}

// 	input_ReqStatMatrix["10.100.23.31"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}

// 	input_ReqStatMatrix["10.100.23.34"] = singleNode_requestStates{{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none},
// 		{STATE_none, STATE_none}}
// 	return input_ReqStatMatrix
// }

// func MakeMsgInput() dt.AllElevDataJSON_withID {
// 	AllElevDat := map[string]dt.ElevDataJSON{
// 		"one": dt.ElevDataJSON{
// 			Behavior:    "moving",
// 			Floor:       2,
// 			Direction:   "up",
// 			CabRequests: []bool{false, false, false, true},
// 		},
// 		"two": dt.ElevDataJSON{
// 			Behavior:    "idle",
// 			Floor:       0,
// 			Direction:   "stop",
// 			CabRequests: []bool{false, false, false, false},
// 		},
// 	}
// 	makeMsg := dt.AllElevDataJSON_withID{
// 		localIP, AllElevDat,
// 	}
// 	return makeMsg
// }

// func MakeMsgOutput() dt.AllElevDataJSON_withID {
// 	AllElevDat := map[string]dt.ElevDataJSON{
// 		"one": dt.ElevDataJSON{
// 			Behavior:    "idle",
// 			Floor:       2,
// 			Direction:   "up",
// 			CabRequests: []bool{false, false, false, true},
// 		},
// 		"two": dt.ElevDataJSON{
// 			Behavior:    "idle",
// 			Floor:       0,
// 			Direction:   "stop",
// 			CabRequests: []bool{false, false, false, false},
// 		},
// 		"three": dt.ElevDataJSON{
// 			Behavior:    "idle",
// 			Floor:       8,
// 			Direction:   "stop",
// 			CabRequests: []bool{false, false, false, false},
// 		},
// 	}
// 	makeMsg := dt.AllElevDataJSON_withID{
// 		localIP, AllElevDat,
// 	}
// 	return makeMsg
// }

// func TestP2P() {
// 	allElevData_fromDist := make(chan dt.AllElevDataJSON_withID, 1)
// 	ReqStateMatrix_fromOrderHandler := make(chan RequestStateMatrix, 1)
// 	allElevData_fromNTW := make(chan dt.AllElevDataJSON_withID, 1)
// 	ReqStateMatrix_fromNTW := make(chan RequestStateMatrix, 1)

// 	allElevData_toDist := make(chan dt.AllElevDataJSON_withID, 1)
// 	ReqStateMatrix_toOrderHandler := make(chan RequestStateMatrix, 1)
// 	allElevData_toNTW := make(chan dt.AllElevDataJSON_withID, 1)
// 	ReqStateMatrix_toNTW := make(chan RequestStateMatrix, 1)

// 	go P2Pntw(
// 		allElevData_fromDist,
// 		ReqStateMatrix_fromOrderHandler,
// 		allElevData_fromNTW,
// 		ReqStateMatrix_fromNTW,
// 		allElevData_toDist,
// 		ReqStateMatrix_toOrderHandler,
// 		allElevData_toNTW,
// 		ReqStateMatrix_toNTW,
// 	)

// 	timer := time.NewTimer(time.Second * 1)

// 	for {
// 		select {
// 		case <-timer.C:
// 			inputMat := MakeMatInput()
// 			inputAllElevState := MakeMsgInput()
// 			outputMat := MakeMatOutput()
// 			outputAllElevData := MakeMsgOutput()
// 			ReqStateMatrix_fromNTW <- inputMat
// 			allElevData_fromNTW <- inputAllElevState
// 			ReqStateMatrix_fromOrderHandler <- outputMat
// 			allElevData_fromDist <- outputAllElevData
// 			timer.Reset(time.Second * 1)

// 		case MatToOrderHandler := <-ReqStateMatrix_toOrderHandler:
// 			fmt.Println(MatToOrderHandler)
// 		case allElevDataDist := <-allElevData_toDist:
// 			fmt.Println(allElevDataDist)
// 		case MatToNTW := <-ReqStateMatrix_toNTW:
// 			fmt.Println(MatToNTW)
// 		case allElevToNTW := <-allElevData_toNTW:
// 			fmt.Println(allElevToNTW)
// 		}
// 	}

// }
