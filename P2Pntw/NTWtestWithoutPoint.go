package P2P

// States for hall requests
const (
	STATE_new       requestState = 0
	STATE_confirmed requestState = 1
	STATE_none      requestState = 2
)

func makeMat() RequestStateMatrix {
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
	return input_ReqStatMatrix
}

// testOrderStatHanlder tests the P2Pntw function
func TestP2Pntw() {
	//Test Case 1: New ReqMatrix from P2P
	//t.Run("New hall button order is received", func(t *testing.T) {
	//input and output channels
	allElevData_fromDist := make(chan AllElevData_withID)
	ReqStateMatrix_fromOrderHandler := make(chan RequestStateMatrix)
	allElevData_fromNTW := make(chan AllElevData_withID)
	ReqStateMatrix_fromNTW := make(chan RequestStateMatrix)

	allElevData_toDist := make(chan AllElevData_withID)
	ReqStateMatrix_toOrderHandler := make(chan RequestStateMatrix)
	allElevData_toNTW := make(chan AllElevData_withID)
	ReqStateMatrix_toNTW := make(chan RequestStateMatrix)

	//Make and send input on channel
	inputMat := makeMat()
	ReqStateMatrix_toOrderHandler <- inputMat

	go P2Pntw(
		allElevData_fromDist,
		ReqStateMatrix_fromOrderHandler,
		allElevData_fromNTW,
		ReqStateMatrix_fromNTW,
		allElevData_toNTW,
		ReqStateMatrix_toNTW,
		allElevData_toDist,
		ReqStateMatrix_toOrderHandler,
	)
}
