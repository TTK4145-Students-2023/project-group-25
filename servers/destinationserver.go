package server

var server_destination = make(chan int)

func DestinationServer() {
	destination := -1
	for {
		select {
		case newDestination := <-server_destination:
			destination = newDestination

		case server_destination <- destination:
		}
	}
}

func GetDestinationFloor() int      { return <-server_destination }
func SetDestinationFloor(floor int) { server_destination <- floor }
