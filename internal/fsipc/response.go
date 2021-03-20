package fsipc

type PingResponse struct {
}

type NetworkResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
