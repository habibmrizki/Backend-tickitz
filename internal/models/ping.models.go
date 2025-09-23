package models

type Body struct {
	Id      int    `json:"id" binding:"required,min=0"`
	Message string `json:"msg"`
	Gender  string `json:"gender"`
}

type PingModel struct {
	Message     string `example:"pong"`
	RequestId   string `example:"123"`
	ContentType string `example:"application/json"`
}
