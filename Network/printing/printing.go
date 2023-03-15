package printing

import (
	"fmt"
	dt "project/commonDataTypes"
	"strings"
)

// States for hall requests
const (
	STATE_none      dt.RequestState = 0
	STATE_new       dt.RequestState = 1
	STATE_confirmed dt.RequestState = 2
)

func WW_toString(WW dt.AllElevDataJSON) string {
	// Create the separator row
	separatorRow := "-------------------------------------------------------------------------------------------------------\n"
	text := "####################################################################################################\n"
	text += "________________________________________WorldView________________________________________________\n\n"

	// Find the maximum length of any ID.
	ColLen := 14

	// Print the header row.
	text = text + fmt.Sprintf("%-*s | Behavior       | Floor          | Direction      | Cab Requests\n", ColLen, "ID")
	text += separatorRow

	// Print each elevator's data.
	for id, elevData := range WW {
		text = text + fmt.Sprintf("%-*s | %-*s | %-*d | %-*s | %v\n",
			ColLen, id,
			ColLen, elevData.Behavior,
			ColLen, elevData.Floor,
			ColLen, elevData.Direction,
			elevData.CabRequests)
	}
	text += separatorRow
	// text = text + fmt.Sprintf("\nHallrequest: %v \n", WW.HallRequests)
	// text += separatorRow
	text = text + "\n####################################################################################################\n"
	return text
}

func RSM_toString(RSM dt.RequestStateMatrix) string {
	text := "####################################################################################################\n"
	text = text + "______________________________________________REQ MAT________________________________________________\n\n"
	separatorRow := "-------------------------------------------------------------------------------------------------------\n"

	// Build the header
	header := []string{"ID", "|1_UP", "|1_DWN", "|2_UP", "|2_DWN", "|3_UP", "|3_DWN", "|4_UP", "|4_DWN"}
	var headerStr string
	for _, v := range header {
		headerStr += fmt.Sprintf("%-12s", v)
	}
	headerStr += "\n"
	text = text + headerStr
	text += separatorRow

	// Iterate over each elevator's data
	for id, reqData := range RSM {
		text += fmt.Sprintf("%-12s", id)
		for _, state := range reqData {
			switch state[0] {
			case STATE_none:
				text += "|NONE       "
			case STATE_new:
				text += "|NEW        "
			case STATE_confirmed:
				text += "|CONF       "
			default:
				text += "|UNDF       "
			}
			switch state[1] {
			case STATE_none:
				text += "|NONE       "
			case STATE_new:
				text += "|NEW        "
			case STATE_confirmed:
				text += "|CONF       "
			default:
				text += "|UNDF       "
			}
		}
		text += "\n"
	}
	text += separatorRow
	text += "#######################################################################################################\n"

	return text
}

func OrdersToString(role dt.MasterSlaveRole, sentOrders map[string][dt.N_FLOORS][2]bool, receivedOrders map[string][dt.N_FLOORS][2]bool) string {

	// Create the header row
	header := []string{"ID", "|1_UP", "|1_DWN", "|2_UP", "|2_DWN", "|3_UP", "|3_DWN", "|4_UP", "|4_DWN"}
	var headerStr string
	for _, v := range header {
		headerStr += fmt.Sprintf("%-12s", v)
	}
	headerStr += "\n"

	// Create the separator row
	separatorRow := "-------------------------------------------------------------------------------------------------------\n"

	// Add the header and separator rows to the text
	text := "############################################ MS_ROLE: " + strings.ToUpper(string(role)) + "  ##########################################\n\n"
	if role == dt.MS_Master {
		text += "___________________________________________ Sent orders _______________________________________________\n\n"
	} else {
		text += "___________________________________________ Recevied orders ______________________________________________\n\n"
	}
	text += headerStr
	text += separatorRow

	// Loop over each node ID
	for id := range sentOrders {
		// Truncate the ID if it is longer than the maximum width

		// Create a new row for the current node
		row := fmt.Sprintf("%-12s", id)

		// Loop over each floor for the current node
		for floor := 0; floor < dt.N_FLOORS; floor++ {
			// Get the up and down orders for the current floor
			up := sentOrders[id][floor][0]
			down := receivedOrders[id][floor][1]

			// Add the up and down orders to the current row
			row += fmt.Sprintf("|%t      |%t      ", up, down)
		}

		// Add the current row to the text
		text += row + "\n"
	}

	// Add the separator row to the text
	text += separatorRow
	text += "#######################################################################################################\n"

	// Return the final text string
	return text
}

// func PrintMessage(message NetworkMessage_t) {
// 	fmt.Printf("id: %+v\n", message.Sender_id)
// 	fmt.Printf("behaviour: %+v\n", Eb_toString(message.Behaviour))
// 	fmt.Printf("floor: %+v\n", message.Floor)
// 	fmt.Printf("direction: %+v\n", Ed_toString(message.Direction))
// 	fmt.Printf("available: %+v\n", message.Available)
// 	fmt.Printf("    Up                                         Down                                       Cab\n")
// 	for i, rq := range message.SenderHallRequests {
// 		fmt.Printf("%d - %s\n", i+1, REQ_toString(rq, message.AllCabRequests[message.Sender_id][i]))
// 	}
// 	fmt.Printf("###################################################################################################################################|\n")
// }

// func PrintPeers(p peers.PeerUpdate) {
// 	fmt.Printf("Peer update:\n")
// 	fmt.Printf("  Peers:    %q\n", p.Peers)
// 	fmt.Printf("  New:      %q\n", p.New)
// 	fmt.Printf("  Lost:     %q\n", p.Lost)
// }
