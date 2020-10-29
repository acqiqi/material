package platform

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"material/lib/e"
	"material/lib/gredis"
	"material/models"
)

type platform_login struct {
	PlatformKey    string `json:"platform_key"`
	PlatformSecret string `json:"platform_secret"`
}

// 获取AccessToken
func GetAccessToken(c *gin.Context) {
	data := platform_login{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	platform, err := models.PlatformGetInfoOrKey(data.PlatformKey)
	if err != nil {
		e.ApiErr(c, "平台不存在")
		return
	}

	//判断秘钥是否正确
	if platform.PlatformSecret != data.PlatformSecret {
		e.ApiErr(c, "非法请求")
		return
	}

	token := uuid.NewV4().String()
	gredis.SetCacheString("PLATFORM"+token, platform.PlatformKey, 7800000000)

	e.ApiOk(c, "登录成功", struct {
		Token string `json:"token"`
	}{
		Token: token,
	})

}
