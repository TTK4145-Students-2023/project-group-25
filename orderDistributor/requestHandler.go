package requestHandler



// States
type RequestState int

const (
	Req_STATE_none      RequestState = 0
	Req_STATE_new       RequestState = 1
	Req_STATE_confirmed RequestState = 2
)

type SingleNode_RequestStates struct {
	Requests [][2]RequestState `json:"requests"`
}


