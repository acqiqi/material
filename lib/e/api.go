package e

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 提供空 obj类型
func GetEmptyStruct() interface{} {
	return struct {
	}{}
}

type ApiJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//返回成功
func ApiOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, ApiJson{
		Code: SUCCESS,
		Msg:  msg,
		Data: data,
	})
	c.Abort()
	return
}

// 返回错误
func ApiErr(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, ApiJson{
		Code: ERROR,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
	return
}

// 返回其他数据类型
func ApiOpt(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, ApiJson{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	c.Abort()
	return
}
