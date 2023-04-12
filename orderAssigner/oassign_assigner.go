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
	masterSlaveRoleCh <-chan dt.MasterSlaveRole,
	costFuncInputCh <-chan dt.CostFuncInputSlice,
	ordersFromMasterCh <-chan []dt.SlaveOrders,
	ordersToSlavesCh chan<- []dt.SlaveOrders,
	assignedOrdersCh chan<- [dt.N_FLOORS][2]bool) {

	var (
		hraExecutable     = ""
		assignerBehaviour = dt.MS_SLAVE

		elevHallOrders = [dt.N_FLOORS][2]bool{}
		ordersToSlaves = map[string][dt.N_FLOORS][2]bool{}

		localOrdersTimer    = time.NewTimer(time.Hour)
		ordersToSlavesTimer = time.NewTimer(time.Hour)
	)
	localOrdersTimer.Stop()
	ordersToSlavesTimer.Stop()

	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}
	for {
		select {
		case assignerBehaviour = <-masterSlaveRoleCh:
		case costFuncInput := <-costFuncInputCh:
			input := dt.SliceToCostFuncInput(costFuncInput)
			switch assignerBehaviour {
			case dt.MS_SLAVE:
			case dt.MS_MASTER:
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
				if newElevHallOrders, ok := output[localIP]; ok {
					elevHallOrders = newElevHallOrders
					localOrdersTimer.Reset(1)
				}
				ordersToSlaves = output
				ordersToSlavesTimer.Reset(1)
			}
		case ordersFromMaster := <-ordersFromMasterCh:

			switch assignerBehaviour {
			case dt.MS_MASTER:
			case dt.MS_SLAVE:

				for _, newHallOrder := range ordersFromMaster {
					if newHallOrder.IP == localIP {
						elevHallOrders = newHallOrder.Orders
						localOrdersTimer.Reset(1)
					}
				}
			}
		case <-ordersToSlavesTimer.C:
			select {
			case ordersToSlavesCh <- dt.SlaveOrdersMapToSlice(ordersToSlaves):
			default:
				ordersToSlavesTimer.Reset(1)
			}
		case <-localOrdersTimer.C:
			select {
			case assignedOrdersCh <- elevHallOrders:
			default:
				localOrdersTimer.Reset(1)
			}
		}
	}
}
