package oassign

import (
	"encoding/json"
	"fmt"
	"os/exec"
	dt "project/commonDataTypes"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func OrderAssigner(localIP string,
	masterSlaveRoleChan <-chan dt.MasterSlaveRole,
	ordersFromDistributor <-chan dt.CostFuncInput, // Input from order distributor
	ordersFromMaster <-chan map[string][dt.N_FLOORS][2]bool, // Input read from Master-Slave network module
	ordersToSlaves chan<- map[string][dt.N_FLOORS][2]bool, // Input written to Master-Slave network module
	localOrders chan<- [dt.N_FLOORS][2]bool) { // Input to local Elevator FSM

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	assignerBehaviour := dt.MS_Slave
	for {
		select {
		case assignerBehaviour = <-masterSlaveRoleChan:
			fmt.Printf("You are now %s! \n", string(assignerBehaviour))
		case input := <-ordersFromDistributor:
			switch assignerBehaviour {
			case dt.MS_Slave:
			case dt.MS_Master:
				jsonBytes, err := json.Marshal(input)
				if err != nil {
					fmt.Println("json.Marshal error: ", err)
					break
				}
				ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
				if err != nil {
					fmt.Println("exec.Command error: ", err)
					fmt.Println(string(ret))
					break
				}
				output := map[string][dt.N_FLOORS][2]bool{}
				err = json.Unmarshal(ret, &output)
				if err != nil {
					fmt.Println("json.Unmarshal error: ", err)
					break
				}
				if localHallOrders, ok := output[localIP]; ok {
					fmt.Printf("OASSIGN, deadlock 1! ")
					localOrders <- localHallOrders
					fmt.Printf("... kidding, no OASSIGN deadlock 1...\n ")
				}
				fmt.Printf("OASSIGN, deadlock 2! ")
				ordersToSlaves <- output
				fmt.Printf("... kidding, no OASSIGN deadlock 2...\n ")
			}
		case newOrders := <-ordersFromMaster:
			switch assignerBehaviour {
			case dt.MS_Master:
			case dt.MS_Slave:
				if localHallOrders, ok := newOrders[localIP]; ok {
					fmt.Printf("OASSIGN, deadlock 3! ")
					localOrders <- localHallOrders
					fmt.Printf("... kidding, no OASSIGN deadlock 3...\n ")
				}
			}
		}
	}
}
