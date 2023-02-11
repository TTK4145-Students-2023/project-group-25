package server

type InputData struct {
	Floors int
	Obstr  bool
	Stop   bool
}

var (
	server_floors = make(chan int)
	server_obstr  = make(chan bool)
	server_stop   = make(chan bool)
)

func InputServer() {
	inputData := InputData{
		Floors: -1,
		Obstr:  false,
		Stop:   false}
	for {
		select {
		case newFloor := <-server_floors:
			{
				inputData.Floors = newFloor
			}
		case newObstructionVal := <-server_obstr:
			{
				inputData.Obstr = newObstructionVal
			}
		case newStopVal := <-server_stop:
			{
				inputData.Stop = newStopVal
			}
		case server_floors <- inputData.Floors:
		case server_obstr <- inputData.Obstr:
		case server_stop <- inputData.Stop:
		}
	}

}

func GetCurrentFloor() int { return <-server_floors }
func GetObstrVal() bool    { return <-server_obstr }
func GetStopVal() bool     { return <-server_stop }

func SetCurrentFloor(floor int) { server_floors <- floor }
func SetObstrVal(obstr bool)    { server_obstr <- obstr }
func SetStopVal(stop bool)      { server_stop <- stop }
