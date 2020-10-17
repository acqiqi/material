package e

import (
	"github.com/gin-gonic/gin"
	"material/lib/utils"
	"net/http"
)

// 提供空 obj类型
func GetEmptyStruct() interface{} {
	return struct {
	}{}
}

type ApiPageLists struct {
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	Total   int                    `json:"total_size"`
	Lists   interface{}            `json:"lists"`
	Map     []utils.WhereData      `json:"map"`
	OptParm map[string]interface{} `json:"opt_parm"`
}

type ApiId struct {
	Id int64 `json:"id"`
}

//检测默认值
func CheckApiPageListDefault(cb *ApiPageLists) {
	if cb.Page <= 0 {
		cb.Page = 1
	}
	if cb.Limit <= 0 {
		cb.Limit = 20
	}
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

// 用于Callback 请求数据集合
type HttpCallbackData struct {
	Code   int         `json:"code"`   // 状态码 默认0成功 1失败 其他专用错误码
	Msg    string      `json:"msg"`    // 主体消息
	Action string      `json:"action"` // 行为 利用行为二次解析Data结构体
	Data   interface{} `json:"data"`   // 主体数据根据Action反序列
}
