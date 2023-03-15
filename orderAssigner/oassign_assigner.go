package oassign

import (
	"encoding/json"
	"fmt"
	"os/exec"
	dt "project/commonDataTypes"
	"runtime"
	"time"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func OrderAssigner(localIP string,
	masterSlaveRoleChan <-chan dt.MasterSlaveRole,
	ordersFromDistributor <-chan dt.CostFuncInput, // Input from order distributor
	ordersFromMaster <-chan map[string][dt.N_FLOORS][2]bool, // Input read from Master-Slave network module
	ordersToSlavesChan chan<- map[string][dt.N_FLOORS][2]bool, // Input written to Master-Slave network module
	localOrders chan<- [dt.N_FLOORS][2]bool) { // Input to local Elevator FSM

	localHallOrders := [dt.N_FLOORS][2]bool{}
	ordersToSlaves := map[string][dt.N_FLOORS][2]bool{}
	localOrdersTimer := time.NewTimer(1)
	localOrdersTimer.Stop()
	ordersToSlavesTimer := time.NewTimer(1)
	ordersToSlavesTimer.Stop()

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
				if newLocalHallOrders, ok := output[localIP]; ok {
					localHallOrders = newLocalHallOrders
					localOrdersTimer.Reset(1)
				}
				ordersToSlaves = output
				ordersToSlavesTimer.Reset(1)
			}
		case newOrders := <-ordersFromMaster:
			switch assignerBehaviour {
			case dt.MS_Master:
			case dt.MS_Slave:
				if newLocalHallOrders, ok := newOrders[localIP]; ok {
					localHallOrders = newLocalHallOrders
					localOrdersTimer.Reset(1)
				}
			}
		case <-ordersToSlavesTimer.C:
			select {
			case ordersToSlavesChan <- ordersToSlaves:
			default:
				ordersToSlavesTimer.Reset(1)
			}
		case <-localOrdersTimer.C:
			select {
			case localOrders <- localHallOrders:
			default:
				localOrdersTimer.Reset(1)
			}
		}
	}
}
