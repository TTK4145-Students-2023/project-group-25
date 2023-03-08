package oassign

import (
	"encoding/json"
	"fmt"
	"os/exec"
	elevfsm "project/localElevator/elev_fsm"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [elevfsm.N_FLOORS]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [elevfsm.N_FLOORS][2]bool `json:"hallRequests"`
	States       map[string]HRAElevState   `json:"states"`
}

type OrderAssignerBehaviour int

const (
	MS_Master OrderAssignerBehaviour = 0
	MS_Slave  OrderAssignerBehaviour = 1
)

func OrderAssigner(OrderAssignerBehaviourChan <-chan OrderAssignerBehaviour,
	localIpAdressChan <-chan string, // Chanel where local IP-adress is fetched
	ordersFromDistributor <-chan HRAInput, // Input from order distributor
	ordersFromMaster <-chan []byte, // Input read from Master-Slave network module
	ordersToSlaves chan<- []byte, // Input written to Master-Slave network module
	localOrders chan<- [elevfsm.N_FLOORS][2]bool) { // Input to local Elevator FSM

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	localIpAdress := ""
	assignerBehaviour := MS_Slave
	for {
		select {
		case localIpAdress = <-localIpAdressChan:
		case assignerBehaviour = <-OrderAssignerBehaviourChan:
		case input := <-ordersFromDistributor:
			fmt.Printf("Inside recc data to from distributor\n")
			switch assignerBehaviour {
			case MS_Slave:
			case MS_Master:
				jsonBytes, err := json.Marshal(input)
				if err != nil {
					fmt.Println("json.Marshal error: ", err)
					return
				}
				ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
				if err != nil {
					fmt.Println("exec.Command error: ", err)
					fmt.Println(string(ret))
					return
				}

				output := map[string][elevfsm.N_FLOORS][2]bool{}
				err = json.Unmarshal(ret, &output)
				if err != nil {
					fmt.Println("json.Unmarshal error: ", err)
					return
				}
				if localHallOrders, ok := output[localIpAdress]; ok {
					localOrders <- localHallOrders
				}

				ordersToSlaves <- ret
			}
		case input := <-ordersFromMaster:
			switch assignerBehaviour {
			case MS_Master:
			case MS_Slave:
				output := map[string][elevfsm.N_FLOORS][2]bool{}
				err := json.Unmarshal(input, &output)
				if err != nil {
					fmt.Println("json.Unmarshal error: ", err)
					return
				}
				if localHallOrders, ok := output[localIpAdress]; ok {
					localOrders <- localHallOrders
				}
			}
		}
	}
}
