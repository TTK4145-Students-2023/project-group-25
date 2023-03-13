package P2P

import (
	"fmt"
	"project/Network/Utilities/bcast"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

func P2Pntw(localIP string,
	localWorldViewChan <-chan dt.AllElevDataJSON_withID,
	localRequestStateMatrixChan <-chan dt.RequestStateMatrix,
	externalWorldViewChan chan<- dt.AllElevDataJSON_withID,
	externalRequestStateMatrixChan chan<- dt.RequestStateMatrix,
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
	externalRequestStateMatrix := dt.RequestStateMatrix{}

	//set timer
	timer1 := time.NewTimer(100 * time.Millisecond)

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
				fmt.Printf("______RSM recieved from P2P__________\n")
				fmt.Printf("Sender ID: %v\n", newRequestStateMatrix.IpAdress)
				fmt.Printf("Data: %v\n", newRequestStateMatrix.RequestMatrix)
				fmt.Printf("_________________________\n")

				externalRequestStateMatrix = newRequestStateMatrix.RequestMatrix
				externalRequestStateMatrixChan <- externalRequestStateMatrix
			}
		case newWorldView := <-receiveWorldView:
			if localIP != newWorldView.ID && !reflect.DeepEqual(newWorldView, externalWorldView) {

				externalWorldView = newWorldView
				externalWorldViewChan <- externalWorldView
			}
		case <-timer1.C:
			//fmt.Printf("timer ticked\n")
			transmittWorldVeiw <- localWorldView
			transmittRequestStateMatrix <- dt.RequestStateMatrix_with_ID{IpAdress: localIP, RequestMatrix: localRequestStateMatrix}
			timer1.Reset(100 * time.Millisecond)
		}

	}
}
