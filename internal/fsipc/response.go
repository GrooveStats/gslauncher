package fsipc

type PingResponse struct {
	Payload string `json:"payload" validate:"required"`
}

type NetworkResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
