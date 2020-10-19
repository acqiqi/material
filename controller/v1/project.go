package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/project_service"
	"time"
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
	//maps["bind_state"] = 1
	if data.OptParm["state"] != nil {
		maps["state"] = data.OptParm["state"]
	}
	if data.OptParm["bind_state"] != nil {
		maps["bind_state"] = data.OptParm["bind_state"]
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
func ProjectCreate(c *gin.Context) {
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
	data.Status = 1 //自动绑定
	cb, err := project_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	p, _ := models.ProjectGetInfo(cb.Id)

	e.ApiOk(c, "创建成功", p)
}

// 编辑项目
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
	if project.IsPlatform == 1 {
		e.ApiErr(c, "公装系统同步无法修改")
		return
	}

	company, _ := c.Get("company")
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		log.Println(company.(models.CompanyUsers))
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

// 获取项目Select
func ProjectSelect(c *gin.Context) {
	company, _ := c.Get("company")
	lists, err := project_service.SelectLists(company.(models.CompanyUsers).Company.Id)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "获取成功", lists)
	return
}

func ProjectInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	project, err := models.ProjectGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info models.Project `json:"info"`
	}{
		Info: *project,
	})
}

//接收项目
func ProjectReceive(c *gin.Context) {
	data := e.ApiId{}
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
	if project.IsPlatform == 0 {
		e.ApiErr(c, "非三方项目无法接收")
		return
	}
	if project.BindState != 0 {
		e.ApiErr(c, "已接收项目请勿重复接收")
		return
	}
	company, _ := c.Get("company")
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		log.Println(company.(models.CompanyUsers))
		return
	}

	project.BindState = 1
	project.ReceiveTime = utils.Time{Time: time.Now()}
	err = models.ProjectEdit(project.Id, project)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//查询Platform
	platform, err := models.PlatformGetInfoOrKey(project.PlatformKey)
	if err == nil {
		callback := e.HttpCallbackData{
			Code:   0,
			Msg:    "Receive Success",
			Action: e.PLATFORM_ACTION_PROJECT_RECEIVE, //接收
			Data: e.PlatformProjectReceiveCallback{
				Id:          project.Id,
				ProjectName: project.ProjectName,
				State:       project.State,
				CompanyId:   project.CompanyId,
				Company:     project.Company,
				BindState:   project.BindState,
				PlatformKey: project.PlatformKey,
				PlatformUid: project.PlatformUid,
				PlatformId:  project.PlatformId,
				CreatedAt:   project.CreatedAt,
				Status:      1,
			},
		}
		cb_url := platform.MessageCallbackUrl
		log.Println(utils.JsonEncode(callback))
		c_data := new(e.HttpCallbackData)
		err = utils.HttpPostJson(cb_url, callback, &c_data)
		if err != nil {
			log.Println(err.Error())
		}
	}
	e.ApiOk(c, "操作成功", e.GetEmptyStruct())
}
