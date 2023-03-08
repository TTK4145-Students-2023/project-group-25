package elevtest

import (
	"fmt"
	btnassign "project/buttonAssigner"
	elevio "project/localElevator/elev_driver"
	elevfsm "project/localElevator/elev_fsm"
	oassign "project/orderAssigner"
	"time"
)

func elevatorDatatoHRAElevState(e elevfsm.ElevatorData) oassign.HRAElevState {
	var (
		behaviour   string
		floor       int
		direction   string
		cabRequests [elevfsm.N_FLOORS]bool
	)
	floor = e.Floor
	cabRequests = e.CabRequests

	switch e.Behaviour {
	case elevfsm.EB_DoorOpen:
		behaviour = "doorOpen"
	case elevfsm.EB_Idle:
		behaviour = "idle"
	case elevfsm.EB_Moving:
		behaviour = "moving"
	}

	switch e.Dirn {
	case elevio.MD_Down:
		direction = "down"
	case elevio.MD_Up:
		direction = "up"
	case elevio.MD_Stop:
		direction = "stop"
	}
	return oassign.HRAElevState{Floor: floor, Behavior: behaviour, Direction: direction, CabRequests: cabRequests}
}

func intermediateOrderDistributor(hallEvent chan elevio.ButtonEvent,
	handler_hallOrdersExecuted chan []elevio.ButtonEvent,
	elev_data chan elevfsm.ElevatorData,
	ordersFromDistributor chan oassign.HRAInput) {

	orderOverview := oassign.HRAInput{
		HallRequests: [elevfsm.N_FLOORS][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]oassign.HRAElevState{
			"127.0.0.1": {
				Behavior:    "idle",
				Floor:       2,
				Direction:   "stop",
				CabRequests: [elevfsm.N_FLOORS]bool{false, false, false, false},
			},
		},
	}

	for {
		select {
		case buttonEvent := <-hallEvent:
			fmt.Printf("HallEvent.\n")
			orderOverview.HallRequests[buttonEvent.Floor][buttonEvent.Button] = true
			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)

		case hallOrdersExecuted := <-handler_hallOrdersExecuted:
			fmt.Printf("OrderExecute.\n")
			for i := 0; i < len(hallOrdersExecuted); i++ {
				if hallOrdersExecuted[i].Button != elevio.BT_Cab {
					orderOverview.HallRequests[hallOrdersExecuted[i].Floor][hallOrdersExecuted[i].Button] = false
					elevio.SetButtonLamp(hallOrdersExecuted[i].Button, hallOrdersExecuted[i].Floor, false)
				}
			}

		case elevatorData := <-elev_data:
			fmt.Printf("newElevData.\n")
			orderOverview.States["127.0.0.1"] = elevatorDatatoHRAElevState(elevatorData)
		}
		ordersFromDistributor <- orderOverview
		fmt.Printf("Data from distro sent!\n")
	}
}

func testSpecDistributor(OrderAssignerBehaviourChan chan oassign.OrderAssignerBehaviour,
	localIpAdressChan chan string) {
	orderAssignerBehaviour := oassign.MS_Master
	localIpAdress := "127.0.0.1"

	for {
		select {
		case localIpAdressChan <- localIpAdress:
		case OrderAssignerBehaviourChan <- orderAssignerBehaviour:
		}
	}
}

var (
	OrderAssignerBehaviourChan = make(chan oassign.OrderAssignerBehaviour)
	localIpAdressChan          = make(chan string) // Chanel where local IP-adress is fetched

	ordersFromDistributor      = make(chan oassign.HRAInput) // Input from order distributor
	ordersFromMaster           = make(chan []byte)           // Input read from Master-Slave network module
	ordersToSlaves             = make(chan []byte)           // Input written to Master-Slave network module
	ordersLocal                = make(chan [elevfsm.N_FLOORS][2]bool)
	handler_hallOrdersExecuted = make(chan []elevio.ButtonEvent)

	btnEvent  = make(chan elevio.ButtonEvent)
	hallEvent = make(chan elevio.ButtonEvent)
	cabEvent  = make(chan elevio.ButtonEvent)

	drv_floors = make(chan int)
	drv_obstr  = make(chan bool)
	elev_data  = make(chan elevfsm.ElevatorData)
)

func RunSingleElevTest() {
	elevio.Init("localhost:15657", elevfsm.N_FLOORS)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollButtons(btnEvent)
	go elevio.PollObstructionSwitch(drv_obstr)

	go btnassign.ButtonHandler(btnEvent, hallEvent, cabEvent)

	go intermediateOrderDistributor(hallEvent,
		handler_hallOrdersExecuted,
		elev_data,
		ordersFromDistributor)

	go testSpecDistributor(OrderAssignerBehaviourChan,
		localIpAdressChan)
	go oassign.OrderAssigner(OrderAssignerBehaviourChan,
		localIpAdressChan,
		ordersFromDistributor,
		ordersFromMaster,
		ordersToSlaves,
		ordersLocal)

	time.Sleep(time.Millisecond * 40)

	go elevfsm.FSM(ordersLocal,
		cabEvent,
		drv_floors,
		drv_obstr,
		elev_data,
		handler_hallOrdersExecuted)

	for {
		// Kill orders that are yet to be handeled!
		<-ordersToSlaves
	}
}
