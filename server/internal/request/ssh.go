package request

type RunCommandReq struct {
	NodeID  uint   `json:"node_id"`
	Command string `json:"command"`
}
