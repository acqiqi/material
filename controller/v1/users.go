package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/dd"
	"material/lib/e"
	"material/lib/setting"
	"material/lib/utils"
	"material/models"
)

// 获取用户信息
func UsersGetInfo(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	company_list, _ := models.CompanyUsersGetMyList(user_info.(models.Users).Cuid)
	e.ApiOk(c, "获取成功", struct {
		UserInfo    models.Users          `json:"user_info"`
		CompanyList []models.CompanyUsers `json:"company_list"`
	}{
		UserInfo:    user_info.(models.Users),
		CompanyList: company_list,
	})
}

// 获取修改密码手机号验证码
func GetRPSMS(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	//处理请求
	headers := make(map[string]string)
	headers["PlatformKey"] = setting.PlatformSetting.PlatformKey

	// 注册接口
	url := setting.PlatformSetting.ApiBaseUrl + dd.DD_API_AUTH_RE_SMS

	cb := e.ApiJson{}
	if err := utils.HttpPostJsonHeader(url, struct {
		Mobile string `json:"mobile"`
	}{Mobile: user_info.(models.Users).Mobile}, headers, &cb); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOpt(c, cb.Code, cb.Msg, cb.Data)
}

// 修改密码
func RePassword(c *gin.Context) {
	data := struct {
		Mobile     string `json:"mobile"`
		Code       string `json:"code"`
		Password   string `json:"password"`
		RePassword string `json:"repassword"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	user_info, _ := c.Get("user_info")
	data.Mobile = user_info.(models.Users).Mobile

	headers := make(map[string]string)
	headers["PlatformKey"] = setting.PlatformSetting.PlatformKey

	// 注册接口
	url := setting.PlatformSetting.ApiBaseUrl + dd.DD_API_AUTH_RE_PASSWORD

	cb := e.ApiJson{}
	if err := utils.HttpPostJsonHeader(url, data, headers, &cb); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOpt(c, cb.Code, cb.Msg, cb.Data)
}
