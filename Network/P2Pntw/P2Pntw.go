package P2P

import (
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
)

func P2Pntw(
	// Receive channels
	allElevData_fromDist chan dt.AllElevDataJSON_withID,
	ReqStateMatrix_fromOrderHandler chan dt.RequestStateMatrix,

	// Sending channels
	allElevData_toDist chan dt.AllElevDataJSON_withID,
	ReqStateMatrix_toOrderHandler chan dt.RequestStateMatrix,

) {
	// Receive from NTW
	go bcast.Receiver(15647, allElevData_toDist)
	go bcast.Receiver(15648, ReqStateMatrix_toOrderHandler)

	// Send to NTW
	go bcast.Transmitter(15647, allElevData_fromDist)
	go bcast.Transmitter(15648, ReqStateMatrix_fromOrderHandler)
}
