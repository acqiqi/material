package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
)

// 产品详情
func ProductInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	product, err := models.ProductGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info models.Product `json:"info"`
	}{
		Info: *product,
	})
}
