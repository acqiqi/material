package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
)

func CompanyGetInfo(c *gin.Context) {
	data := struct {
		CompanyKey string `json:"company_key"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	company, err := models.CompanyGetInfoOrKey(data.CompanyKey)
	if err != nil {
		e.ApiErr(c, "查无此企业")
		return
	}
	e.ApiOk(c, "获取成功", company)
}

// 链接
func CompanyBindLink(c *gin.Context) {
	data := struct {
		PlatformUid string `json:"platform_uid"`
		CompanyKey  string `json:"company_key"`
		DataOrigin  string `json:"data_origin"`
		SupplierId  string `json:"supplier_id"`
		Opt         string `json:"opt"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	company, err := models.CompanyGetInfoOrKey(data.CompanyKey)
	if err != nil {
		e.ApiErr(c, "查无此企业")
		return
	}

	platform_key, _ := c.Get("platform_key")

	//检测是否已经绑定
	_, err = models.PlatformCompanyCheck(company.Id, platform_key.(string), data.PlatformUid)
	if err == nil {
		e.ApiErr(c, "已经绑定，请勿重复绑定")
		return
	}

	m := models.PlatformCompany{
		CompanyId:   company.Id,
		Company:     models.Company{},
		PlatformKey: platform_key.(string),
		PlatformUid: data.PlatformUid,
		CompanyKey:  data.CompanyKey,
		DataOrigin:  data.DataOrigin,
		SupplierId:  data.SupplierId,
	}

	err = models.PlatformCompanyAdd(&m)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "绑定成功", company)
}

// 获取企业列表
func CompanyList(c *gin.Context) {
	data := struct {
		PlatformUid string `json:"platform_uid"`
		CompanyKey  string `json:"company_key"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	maps := utils.WhereToMap(nil)
	maps["platform_uid"] = data.PlatformUid
	maps["company_key"] = data.CompanyKey
	maps["flag"] = 1
	lists, err := models.PlatformCompanyGetLists(utils.BuildWhere(maps))
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Lists interface{} `json:"lists"`
	}{Lists: lists})
	return
}

// 删除企业连接
func CompanyDelete(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	// 检测是否存在
	info, err := models.PlatformCompanyGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "企业不存在")
		return
	}

	platform_key, _ := c.Get("platform_key")
	if platform_key.(string) != info.PlatformKey {
		e.ApiErr(c, "非法请求")
		return
	}

	models.PlatformCompanyDelete(info)
	e.ApiOk(c, "操作成功", e.GetEmptyStruct())
}
