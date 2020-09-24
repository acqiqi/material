package controller

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/app"
	"material/lib/dd"
	"material/lib/e"
	"material/lib/setting"
	"material/lib/utils"
	"material/service/auth_service"
)

type auth struct {
	Username string `json:"username" valid:"Required; MaxSize(32)"`
	Password string `json:"password" valid:"Required; MaxSize(32)"`
}

type mobile_login struct {
	Mobile string `json:"mobile" valid:"Required; MaxSize(32)"`
	Code   string `json:"code" valid:"Required; MaxSize(32)"`
}

func AutoLogin(c *gin.Context) {
	data := mobile_login{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	log.Println(setting.PlatformSetting.PlatformKey)
	// 表单验证
	valid := validation.Validation{}
	ok, _ := valid.Valid(&data)
	if !ok {
		app.MarkErrors(valid.Errors)
		e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
		return
	}

	//处理请求
	headers := make(map[string]string)
	headers["PlatformKey"] = setting.PlatformSetting.PlatformKey
	// 注册接口
	url := setting.PlatformSetting.ApiBaseUrl + dd.DD_API_AUTOLOGIN

	cb := e.ApiJson{}

	if err := utils.HttpPostJsonHeader(url, data, headers, &cb); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	if cb.Code != 0 {
		e.ApiErr(c, cb.Msg)
		return
	}
	cc := cb.Data.(map[string]interface{})["token"].(string)
	log.Println(cc)
	e.ApiOk(c, cb.Msg, cb.Data)
}

//获取注册手机验证码
func GetAutoLoginSMSCode(c *gin.Context) {
	data := struct {
		Mobile string `json:"mobile"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//处理请求
	headers := make(map[string]string)
	headers["PlatformKey"] = setting.PlatformSetting.PlatformKey
	// 注册接口
	url := setting.PlatformSetting.ApiBaseUrl + dd.DD_API_AUTOLOGIN_GETSMS

	cb := e.ApiJson{}
	if err := utils.HttpPostJsonHeader(url, data, headers, &cb); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOpt(c, cb.Code, cb.Msg, cb.Data)
}

func Login(c *gin.Context) {
	log.Print(utils.PasswordEncode("123qwe"))
	data := auth{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	valid := validation.Validation{}
	ok, _ := valid.Valid(&data)

	if !ok {
		app.MarkErrors(valid.Errors)
		e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
		return
	}

	authService := auth_service.Auth{Username: data.Username, Password: data.Password}
	isExist, err := authService.Check()
	if err != nil {
		e.ApiOpt(c, e.ERROR_AUTH_CHECK_PASSWORD, e.GetMsg(e.ERROR_AUTH_CHECK_PASSWORD), e.GetEmptyStruct())
		return
	}

	if isExist == 0 {
		e.ApiOpt(c, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH), e.GetEmptyStruct())
		return
	}

	token, err := utils.GenerateToken(isExist)
	if err != nil {
		e.ApiOpt(c, e.ERROR_AUTH_TOKEN, e.GetMsg(e.ERROR_AUTH_TOKEN), e.GetEmptyStruct())

		return
	}

	e.ApiOk(c, "登录成功", struct {
		Token string `json:"token"`
	}{Token: token})

}
