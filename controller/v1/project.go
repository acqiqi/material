package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/project_service"
)

// 获取我的项目列表
func ProjectList(c *gin.Context) {
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

	if data.OptParm["state"] != nil {
		maps["state"] = data.OptParm["state"]
	}

	lists, _ := project_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.ProjectGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

// 创建项目
func ProjecttCreate(c *gin.Context) {
	data := project_service.ProjectAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	user_info, _ := c.Get("user_info")
	data.Cuid = int(user_info.(models.Users).Id)
	company, _ := c.Get("company")
	data.CompanyId = company.(models.CompanyUsers).Company.Id
	//手动创建默认绑定
	data.BindState = 1
	data.BindType = 0
	data.DataOrigin = "自建"

	cb, err := project_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	p, _ := models.ProjectGetInfo(cb.Id)

	e.ApiOk(c, "创建成功", p)
}

func ProjectEdit(c *gin.Context) {
	data := project_service.ProjectAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	// 检测是否有项目
	project, err := models.ProjectGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}
	company, _ := c.Get("company")
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	cb, err := project_service.Edit(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	p, _ := models.ProjectGetInfo(cb.Id)
	e.ApiOk(c, "编辑成功", p)
}
