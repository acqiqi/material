package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"material/service/product_service"
)

// 获取大分类列表
func ProductCatsList(c *gin.Context) {
	list, err := product_service.CatsSelect()
	if err != nil {
		e.ApiErr(c, "请求失败 "+err.Error())
		return
	}
	e.ApiOk(c, "获取成功", struct {
		Lists interface{} `json:"list"`
	}{
		Lists: list,
	})
}

// 同步下料单
func MaterialSync(c *gin.Context) {
	data := struct {
		MData product_service.MaterialAdd  `json:"m_data"`
		PData []product_service.ProductAdd `json:"p_data"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	platform, _ := c.Get("platform")

	cb, err := product_service.ProductSync(&data.MData, data.PData, platform.(models.Platform))
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "操作成功", cb)
}
