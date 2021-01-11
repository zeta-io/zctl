package api

import(
	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int 	`json:"code"`
	Msg  string `json:"msg"`
}

func suc(c *gin.Context, data interface{}){
	c.JSON(200, data)
}

func fail(c *gin.Context, data interface{}){
	c.JSON(500, data)
}