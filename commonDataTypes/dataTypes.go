package dt

import (
	"time"
)

type AssignerBehaviour string

// ****** ELEVATOR CONSTANTS ********
const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

// ****** NETWORK CONSTANTS ********
const (
	MS_PORT                             = 15480
	PEER_LIST_PORT                      = 15489
	ORDERSTATE_PORT                     = 15488
	DATA_DISTRIBUTOR_PORT               = 15487
	BROADCAST_PERIOD      time.Duration = 100 * time.Millisecond
)

const (
	MASTER AssignerBehaviour = "master"
	SLAVE  AssignerBehaviour = "slave"
)

type ElevData struct {
	Behavior    string         `json:"behaviour"`
	Floor       int            `json:"floor"`
	Direction   string         `json:"direction"`
	CabRequests [N_FLOORS]bool `json:"cabRequests"`
}

type NodeInfo struct {
	IP   string   `json:"ip"`
	Data ElevData `json:"data"`
}

type AllNodeInfoWithSenderIP struct {
	SenderIP    string              `json:"senderIp"`
	AllNodeInfo map[string]ElevData `json:"allNodeInfo"`
}

type CostFuncInput struct {
	HallRequests [N_FLOORS][2]bool   `json:"hallRequests"`
	States       map[string]ElevData `json:"states"`
}

type CostFuncInputSlice struct {
	HallRequests [N_FLOORS][2]bool `json:"hallRequests"`
	States       []NodeInfo        `json:"states"`
}

type SlaveOrders struct {
	IP     string            `json:"slaveIp"`
	Orders [N_FLOORS][2]bool `json:"slaveOrders"`
}

func SlaveOrdersSliceToMap(slaveOrdersSlice []SlaveOrders) map[string][N_FLOORS][2]bool {
	slaveOrdersMap := map[string][N_FLOORS][2]bool{}
	for _, slaveOrders := range slaveOrdersSlice {
		slaveOrdersMap[slaveOrders.IP] = slaveOrders.Orders
	}
	return slaveOrdersMap
}

func SlaveOrdersMapToSlice(slaveOrdersMap map[string][N_FLOORS][2]bool) []SlaveOrders {
	slaveOrdersSlice := make([]SlaveOrders, len(slaveOrdersMap))
	for ip, slaveOrders := range slaveOrdersMap {
		slaveOrdersSlice = append(slaveOrdersSlice, SlaveOrders{IP: ip, Orders: slaveOrders})
	}
	return slaveOrdersSlice
}

func NodeInfoMapToSlice(nodesInfoMap map[string]ElevData) []NodeInfo {
	nodesInfoSlice := []NodeInfo{}
	for ip, elevData := range nodesInfoMap {
		nodesInfoSlice = append(nodesInfoSlice, NodeInfo{IP: ip, Data: elevData})
	}
	return nodesInfoSlice
}

func CostFuncInputSliceToMap(costFuncInputSlice CostFuncInputSlice) CostFuncInput {
	allNodeInfo := map[string]ElevData{}
	for _, nodeInfo := range costFuncInputSlice.States {
		allNodeInfo[nodeInfo.IP] = nodeInfo.Data
	}
	costFuncInput := CostFuncInput{HallRequests: costFuncInputSlice.HallRequests, States: allNodeInfo}
	return costFuncInput
}
