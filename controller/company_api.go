package controller

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"material/lib/e"
	"material/lib/gredis"
	"material/models"
)

type platform_api struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

// 获取AccessToken
func GetAccessTokenApi(c *gin.Context) {
	data := platform_api{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	if data.Ak == "" {
		e.ApiErr(c, "参数有误")
		return
	}

	company, err := models.CompanyGetInfoOrAk(data.Ak)
	if err != nil {
		e.ApiErr(c, "平台不存在")
		return
	}

	//判断秘钥是否正确
	if company.Sk != data.Sk {
		e.ApiErr(c, "非法请求")
		return
	}

	token := uuid.NewV4().String()
	gredis.SetCacheString("COMPANY_API"+token, company.Ak, 7800000000)

	e.ApiOk(c, "登录成功", struct {
		Token string `json:"token"`
	}{
		Token: token,
	})

}
