package P2P

import (
	"fmt"
	"project/Network/Utilities/bcast"
	PP "project/Network/printing"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

const BROADCAST_FREQ = 100 * time.Millisecond

func P2Pntw(localIP string,
	localWorldViewChan <-chan dt.AllElevDataJSON,
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

	localWorldView := dt.AllElevDataJSON{}
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

	RSM := ""
	WW := ""

	for {
		select {
		case localRequestStateMatrix = <-localRequestStateMatrixChan:
			RSM = PP.RSM_toString(localRequestStateMatrix)
			fmt.Printf(RSM + "/n" + WW)
		case localWorldView = <-localWorldViewChan:
			WW = PP.WW_toString(localWorldView)
			fmt.Printf(RSM + "/n" + WW)
		case newRequestStateMatrix := <-receiveRequestStateMatrix:
			senderData := newRequestStateMatrix.RequestMatrix
			senderIP := newRequestStateMatrix.IpAdress
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localRequestStateMatrix[senderIP]) {
				requestStateMatrix = newRequestStateMatrix
				reqStateMatrixTimer.Reset(1)
			}
		case newWorldView := <-receiveWorldView:
			senderData := newWorldView.AllData
			senderIP := newWorldView.ID
			if localIP != senderIP && !reflect.DeepEqual(senderData[senderIP], localWorldView[senderIP]) {
				worldView = newWorldView
				worldViewTimer.Reset(1)
			}
		case <-broadCastTimer.C:
			transmittWorldVeiw <- dt.AllElevDataJSON_withID{ID: localIP, AllData: localWorldView}
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
