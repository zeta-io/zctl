package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zeta-io/ginx"
	"github.com/zeta-io/zeta"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func suc(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

func fail(c *gin.Context, data interface{}) {
	c.JSON(500, data)
}


func list(c *gin.Context, input interface{}) (string, error){
	return "hello zeta", nil
}

func main() {
	router := zeta.Router("/api/:version/users")
	router.Get("", list)

	e := zeta.New(router, ginx.New(gin.New())).Run(":8080")
	if e != nil{
		panic(e)
	}
}