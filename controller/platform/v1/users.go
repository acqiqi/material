package v1

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/setting"
	"material/lib/utils"
	"material/models"
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

	_, err := models.GetUsersInfoDD(data.PlatformUid)
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
