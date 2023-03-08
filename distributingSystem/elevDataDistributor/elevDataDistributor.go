package elevDataDistributor

var localID string = "ID1"

// Datatypes
type ElevData struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type AllElevData map[string]ElevData

type AllElevData_withID struct {
	ID      string      `json:"id"`
	AllData AllElevData `json:"allData"`
}

type WorldView struct {
	HallRequests [][2]bool   `json:"hallRequests"`
	AllData      AllElevData `json:"elevStates"`
}

// input channels
var (
	allElevData_fromP2P = make(chan AllElevData_withID)
	localElevData       = make(chan ElevData)
	HallOrderArray      = make(chan [][2]bool)
)

// output channels
var (
	allElevData_toP2P   = make(chan AllElevData)
	WorlView_toAssigner = make(chan WorldView)
)
