package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/depository_service"
)

// 仓库列表
func DepositoryList(c *gin.Context) {
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

	lists, _ := depository_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.ProductGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

// 创建仓库
func DepositoryCreate(c *gin.Context) {
	data := depository_service.DepositoryAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	company, _ := c.Get("company")
	data.CompanyId = company.(models.CompanyUsers).Company.Id
	cb, err := depository_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "创建成功", cb)
}

// 编辑仓库
func DepositoryEdit(c *gin.Context) {
	data := depository_service.DepositoryAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	// 检测是否有项目
	d, err := models.DepositoryGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	company, _ := c.Get("company")
	if d.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	err = depository_service.Edit(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	p, _ := models.DepositoryGetInfo(data.Id)
	e.ApiOk(c, "编辑成功", p)
}
