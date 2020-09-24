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
	"material/models"
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

type user_info_api struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Mobile   string `json:"mobile"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		Gender   string `json:"gender"`
		Status   int    `json:"status"`
		Score    int    `json:"score"`
		Money    string `json:"money"`
		OkMoney  string `json:"ok_money"`
		NoMoney  string `json:"no_money"`
	} `json:"data"`
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
	//获取UserInfo
	user_info := user_info_api{}
	url = setting.PlatformSetting.ApiBaseUrl + dd.DD_API_AUTH_GET_USER_INFO
	headers["Authorization"] = cb.Data.(map[string]interface{})["token"].(string) //插入Token
	if err := utils.HttpPostJsonHeader(url, e.GetEmptyStruct(), headers, &user_info); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	if user_info.Code != 0 {
		e.ApiErr(c, user_info.Msg)
		return
	}

	uid := int64(0)
	//查询本地是否注册
	my_user_info, err := models.GetUsersInfoCuid(user_info.Data.ID)
	if err != nil {
		//直接注册
		user_model := models.Users{
			Cuid:     int(user_info.Data.ID),
			Nickname: user_info.Data.Nickname,
			Avatar:   user_info.Data.Avatar,
			MUserKey: models.GetMUserKey(),
		}
		if err := models.AddUsers(&user_model); err != nil {
			e.ApiErr(c, err.Error())
			return
		}
		uid = int64(user_model.Cuid)
	} else {
		uid = int64(my_user_info.Cuid)
	}

	log.Println(uid)
	token, err := utils.GenerateToken(uid)
	if err != nil {
		e.ApiOpt(c, e.ERROR_AUTH_TOKEN, e.GetMsg(e.ERROR_AUTH_TOKEN), e.GetEmptyStruct())
		return
	}

	e.ApiOk(c, "登录成功", struct {
		Token string `json:"token"`
	}{Token: token})
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
