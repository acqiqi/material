package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/pr_service"
)

// 获取请款TypeSelect
func PrTypeSelect(c *gin.Context) {
	e.ApiOk(c, "获取成功", pr_service.PRTypeData)
}

func PrCheckPrice(c *gin.Context) {
	data := struct {
		Id   int64 `json:"id"`
		Type int   `json:"type"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	price, sw, err := pr_service.CheckPrice(data.Id, data.Type)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "获取成功", struct {
		Price float64 `json:"price"`
		Sw    bool    `json:"sw"`
	}{
		Price: price,
		Sw:    sw,
	})
}

// 请款
func PrCreate(c *gin.Context) {
	data := pr_service.PrAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	company, _ := c.Get("company")
	//= company.(models.CompanyUsers).Company.Id

	//查询项目
	project, err := models.ProjectGetInfo(data.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}
	data.CompanyId = project.CompanyId
	user_info, _ := c.Get("user_info")
	data.Cuid = user_info.(models.Users).Cuid
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}
	if cb, err := pr_service.Add(&data); err == nil {
		e.ApiOk(c, "提交成功", cb)
	} else {
		e.ApiErr(c, "请款失败")
	}
}

// 请款列表
func PrList(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	log.Print(user_info.(models.Users))
	company, _ := c.Get("company")
	log.Println(company.(models.CompanyUsers))

	data := e.ApiPageLists{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.CheckApiPageListDefault(&data) //处理页数据

	maps := utils.WhereToMap(data.Map)
	maps["company_id"] = company.(models.CompanyUsers).Company.Id
	maps["flag"] = 1

	lists, _ := pr_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.DepositoryGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}
