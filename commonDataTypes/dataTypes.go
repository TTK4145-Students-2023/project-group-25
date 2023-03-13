package dt

type MasterSlaveRole string

const (
	MS_Master MasterSlaveRole = "master"
	MS_Slave  MasterSlaveRole = "slave"
)

type ElevDataJSON struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type AllElevDataJSON map[string]ElevDataJSON

type AllElevDataJSON_withID struct {
	ID      string          `json:"id"`
	AllData AllElevDataJSON `json:"allData"`
}

type CostFuncInput struct {
	HallRequests [][2]bool       `json:"hallRequests"`
	States       AllElevDataJSON `json:"states"`
}

type RequestState int

type SingleNode_requestStates [][2]RequestState

type RequestStateMatrix map[string]SingleNode_requestStates

type RequestStateMatrix_with_ID struct {
	IpAdress      string             `json:"ipAdress"`
	RequestMatrix RequestStateMatrix `json:"requestMatrix"`
}
