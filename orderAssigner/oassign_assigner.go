package oassign

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"project/Network/Utilities/localip"
	dt "project/commonDataTypes"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func OrderAssigner(masterSlaveRoleChan <-chan dt.MasterSlaveRole,
	ordersFromDistributor <-chan dt.CostFuncInput, // Input from order distributor
	ordersFromMaster <-chan []byte, // Input read from Master-Slave network module
	ordersToSlaves chan<- []byte, // Input written to Master-Slave network module
	localOrders chan<- [][2]bool) { // Input to local Elevator FSM

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	localIpAdress, _ := localip.LocalIP()

	assignerBehaviour := dt.MS_Slave
	for {
		select {
		case assignerBehaviour = <-masterSlaveRoleChan:
			fmt.Printf("We are now the %s\n", string(assignerBehaviour))
		case input := <-ordersFromDistributor:
			fmt.Printf("We have recieved data from Distributor\n")
			switch assignerBehaviour {
			case dt.MS_Slave:
			case dt.MS_Master:
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
				output := map[string][][2]bool{}
				err = json.Unmarshal(ret, &output)
				if err != nil {
					fmt.Println("json.Unmarshal error: ", err)
					return
				}
				if localHallOrders, ok := output[localIpAdress]; ok {
					localOrders <- localHallOrders
				}
				fmt.Printf("Sending and ... ")
				ordersToSlaves <- ret
				fmt.Printf("...sendt! Data to slave\n\n")
			}
		case input := <-ordersFromMaster:
			switch assignerBehaviour {
			case dt.MS_Master:
			case dt.MS_Slave:
				output := map[string][][2]bool{}
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
