package dt

type MasterSlaveRole string

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

const (
	MS_Master MasterSlaveRole = "master"
	MS_Slave  MasterSlaveRole = "slave"
)

type ElevDataJSON struct {
	Behavior    string         `json:"behaviour"`
	Floor       int            `json:"floor"`
	Direction   string         `json:"direction"`
	CabRequests [N_FLOORS]bool `json:"cabRequests"`
}

type AllElevDataJSON_withID struct {
	ID      string          `json:"id"`
	AllData AllElevDataJSON `json:"allData"`
}
type AllElevDataJSON map[string]ElevDataJSON


type CostFuncInput struct {
	HallRequests [N_FLOORS][2]bool `json:"hallRequests"`
	States       AllElevDataJSON   `json:"states"`
}

type RequestState int

type SingleNode_requestStates [N_FLOORS][2]RequestState

type RequestStateMatrix map[string]SingleNode_requestStates

type RequestStateMatrix_with_ID struct {
	IpAdress      string             `json:"ipAdress"`
	RequestMatrix RequestStateMatrix `json:"requestMatrix"`
}
