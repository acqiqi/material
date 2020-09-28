package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"material/service/project_service"
)

func ProjectSync(c *gin.Context) {
	data := project_service.ProjectAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	platform, _ := c.Get("platform")
	_, err := models.ProjectCheck(data.PlatformId, platform.(models.Platform).PlatformKey, data.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经同步过，请勿重复同步")
		return
	}
	// 查询企业是否存在
	company, err := models.CompanyGetInfo(data.CompanyId)
	if err != nil {
		e.ApiErr(c, "请选择正确的材料商")
		return
	}

	//手动创建默认绑定
	data.BindState = 0
	data.BindType = 1
	data.DataOrigin = platform.(models.Platform).PlatformName
	data.PlatformKey = platform.(models.Platform).PlatformKey
	data.Cuid = company.Cuid
	data.IsPlatform = 1
	cb, err := project_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	p, _ := models.ProjectGetInfo(cb.Id)

	e.ApiOk(c, "创建成功", p)
}
