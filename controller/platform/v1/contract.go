package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"material/service/contract_service"
	"material/service/receiver_service"
)

// 合同同步
func ContractSync(c *gin.Context) {
	data := struct {
		Info   contract_service.ContractAdd       `json:"info"`
		Config contract_service.ContractConfigAdd `json:"config"`
		//Users  []receiver_service.UserAdd         `json:"users"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	platform, _ := c.Get("platform")
	project, err := models.ProjectGetInfo(data.Info.ProjectId)
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
	data.Info.CompanyId = company.Id

	_, err = models.ContractCheck(data.Info.PlatformId, platform.(models.Platform).PlatformKey, data.Info.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经同步过，请勿重复同步")
		return
	}

	//手动创建默认绑定
	data.Info.PlatformKey = platform.(models.Platform).PlatformKey
	data.Info.Cuid = company.Cuid
	data.Info.IsPlatform = 1
	cb, err := contract_service.Add(&data.Info)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//同步付款信息
	data.Config.ContractId = cb.Id
	conf_cb, err := contract_service.AddConfig(&data.Config)
	p, _ := models.ContractInfo(cb.Id)

	////检查和绑定用户
	//user_cb, err := receiver_service.SyncUsers(data.Users, p, platform.(models.Platform).PlatformKey)
	//if err != nil {
	//	e.ApiErr(c, err.Error())
	//	return
	//}

	e.ApiOk(c, "创建成功", struct {
		Info   models.Contract                          `json:"info"`
		Config models.ContractConfig                    `json:"config"`
		Users  []receiver_service.ReceiverUsersCallback `json:"users"`
	}{
		Info:   *p,
		Config: *conf_cb,
		//Users:  user_cb,
	})
}
