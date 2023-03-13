package masterSlaveNTW

import (
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

// datatypes
type AssignedOrders map[string][]bool

func MasterSlaveNTW(localIP string,
	peerUpdate_MS chan peers.PeerUpdate,
	localOrdersChan <-chan map[string][][2]bool,
	externalOrdersChan chan<- map[string][][2]bool,
	masterOrSlave chan dt.MasterSlaveRole,
) {
	var (
		receiveOrdersChan   = make(chan map[string][][2]bool)
		transmittOrdersChan = make(chan map[string][][2]bool)
	)

	go bcast.Receiver(15660, receiveOrdersChan)
	go bcast.Transmitter(15660, transmittOrdersChan)

	local_MS_role := dt.MS_Slave
	localOrders := map[string][][2]bool{}
	externalOrders := map[string][][2]bool{}

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
