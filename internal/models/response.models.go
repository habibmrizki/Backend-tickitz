package models

type Response struct {
	Message string
	Status  string
	Data    interface{} `json:"data,omitempty"`
}

type ResponseProfile struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseWithData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
