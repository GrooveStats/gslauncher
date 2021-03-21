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
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
