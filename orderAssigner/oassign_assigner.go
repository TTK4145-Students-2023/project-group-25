package oassign

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"project/Network/Utilities/bcast"
	peers "project/Network/Utilities/peers"
	dt "project/commonDataTypes"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func OrderAssigner(localIP string,
	peerUpdateCh <-chan peers.PeerUpdate,
	costFuncInputCh <-chan dt.CostFuncInputSlice,
	assignedOrdersCh chan<- [dt.N_FLOORS][2]bool) {

	var (
		receiveOrdersCh   = make(chan map[string][dt.N_FLOORS][2]bool)
		transmittOrdersCh = make(chan map[string][dt.N_FLOORS][2]bool)

		hraExecutable     = ""
		assignerBehaviour = dt.MASTER

		localHallOrders       = [dt.N_FLOORS][2]bool{}
		ordersToExternalNodes = map[string][dt.N_FLOORS][2]bool{}
		//ordersFromMaster      = [dt.N_FLOORS][2]bool{}

		broadCastTimer   = time.NewTimer(1)
		localOrdersTimer = time.NewTimer(time.Hour)
	)
	broadCastTimer.Stop()
	localOrdersTimer.Stop()

	go bcast.Receiver(dt.MS_PORT, receiveOrdersCh)
	go bcast.Transmitter(dt.MS_PORT, transmittOrdersCh)

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
		case peerUpdate := <-peerUpdateCh:
			fmt.Printf("p update in O.Ass: %+v\n", peerUpdate)
			assignerBehaviour = assignRole(localIP, peerUpdate.Peers)
			fmt.Println(assignerBehaviour)
		case costFuncInput := <-costFuncInputCh:
			fmt.Println("Received CostFuncInput")
			input := dt.CostFuncInputSliceToMap(costFuncInput)
			switch assignerBehaviour {
			case dt.SLAVE:
			case dt.MASTER:
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
					localHallOrders = newElevHallOrders
					localOrdersTimer.Reset(1)
				}
				ordersToExternalNodes = output
				broadCastTimer.Reset(1)
			}
		case newOrders := <-receiveOrdersCh:

			if !reflect.DeepEqual(newOrders[localIP], localHallOrders) {
				fmt.Println("New order from external Master")
				fmt.Printf("new: %+v\n stored: %+v\n", newOrders[localIP], localHallOrders)
				switch assignerBehaviour {
				case dt.MASTER:
				case dt.SLAVE:
					localHallOrders = newOrders[localIP]
				}
			}
		case <-broadCastTimer.C:

			broadCastTimer.Reset(dt.BROADCAST_PERIOD)
			switch assignerBehaviour {
			case dt.SLAVE:
			case dt.MASTER:

				transmittOrdersCh <- ordersToExternalNodes
				fmt.Println("Orders Transmitted")
				fmt.Println(ordersToExternalNodes)
			}
		case <-localOrdersTimer.C:
			select {
			case assignedOrdersCh <- localHallOrders:
			default:
				localOrdersTimer.Reset(1)
			}
		}
	}
}

func assignRole(localIP string, peers []string) dt.AssignerBehaviour {
	if len(peers) == 0 {
		return dt.MASTER
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
		return dt.MASTER
	}
	return dt.SLAVE
}
