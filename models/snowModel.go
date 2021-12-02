package models

type Results struct {
	Result Snow_Response `json:"result,omitempty"`
}
type Snow_Response struct {
	SysID  string `json:"sys_id,omitempty"`
	Status string `json:"status,omitempty"`
}

type SysIDs struct {
	SysID  string `json:"sys_id,omitempty"`
	Status string `json:"status,omitempty"`
}
