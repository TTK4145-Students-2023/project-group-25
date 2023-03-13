package P2P

import (
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

func P2Pntw(
	localWorldViewChan <-chan dt.AllElevDataJSON_withID,
	localRequestStateMatrixChan <-chan dt.RequestStateMatrix,
	externalWorldViewChan chan<- dt.AllElevDataJSON_withID,
	externalRequestStateMatrixChan chan<- dt.RequestStateMatrix,

) {
	var (
		transmittWorldVeiw          = make(chan<- dt.AllElevDataJSON_withID)
		transmittRequestStateMatrix = make(chan<- dt.RequestStateMatrix)
		receiveWorldView            = make(<-chan dt.AllElevDataJSON_withID)
		receiveRequestStateMatrix   = make(<-chan dt.RequestStateMatrix)
	)

	localWorldView := dt.AllElevDataJSON_withID{}
	localRequestStateMatrix := dt.RequestStateMatrix{}

	externalWorldView := dt.AllElevDataJSON_withID{}
	externalRequestStateMatrix := dt.RequestStateMatrix{}

	// Receive from NTW
	go bcast.Receiver(15647, receiveWorldView)
	go bcast.Receiver(15648, receiveRequestStateMatrix)

	// Send to NTW
	go bcast.Transmitter(15647, transmittWorldVeiw)
	go bcast.Transmitter(15648, transmittRequestStateMatrix)

	for {
		select {
		case localRequestStateMatrix = <-localRequestStateMatrixChan:
		case localWorldView = <-localWorldViewChan:
		case newRequestStateMatrix := <-receiveRequestStateMatrix:
			if !reflect.DeepEqual(newRequestStateMatrix, externalRequestStateMatrix) {
				externalRequestStateMatrix = newRequestStateMatrix
				externalRequestStateMatrixChan <- externalRequestStateMatrix
			}
		case newWorldView := <-receiveWorldView:
			if !reflect.DeepEqual(newWorldView, externalWorldView) {
				externalWorldView = newWorldView
				externalWorldViewChan <- externalWorldView
			}
		default:
			time.Sleep(time.Millisecond * 40)
		}
		transmittWorldVeiw <- localWorldView
		transmittRequestStateMatrix <- localRequestStateMatrix
	}
}
