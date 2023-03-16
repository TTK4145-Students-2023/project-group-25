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
	peerUpdateChan chan peers.PeerUpdate,
	ordersToSlavesChan <-chan []dt.SlaveOrders,
	ordersFromMasterChan chan<- []dt.SlaveOrders,
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

	MS_role := dt.MS_SLAVE
	ordersToSlaves := map[string][dt.N_FLOORS][2]bool{}
	ordersFromMaster := map[string][dt.N_FLOORS][2]bool{}

	for {
		select {
		case peerUpdate := <-peerUpdateChan:
			if newMS_Role := MS_Assigner(localIP, peerUpdate.Peers); newMS_Role != MS_role {
				MS_role = newMS_Role
				masterSlaveRoleTimer.Reset(1)
			}
		case newOrdersToSlaves := <-ordersToSlavesChan:
			ordersToSlaves = dt.SlaveOrdersSliceToMap(newOrdersToSlaves)
		case newOrdersFromMaster := <-receiveOrdersChan:
			if !reflect.DeepEqual(newOrdersFromMaster, ordersFromMaster) {
				switch MS_role {
				case dt.MS_MASTER:
				case dt.MS_SLAVE:
					ordersFromMaster = newOrdersFromMaster
					ordersFromMasterTimer.Reset(1)
				}
			}
		case <-broadCastTimer.C:
			broadCastTimer.Reset(dt.BROADCAST_RATE)
			switch MS_role {
			case dt.MS_SLAVE:
			case dt.MS_MASTER:
				transmittOrdersChan <- ordersToSlaves
			}
		case <-masterSlaveRoleTimer.C:
			select {
			case masterOrSlaveChan <- MS_role:
			default:
				masterSlaveRoleTimer.Reset(1)
			}
		case <-ordersFromMasterTimer.C:
			select {
			case ordersFromMasterChan <- dt.SlaveOrdersMapToSlice(ordersFromMaster):
			default:
				ordersFromMasterTimer.Reset(1)
			}
		}
	}
}
