package server

type InputData struct {
	Floors int
	Obstr  bool
	Stop   bool
}

var (
	SetCurrentFloorChan = make(chan int)
	SetObstrValChan     = make(chan bool)
	SetStopValChan      = make(chan bool)

	getCurrentFloorChan = make(chan int)
	getObstrValChan     = make(chan bool)
	getStopValChan      = make(chan bool)
)

func InputServer() {
	inputData := InputData{
		Floors: -1,
		Obstr:  false,
		Stop:   false}
	for {
		select {
		case inputData.Floors = <-SetCurrentFloorChan:
		case inputData.Obstr = <-SetObstrValChan:
		case inputData.Stop = <-SetStopValChan:

		case getCurrentFloorChan <- inputData.Floors:
		case getObstrValChan <- inputData.Obstr:
		case getStopValChan <- inputData.Stop:
		}
	}
}

func GetCurrentFloor() int { return <-getCurrentFloorChan }
func GetObstrVal() bool    { return <-getObstrValChan }
func GetStopVal() bool     { return <-getStopValChan }
