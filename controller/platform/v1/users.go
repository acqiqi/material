package v1

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/setting"
	"material/lib/utils"
	"material/models"
	"material/service/receiver_service"
)

// 去小程序绑定的Qrcode 是没有到任何子企业的
func UsersDDSyncQrcode(c *gin.Context) {
	data := struct {
		PlatformUid string `json:"platform_uid"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	if data.PlatformUid == "" {
		e.ApiErr(c, "请输入正确的用户id")
		return
	}

	platform, _ := c.Get("platform")
	_, err := models.PlatformUsersCheckUser(data.PlatformUid, platform.(models.Platform).PlatformKey)
	if err == nil {
		e.ApiErr(c, "已经绑定，请勿重复绑定")
		return
	}

	wechat_utils := utils.WechatUtils{}
	wechat_utils.Init(setting.WechatSetting.SmallAppID, setting.WechatSetting.AppSecret)
	wechat_utils.SmallQrcodeData.Width = 430
	wechat_utils.SmallQrcodeData.Path = "/page/index/index"
	if err := wechat_utils.GetAccessToken(); err != nil {
		e.ApiErr(c, "服务器异常")
		return
	}
	if cb, err := wechat_utils.GetSmallQrcode(); err != nil {
		e.ApiErr(c, "获取失败")
		return
	} else {
		b64 := base64.StdEncoding.EncodeToString(cb)
		e.ApiOk(c, "获取成功", struct {
			Qrcode string `json:"qrcode"`
		}{
			Qrcode: b64,
		})
	}
}

func UsersSync(c *gin.Context) {

	data := struct {
		//CompanyId  int64                      `json:"company_id"`
		ContractId int64                      `json:"contract_id"`
		Users      []receiver_service.UserAdd `json:"users"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	contract, err := models.ContractInfo(data.ContractId)
	if err != nil {
		e.ApiErr(c, "合同不存在")
		return
	}

	platform, _ := c.Get("platform")
	if contract.PlatformKey != platform.(models.Platform).PlatformKey {
		e.ApiErr(c, "非法请求")
		return
	}

	//检查和绑定用户
	cb, err := receiver_service.SyncUsers(data.Users, contract, platform.(models.Platform).PlatformKey)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "同步成功", cb)
}
