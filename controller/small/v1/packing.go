package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
)

func PackingInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	packing, err := models.PackingGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "包装不存在")
		return
	}

	//查询所有商品
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["packing_id"] = packing.Id
	products, err := models.PackingProductGetLists(0, 9999, utils.BuildWhere(maps))
	if err != nil {
		e.ApiErr(c, "导出列表失败")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Packing  models.Packing           `json:"packing"`
		Products []*models.PackingProduct `json:"products"`
	}{
		Packing:  *packing,
		Products: products,
	})
}

func PackingProductInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	pp, err := models.PackingProductGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	//检测当前用户是否有权限
	authFlag := false
	user_info, _ := c.Get("user_info")
	_, err = models.ReceiverUsersCheckAuth(user_info.(models.Users).Cuid, pp.ProjectId)
	if err == nil {
		authFlag = true
	}

	isRet := false
	if pp.Packing.SendId > 0 {
		send, err := models.SendGetInfo(pp.Packing.SendId)
		if err == nil {
			if send.Status == 0 {
				isRet = true
			}
		}
	}

	project, err := models.ProjectGetInfo(pp.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目有误")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info    interface{} `json:"info"`
		IsAuth  bool        `json:"is_auth"`
		Project interface{} `json:"project"`
		IsRet   bool        `json:"is_ret"`
	}{
		Info:    pp,
		IsAuth:  authFlag,
		Project: project,
		IsRet:   isRet,
	})
}

// 打包产品退货详情
func PackingProductReturnInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	pp, err := models.PackingProductGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	isRet := false
	if pp.Packing.SendId > 0 {
		send, err := models.SendGetInfo(pp.Packing.SendId)
		if err == nil {
			if send.Status == 0 {
				isRet = true
			}
		}
	}

	//检测当前用户是否有权限
	user_info, _ := c.Get("user_info")
	_, err = models.ReceiverUsersCheckAuth(user_info.(models.Users).Cuid, pp.ProjectId)
	if err != nil {
		e.ApiErr(c, "非法请求")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info    interface{} `json:"info"`
		OkCount float64     `json:"ok_count"` //可退货数量
		IsRet   bool        `json:"is_ret"`
	}{
		Info:    pp,
		OkCount: pp.Count - pp.ReturnCount,
		IsRet:   isRet,
	})
}
