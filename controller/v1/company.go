package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/company_service"
	"material/service/depository_service"
)

// 创建公司
func CompanyCreate(c *gin.Context) {
	data := company_service.CompanyAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	user_info, _ := c.Get("user_info")
	data.Cuid = int(user_info.(models.Users).Id)
	cb, err := company_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	//创建默认仓库
	d := depository_service.DepositoryAdd{
		Name:      "默认仓库",
		Desc:      "",
		Address:   "",
		CompanyId: cb.Id,
		Status:    1,
	}
	depository_service.Add(&d)

	e.ApiOk(c, "创建成功", cb)
}

// 获取我的公司详情
func CompanyMyInfo(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.CompanyUsersGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "公司不存在")
		return
	}

	user_info, _ := c.Get("user_info")
	if info.Cuid != int(user_info.(models.Users).Id) {
		e.ApiErr(c, "非法请求")
		return
	}
	e.ApiOk(c, "获取成功", struct {
		Info     models.CompanyUsers `json:"info"`
		AuthPics []string            `json:"auth_pics"`
	}{
		Info:     *info,
		AuthPics: utils.JsonDecodeUrls(info.Company.AuthPics),
	})
}

// 获取我的公司
func CompanyMyList(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	list, err := models.CompanyUsersGetMyList(user_info.(models.Users).Id)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  1,
		Limit: 9999,
		Lists: list,
		Total: len(list),
	})
	return
}
