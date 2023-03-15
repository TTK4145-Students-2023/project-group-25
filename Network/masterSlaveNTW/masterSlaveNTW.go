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

const BROADCAST_FREQ = 100 * time.Millisecond //ms

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

	broadCastTimer := time.NewTimer(1)
	masterSlaveRoleTimer := time.NewTimer(1)
	masterSlaveRoleTimer.Stop()
	ordersFromMasterTimer := time.NewTimer(1)
	ordersFromMasterTimer.Stop()

	MS_role := dt.MS_Slave
	ordersToSlaves := map[string][dt.N_FLOORS][2]bool{}
	ordersFromMaster := map[string][dt.N_FLOORS][2]bool{}

	for {
		select {
		case peerUpdate := <-peerUpdateChan:
			if newMS_Role := MS_Assigner(localIP, peerUpdate.Peers); newMS_Role != MS_role {
				MS_role = newMS_Role
				masterSlaveRoleTimer.Reset(1)
			}
		case ordersToSlaves = <-ordersToSlavesChan:
		case newOrdersFromMaster := <-receiveOrdersChan:
			if !reflect.DeepEqual(newOrdersFromMaster, ordersFromMaster) {
				switch MS_role {
				case dt.MS_Master:
				case dt.MS_Slave:
					ordersFromMaster = newOrdersFromMaster
					ordersFromMasterTimer.Reset(1)
				}
			}
		case <-broadCastTimer.C:
			switch MS_role {
			case dt.MS_Slave:
			case dt.MS_Master:
				transmittOrdersChan <- ordersToSlaves
				broadCastTimer.Reset(BROADCAST_FREQ)
			}
		case <-masterSlaveRoleTimer.C:
			select {
			case masterOrSlaveChan <- MS_role:
			default:
				masterSlaveRoleTimer.Reset(1)
			}
		case <-ordersFromMasterTimer.C:
			select {
			case ordersFromMasterChan <- ordersFromMaster:
			default:
				ordersFromMasterTimer.Reset(1)
			}
		}
	}
}
