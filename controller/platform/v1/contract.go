package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"material/service/contract_service"
)

// 合同同步
func ContractSync(c *gin.Context) {
	data := contract_service.ContractAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	platform, _ := c.Get("platform")
	project, err := models.ProjectGetInfo(data.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}
	if project.PlatformKey != platform.(models.Platform).PlatformKey {
		e.ApiErr(c, "非法请求 Err Auth PlatformKey")
		return
	}
	// 查询企业是否存在
	company, err := models.CompanyGetInfo(project.CompanyId)
	if err != nil {
		e.ApiErr(c, "请选择正确的材料商")
		return
	}
	data.CompanyId = company.Id

	_, err = models.ContractCheck(data.PlatformId, platform.(models.Platform).PlatformKey, data.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经同步过，请勿重复同步")
		return
	}

	//手动创建默认绑定
	data.PlatformKey = platform.(models.Platform).PlatformKey
	data.Cuid = company.Cuid
	data.IsPlatform = 1
	cb, err := contract_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	p, _ := models.ContractInfo(cb.Id)
	e.ApiOk(c, "创建成功", p)
}
