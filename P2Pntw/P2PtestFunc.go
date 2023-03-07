package P2P

import (
	"Driver-go/localip"
	"fmt"
	"time"
)

// States for hall requests
const (
	STATE_new       requestState = 0
	STATE_confirmed requestState = 1
	STATE_none      requestState = 2
)

// Test functions for sending different datatypes
func MakeMat() RequestStateMatrix {
	input_ReqStatMatrix := make(RequestStateMatrix)
	input_ReqStatMatrix["ID1"] = singleNode_requestStates{{STATE_none, STATE_none},
		{STATE_new, STATE_none},
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
	return input_ReqStatMatrix
}

func MakeMsg() AllElevData_withID {
	localIP, _ := localip.LocalIP()
	AllElevDat := map[string]ElevData{
		"one": ElevData{
			Behavior:    "moving",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{false, false, false, true},
		},
		"two": ElevData{
			Behavior:    "idle",
			Floor:       0,
			Direction:   "stop",
			CabRequests: []bool{false, false, false, false},
		},
	}
	makeMsg := AllElevData_withID{
		localIP, AllElevDat,
	}
	return makeMsg
}

func TestP2P() {
	allElevData_fromDist := make(chan AllElevData_withID)
	ReqStateMatrix_fromOrderHandler := make(chan RequestStateMatrix)
	allElevData_fromNTW := make(chan AllElevData_withID)
	ReqStateMatrix_fromNTW := make(chan RequestStateMatrix)

	allElevData_toDist := make(chan AllElevData_withID)
	ReqStateMatrix_toOrderHandler := make(chan RequestStateMatrix)
	allElevData_toNTW := make(chan AllElevData_withID)
	ReqStateMatrix_toNTW := make(chan RequestStateMatrix)

	go P2Pntw(
		allElevData_fromDist,
		ReqStateMatrix_fromOrderHandler,
		allElevData_fromNTW,
		ReqStateMatrix_fromNTW,
		allElevData_toDist,
		ReqStateMatrix_toOrderHandler,
		allElevData_toNTW,
		ReqStateMatrix_toNTW,
	)

	timer := time.NewTimer(time.Second * 2)

	for {
		select {
		case tmp := <-ReqStateMatrix_fromNTW:
			fmt.Println(tmp)
			ReqStateMatrix_toOrderHandler <- tmp

		case <-timer.C:
			inputMat := MakeMat()
			ReqStateMatrix_fromNTW <- inputMat
			inputAllElevState := MakeMsg()
			allElevData_fromNTW <- inputAllElevState
			timer.Reset(time.Second * 2)
		}
	}
}
