package btnEventSplitter

import (
	elevio "project/localElevator/elev_driver"
)

func BtnEventSplitter(btnEvent chan elevio.ButtonEvent,
	hallEvent chan elevio.ButtonEvent,
	cabEvent chan elevio.ButtonEvent) {
	for {
		select {
		case event := <-btnEvent:
			if event.Button == elevio.BT_Cab {
				cabEvent <- event
			} else {
				hallEvent <- event
			}
		}
	}
}
