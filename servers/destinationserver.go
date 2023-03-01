package server

var (
	setDestinationChan     = make(chan int)
	getDestinationChan     = make(chan int)
	destinationChangedChan = make(chan bool)
)

func DestinationServer() {
	destination := -1
	destinationChanged := false
	for {
		select {
		case destination = <-setDestinationChan:
			destinationChanged = true

		case destinationChanged = <-destinationChangedChan:
		case destinationChangedChan <- destinationChanged:
		case getDestinationChan <- destination:
		}
	}
}

func SetDestinationFloor(floor int) { setDestinationChan <- floor }

func GetDestinationFloor() int     { return <-getDestinationChan }
func DestinationChangeIsRecieved() { destinationChangedChan <- false }
func DestinationHasChanged() bool  { return <-destinationChangedChan }
