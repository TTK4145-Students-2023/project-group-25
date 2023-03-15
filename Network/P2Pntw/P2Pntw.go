package P2P

import (
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

const BROADCAST_FREQ = 100 * time.Millisecond

func P2Pntw(localIP string,
	localWorldViewChan <-chan dt.AllElevDataJSON_withID,
	localRequestStateMatrixChan <-chan dt.RequestStateMatrix,
	externalWorldViewChan chan<- dt.AllElevDataJSON_withID,
	externalRequestStateMatrixChan chan<- dt.RequestStateMatrix_with_ID,
) {
	var (
		transmittWorldVeiw          = make(chan dt.AllElevDataJSON_withID)
		transmittRequestStateMatrix = make(chan dt.RequestStateMatrix_with_ID)
		receiveWorldView            = make(chan dt.AllElevDataJSON_withID)
		receiveRequestStateMatrix   = make(chan dt.RequestStateMatrix_with_ID)
	)

	localWorldView := dt.AllElevDataJSON_withID{}
	localRequestStateMatrix := dt.RequestStateMatrix{}

	worldView := dt.AllElevDataJSON_withID{}
	requestStateMatrix := dt.RequestStateMatrix_with_ID{}

	//set timer
	broadCastTimer := time.NewTimer(BROADCAST_FREQ)
	reqStateMatrixTimer := time.NewTimer(1)
	reqStateMatrixTimer.Stop()
	worldViewTimer := time.NewTimer(1)
	worldViewTimer.Stop()

	// Receive from NTW
	go bcast.Receiver(15667, receiveWorldView)
	go bcast.Receiver(15668, receiveRequestStateMatrix)

	// Send to NTW
	go bcast.Transmitter(15667, transmittWorldVeiw)
	go bcast.Transmitter(15668, transmittRequestStateMatrix)

	for {
		select {
		case localRequestStateMatrix = <-localRequestStateMatrixChan:
		case localWorldView = <-localWorldViewChan:
		case newRequestStateMatrix := <-receiveRequestStateMatrix:
			if localIP != newRequestStateMatrix.IpAdress && !reflect.DeepEqual(newRequestStateMatrix, requestStateMatrix) {
				requestStateMatrix = newRequestStateMatrix
				reqStateMatrixTimer.Reset(1)
			}
		case newWorldView := <-receiveWorldView:
			if localIP != newWorldView.ID && !reflect.DeepEqual(newWorldView, worldView) {
				worldView = newWorldView
				worldViewTimer.Reset(1)
			}
		case <-broadCastTimer.C:
			transmittWorldVeiw <- localWorldView
			transmittRequestStateMatrix <- dt.RequestStateMatrix_with_ID{IpAdress: localIP, RequestMatrix: localRequestStateMatrix}
			broadCastTimer.Reset(BROADCAST_FREQ)
		case <-worldViewTimer.C:
			select {
			case externalWorldViewChan <- worldView:
			default:
				worldViewTimer.Reset(1)
			}
		case <-reqStateMatrixTimer.C:
			select {
			case externalRequestStateMatrixChan <- requestStateMatrix:
			default:
				reqStateMatrixTimer.Reset(1)
			}
		}
	}
}
