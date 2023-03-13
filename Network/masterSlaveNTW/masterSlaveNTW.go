package masterSlaveNTW

import (
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/localip"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

// datatypes
type AssignedOrders map[string][]bool

func MasterSlaveNTW(

	peerUpdate_MS chan peers.PeerUpdate,
	localOrdersChan <-chan []byte,
	externalOrdersChan chan<- []byte,
	masterOrSlave chan dt.MasterSlaveRole,
) {
	var (
		receiveOrdersChan   = make(chan []byte)
		transmittOrdersChan = make(chan []byte)
	)

	localIP, _ := localip.LocalIP()

	go bcast.Receiver(15640, receiveOrdersChan)
	go bcast.Transmitter(15640, transmittOrdersChan)

	local_MS_role := dt.MS_Slave
	localOrders := []byte{}
	externalOrders := []byte{}

	for {
		select {
		case peerList := <-peerUpdate_MS:
			local_MS_role = MS_Assigner(localIP, peerList.Peers)
			masterOrSlave <- local_MS_role

		case localOrders = <-localOrdersChan:
		case newOrders := <-receiveOrdersChan:
			if !reflect.DeepEqual(newOrders, externalOrders) {
				externalOrders = newOrders
				switch local_MS_role {
				case dt.MS_Master:
				case dt.MS_Slave:
					externalOrdersChan <- externalOrders
				}
			}
		default:
			time.Sleep(time.Millisecond * 40)
		}
		switch local_MS_role {
		case dt.MS_Slave:
		case dt.MS_Master:
			transmittOrdersChan <- localOrders
		}
	}
}
