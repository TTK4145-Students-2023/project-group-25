package masterSlaveNTW

import (
	"fmt"
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"time"
)

// datatypes
type AssignedOrders map[string][]bool

const BROADCAST_FREQ = 100 //ms

func MasterSlaveNTW(localIP string,
	peerUpdateChan chan peers.PeerUpdate,
	ordersToSlavesChan <-chan map[string][dt.N_FLOORS][2]bool,
	ordersFromMasterChan chan<- map[string][dt.N_FLOORS][2]bool,
	masterOrSlaveChan chan dt.MasterSlaveRole,
) {
	var (
		receiveOrdersChan   = make(chan map[string][dt.N_FLOORS][2]bool)
		transmittOrdersChan = make(chan map[string][dt.N_FLOORS][2]bool)
	)

	go bcast.Receiver(15660, receiveOrdersChan)
	go bcast.Transmitter(15660, transmittOrdersChan)

	timer := time.NewTimer(BROADCAST_FREQ * time.Millisecond)
	MS_role := dt.MS_Slave
	ordersToSlaves := map[string][dt.N_FLOORS][2]bool{}
	ordersFromMaster := map[string][dt.N_FLOORS][2]bool{}

	for {
		select {
		case peerUpdate := <-peerUpdateChan:
			if newRole := MS_Assigner(localIP, peerUpdate.Peers); newRole != MS_role {
				MS_role = newRole
				fmt.Printf("MS, deadlock 1! ")
				masterOrSlaveChan <- MS_role
				fmt.Printf("... kidding, no MS deadlock 1...\n ")
			}
		case ordersToSlaves = <-ordersToSlavesChan:
		case newOrdersFromMaster := <-receiveOrdersChan:
			if !reflect.DeepEqual(newOrdersFromMaster, ordersFromMaster) {
				ordersFromMaster = newOrdersFromMaster
				switch MS_role {
				case dt.MS_Master:
				case dt.MS_Slave:
					fmt.Printf("MS, deadlock 2! ")
					ordersFromMasterChan <- ordersFromMaster
					fmt.Printf("... kidding, no MS deadlock 2...\n ")
				}
			}
		case <-timer.C:
			switch MS_role {
			case dt.MS_Slave:
			case dt.MS_Master:
				transmittOrdersChan <- ordersToSlaves
				timer.Reset(BROADCAST_FREQ * time.Millisecond)
			}
		}
	}
}
