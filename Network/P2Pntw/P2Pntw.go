package P2P

import (
	"fmt"
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

const BROADCAST_FREQ = 100 //ms

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

	externalWorldView := dt.AllElevDataJSON_withID{}
	externalRequestStateMatrix := dt.RequestStateMatrix_with_ID{}

	//set timer
	timer := time.NewTimer(BROADCAST_FREQ * time.Millisecond)

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

			if localIP != newRequestStateMatrix.IpAdress && !reflect.DeepEqual(newRequestStateMatrix, externalRequestStateMatrix) {
				externalRequestStateMatrix = newRequestStateMatrix
				fmt.Printf("P2P, deadlock 1! ")
				externalRequestStateMatrixChan <- externalRequestStateMatrix
				fmt.Printf("... kidding, no P2P deadlock 1...\n ")
			}
		case newWorldView := <-receiveWorldView:
			if localIP != newWorldView.ID && !reflect.DeepEqual(newWorldView, externalWorldView) {

				externalWorldView = newWorldView
				fmt.Printf("P2P, deadlock 2! ")
				externalWorldViewChan <- externalWorldView
				fmt.Printf("... kidding, no P2P deadlock 2...\n ")
			}
		case <-timer.C:
			transmittWorldVeiw <- localWorldView
			transmittRequestStateMatrix <- dt.RequestStateMatrix_with_ID{IpAdress: localIP, RequestMatrix: localRequestStateMatrix}
			timer.Reset(BROADCAST_FREQ * time.Millisecond)
		}
	}
}
