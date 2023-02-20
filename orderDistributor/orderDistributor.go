import (
	"elvio/elevator_io.go"
)

latest_ordermatrix := [] int 

type DistributorState int

const (
	STATE_updateLocalData           DistributorState = 0
	STATE_distributeLocalChanges              		 = 1
)


//input channels 
var(
allElevData_fromP2P = make(chan ElevData)
btnPress = make(chan ButtonEvent)
orderExecuted = make(chan bool)
localElevData = make(chan ElevData)
)
//output channels 
allElevData_toP2P := make(chan ElevData)
allElevData_toAssigner := make(chan ElevData)


func dataDistributor(
	allElevData_fromP2P 		<-chan ElevData,
	btnPress 					<-chan ButtonEvent,
	orderExecuted 				<-chan bool,
	localElevData 				<-chan ElevData,         //not used as fsm trigger
	allElevData_toP2P 			chan<- ElevData,
	allElevData_toAssigner 		chan<- ElevData,

){

	for {
		distributor_state = STATE_updateData
		select {


		case allElevData := <-allElevData_fromP2P:
			switch distributor_state {
				case STATE_updateData: 
					//update local storage of data and send json to assigner

					latestValid_ordermatrix, _  = validateData(allElevData)
					allElevData_toAssigner <- latestValid_ordermatrix.json()
					
				case STATE_distributeLocalChange: 
					// check if input==output
					// if true: switch state and turn on/off light 
					//if false: send change again 

					_, succesfully_distributed := validateData(allElevData)

					if succesfully_distributed{
						distributor_state = STATE_updateLocalData
						toggle light
					}

		
							
			}

		case executedOrder := <- orderExecuted:
			switch distributor_state {
				case STATE_updateData: 
					//add order change to elevdata and send to p2p
					//switch state
					
				case STATE_distributeLocalChange: 
					// dont accept orderchanges when distrubiting 

		
							
			}
			
		case pressedBtn := <- ButtonEvent:
			switch dist_state {
				case STATE_updateData: 
					//add order change to elevdata and send to p2p
					//switch state
					
				case STATE_distributeLocalChange: 
					// dont accept orderchanges when distrubiting 
									
			}
		}
	}
}
		




