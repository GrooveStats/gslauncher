package fsipc

type PingVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type PingResponse struct {
	Version PingVersion `json:"version"`
}

type NetworkResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
