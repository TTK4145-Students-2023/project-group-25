package masterSlaveNTW

import (
	"project/Network/Utilities/bcast"
	"project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"strconv"
	"strings"
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

		MS_role          = dt.MS_SLAVE
		ordersToSlaves   = map[string][dt.N_FLOORS][2]bool{}  	// orders to be sent  --  store as slice?
		ordersFromMaster = map[string][dt.N_FLOORS][2]bool{}  	// received orders   --  store as slice?

		broadCastTimer        = time.NewTimer(1)
		masterSlaveRoleTimer  = time.NewTimer(time.Hour)     	// roleAssignerTrigger
		ordersFromMasterTimer = time.NewTimer(time.Hour)		// received orders to O.ASS
	)
	masterSlaveRoleTimer.Stop()
	ordersFromMasterTimer.Stop()

	go bcast.Receiver(dt.MS_PORT, receiveOrdersChan)
	go bcast.Transmitter(dt.MS_PORT, transmittOrdersChan)

	for {
		select {
		case peerUpdate := <-peerUpdateChan:
			if newMS_Role := assignRole(localIP, peerUpdate.Peers); newMS_Role != MS_role {
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
		case <-masterSlaveRoleTimer.C:
			select {
			case masterOrSlaveChan <- MS_role: 		// reasign roles 
			default:
				masterSlaveRoleTimer.Reset(1)
			}
		case <-ordersFromMasterTimer.C:
			select {
			case ordersFromMasterChan <- dt.SlaveOrdersMapToSlice(ordersFromMaster):  // send to O.Ass
			default:
				ordersFromMasterTimer.Reset(1)
			}
		case <-broadCastTimer.C:
			broadCastTimer.Reset(dt.BROADCAST_PERIOD)
			switch MS_role {
			case dt.MS_SLAVE:
			case dt.MS_MASTER:
				transmittOrdersChan <- ordersToSlaves
			}
		}
	}
}

func assignRole(localIP string, peers []string) dt.MasterSlaveRole {
	if len(peers) == 0 {
		return dt.MS_MASTER
	}

	localIPArr := strings.Split(localIP, ".")
	LocalLastByte, _ := strconv.Atoi(localIPArr[len(localIPArr)-1])

	maxIP := LocalLastByte
	for _, externalIP := range peers {
		externalIPArr := strings.Split(externalIP, ".")
		externalLastByte, _ := strconv.Atoi(externalIPArr[len(externalIPArr)-1])
		if externalLastByte > maxIP {
			maxIP = externalLastByte
		}
	}

	if maxIP <= LocalLastByte {
		return dt.MS_MASTER
	}
	return dt.MS_SLAVE
}
