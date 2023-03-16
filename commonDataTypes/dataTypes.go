package dt

import (
	"time"
)

type MasterSlaveRole string
type OrderState string

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

const (
	MS_MASTER      MasterSlaveRole = "master"
	MS_SLAVE       MasterSlaveRole = "slave"
	BROADCAST_RATE time.Duration   = 100 * time.Millisecond
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
	SenderIP    string     `json:"senderIp"`
	AllNodeInfo []NodeInfo `json:"allNodeInfo"`
}

type CostFuncInput struct {
	HallRequests [N_FLOORS][2]bool   `json:"hallRequests"`
	States       map[string]ElevData `json:"states"`
}

type CostFuncInputSlice struct {
	HallRequests [N_FLOORS][2]bool `json:"hallRequests"`
	States       []NodeInfo        `json:"states"`
}

type NodeOrderStates struct {
	IP          string                  `json:"ip"`
	OrderStates [N_FLOORS][2]OrderState `json:"orderStates"`
}

type AllNOS_WithSenderIP struct {
	SenderIP string            `json:"ip"`
	AllNOS   []NodeOrderStates `json:"allNOS"`
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

func NOSSliceToMap(NOSSlice []NodeOrderStates) map[string][N_FLOORS][2]OrderState {
	NOSMap := map[string][N_FLOORS][2]OrderState{}
	for _, NodeOrderStates := range NOSSlice {
		NOSMap[NodeOrderStates.IP] = NodeOrderStates.OrderStates
	}
	return NOSMap
}

func NOSMapToSlice(NOSMap map[string][N_FLOORS][2]OrderState) []NodeOrderStates {
	NOSSlice := []NodeOrderStates{}
	for ip, orderStates := range NOSMap {
		NOSSlice = append(NOSSlice, NodeOrderStates{IP: ip, OrderStates: orderStates})
	}
	return NOSSlice
}

func NodeInfoSliceToMap(nodesInfoSlice []NodeInfo) map[string]ElevData {
	nodesInfoMap := map[string]ElevData{}
	for _, nodeInfo := range nodesInfoSlice {
		nodesInfoMap[nodeInfo.IP] = nodeInfo.Data
	}
	return nodesInfoMap
}

func NodeInfoMapToSlice(nodesInfoMap map[string]ElevData) []NodeInfo {
	nodesInfoSlice := []NodeInfo{}
	for ip, elevData := range nodesInfoMap {
		nodesInfoSlice = append(nodesInfoSlice, NodeInfo{IP: ip, Data: elevData})
	}
	return nodesInfoSlice
}

func CostFuncInputToSlice(costFuncInput CostFuncInput) CostFuncInputSlice {
	allNodeInfo := []NodeInfo{}
	costFuncInputSlice := CostFuncInputSlice{HallRequests: costFuncInput.HallRequests, States: allNodeInfo}
	for ip, elevData := range costFuncInput.States {
		allNodeInfo = append(allNodeInfo, NodeInfo{IP: ip, Data: elevData})
	}
	return costFuncInputSlice
}

func SliceToCostFuncInput(costFuncInputSlice CostFuncInputSlice) CostFuncInput {
	allNodeInfo := map[string]ElevData{}
	for _, nodeInfo := range costFuncInputSlice.States {
		allNodeInfo[nodeInfo.IP] = nodeInfo.Data
	}
	costFuncInput := CostFuncInput{HallRequests: costFuncInputSlice.HallRequests, States: allNodeInfo}
	return costFuncInput
}
