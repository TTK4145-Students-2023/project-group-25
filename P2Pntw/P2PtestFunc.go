package P2P

import (
	"fmt"
	"time"
)

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
