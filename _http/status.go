package _http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, msg string, data ...interface{}) {
	var res any
	if len(data) > 0 {
		res = data[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
		Data: res,
	})
}

func Failed(c *gin.Context, msg string, code ...int) {
	status := 1
	if len(code) > 0 {
		status = code[1]
	}
	c.JSON(http.StatusOK, Response{
		Code: status,
		Msg:  msg,
	})
}
