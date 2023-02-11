package orderhandler

import (
	"Driver-go/elevio"
)

type InputServerChan struct {
	DRV_floors  chan int
	DRV_obstr   chan bool
	DRV_buttons chan elevio.ButtonEvent
	DRV_stop    chan bool
}

type InputServerData struct {
	DRV_floors int
	DRV_obstr  bool
	DRV_stop   bool
}
