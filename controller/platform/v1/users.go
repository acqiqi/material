package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
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
	_, err := models.GetUsersInfoDD(data.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经绑定，请勿重复绑定")
		return
	}

}
