package peers

import (
	"fmt"
	"net"
	"project/Network/Utilities/conn"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func PeerListHandler(localIP string,
	peerUpdate_MS chan<- PeerUpdate,
	peerUpdate_DataDistributor chan<- PeerUpdate,
	peerUpdate_OrderHandler chan<- PeerUpdate,
) {
	peerUpdateCh := make(chan PeerUpdate) // channel for receiving updates on the id of the peers that are alive on the network
	peerTxEnable := make(chan bool)       // disable/enable the transmitter after started

	go Transmitter(15669, localIP, peerTxEnable)
	go Receiver(15669, peerUpdateCh)

	for {
		peerList := <-peerUpdateCh
		fmt.Printf("Do we not get peerUpdate on channel?\n\n")
		peerListSendt := [3]bool{false, false, false}
		for peerListSendt != [3]bool{true, true, true} {
			if !peerListSendt[0] {
				select {
				case peerUpdate_MS <- peerList:
					peerListSendt[0] = true
				case peerList = <-peerUpdateCh:
					peerListSendt = [3]bool{false, false, false}
				default:
				}
			}
			if !peerListSendt[1] {
				select {
				case peerUpdate_DataDistributor <- peerList:
					peerListSendt[1] = true
				case peerList = <-peerUpdateCh:
					peerListSendt = [3]bool{false, false, false}
				default:
				}
			}
			if !peerListSendt[2] {
				fmt.Printf("in here\n\n")
				select {
				case peerUpdate_OrderHandler <- peerList:
					peerListSendt[2] = true
					fmt.Printf("not in here\n\n")
				case peerList = <-peerUpdateCh:
					peerListSendt = [3]bool{false, false, false}
				default:
				}
			}
		}
	}
}

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}
