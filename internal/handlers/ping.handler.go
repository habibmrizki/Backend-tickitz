package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

type PingWithParam struct {
	Param  string `json:"param"`
	Param2 string `json:"param2"`
	Q      string `json:"q"`
}

// GetPing
// @tags		ping
// @router	/ping [GET]
// @Param		X-Request-ID	header	string	true	"Header for requestID"
// @Param		Content-Type	header	string	true	"Header for requestID"
// @produce	json
// @success	200 {object} models.PingModel
func (p *PingHandler) GetPing(ctx *gin.Context) {
	requestId := ctx.GetHeader("X-Request-ID")
	contentType := ctx.GetHeader("Content-Type")
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message":     "pong",
	// 	"requestId":   requestId,
	// 	"contentType": contentType,
	// })

	ctx.JSON(http.StatusOK, models.PingModel{
		Message:     "pong",
		RequestId:   requestId,
		ContentType: contentType,
	})
}

// GetPingWithParam
// @tags		ping
// @router	/ping/:id/:param2 [GET]
// @Param		id		path	integer	true	"path params for id"
// @Param		param2	path	string	true	"path params for param2"
// @Param 		q 		query	string 	false	"query for q"
// @produce	json
// @success	200
func (p *PingHandler) GetPingWithParam(ctx *gin.Context) {
	pingId := ctx.Param("id")
	param2 := ctx.Param("param2")
	q := ctx.Query("q")
	log.Printf("%s, %s, %s\n", pingId, param2, q)
	ctx.JSON(http.StatusOK, gin.H{
		"param":  pingId,
		"param2": param2,
		"q":      q,
	})
}

// func (p *PingHandler) PostPing(ctx *gin.Context) {
// 	body := models.Body{}
// 	// data-binding, memasukkan body ke dalam variabel golang
// 	if err := ctx.ShouldBind(&body); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   err.Error(),
// 			"success": false,
// 		})
// 		return
// 	}
// 	if err := utils.ValidateBody(body); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	log.Println(body)
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"body":    body,
// 	})
// }
