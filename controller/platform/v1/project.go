package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/project_service"
	"material/service/receiver_service"
	"strconv"
)

// 同步项目
func ProjectSync(c *gin.Context) {
	data := struct {
		Project project_service.ProjectAdd `json:"project"`
		Users   []receiver_service.UserAdd `json:"users"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	d_str := utils.JsonEncode(data)
	log.Println(d_str)

	platform, _ := c.Get("platform")
	_, err := models.ProjectCheck(data.Project.PlatformId, platform.(models.Platform).PlatformKey, data.Project.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经同步过，请勿重复同步")
		return
	}
	// 查询企业是否存在
	company, err := models.CompanyGetInfo(data.Project.CompanyId)
	if err != nil {
		e.ApiErr(c, "请选择正确的材料商")
		return
	}

	//手动创建默认绑定
	data.Project.BindState = 1
	data.Project.BindType = 1
	data.Project.DataOrigin = platform.(models.Platform).PlatformName
	data.Project.PlatformKey = platform.(models.Platform).PlatformKey
	data.Project.Cuid = company.Cuid
	data.Project.IsPlatform = 1

	cb, err := project_service.Add(&data.Project)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	p, _ := models.ProjectGetInfo(cb.Id)
	//检查和绑定用户
	user_cb, err := receiver_service.SyncUsers(data.Users, p, platform.(models.Platform).PlatformKey)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//查询Platform
	if platform, err := models.PlatformGetInfoOrKey(p.PlatformKey); err == nil {
		callback := e.HttpCallbackData{
			Code:        0,
			Msg:         "Receive Success",
			Action:      e.PLATFORM_ACTION_PROJECT_RECEIVE, //接收
			CallbackUrl: platform.MessageCallbackUrl,
			Data: e.PlatformProjectReceiveCallback{
				Id:          strconv.FormatInt(p.Id, 10),
				ProjectName: p.ProjectName,
				State:       p.State,
				CompanyId:   p.CompanyId,
				Company:     p.Company,
				BindState:   p.BindState,
				PlatformKey: p.PlatformKey,
				PlatformUid: p.PlatformUid,
				PlatformId:  p.PlatformId,
				CreatedAt:   p.CreatedAt,
				Status:      1,
				SupplierId:  p.SupplierId,
			},
		}
		c_data := new(e.HttpCallbackData)

		c_str := utils.JsonEncode(callback)
		log.Println(c_str)

		err = callback.RequestCallback(&c_data)
		if err != nil {
			log.Println(err.Error())
		}
	}
	e.ApiOk(c, "创建成功", struct {
		Info  interface{}                              `json:"info"`
		Users []receiver_service.ReceiverUsersCallback `json:"users"`
	}{
		Info:  *p,
		Users: user_cb,
	})
}
