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

func OrderAssigner(OrderAssignerBehaviourChan <-chan dt.OrderAssignerBehaviour,
	localIpAdressChan <-chan string, // Chanel where local IP-adress is fetched
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

	localIpAdress := ""
	assignerBehaviour := dt.OA_Slave
	for {
		select {
		case localIpAdress = <-localIpAdressChan:
		case assignerBehaviour = <-OrderAssignerBehaviourChan:
		case input := <-ordersFromDistributor:
			switch assignerBehaviour {
			case dt.OA_Slave:
			case dt.OA_Master:
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
				ordersToSlaves <- ret
			}
		case input := <-ordersFromMaster:
			switch assignerBehaviour {
			case dt.OA_Master:
			case dt.OA_Slave:
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
