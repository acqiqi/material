package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/packing_service"
	"material/service/send_service"
)

// 发货列表
func SendList(c *gin.Context) {
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

	if data.OptParm["status"] != nil {
		maps["status"] = data.OptParm["status"]
	}

	lists, _ := send_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.DepositoryGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

func SendCreate(c *gin.Context) {
	data := struct {
		Send  send_service.SendAdd `json:"send"`
		Links []int64              `json:"links"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	if len(data.Links) == 0 {
		e.ApiErr(c, "请选择打包数据")
		return
	}

	project, err := models.ProjectGetInfo(data.Send.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	company, _ := c.Get("company")
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	data.Send.ProjectId = project.Id
	data.Send.CompanyId = project.CompanyId

	//检测links
	maps := utils.WhereToMap(nil)
	maps["id__in"] = data.Links
	maps["company_id"] = company.(models.CompanyUsers).Company.Id
	maps["project_id"] = project.Id
	maps["flag"] = 1
	packings, err := packing_service.Select(utils.BuildWhere(maps))
	if err != nil {
		e.ApiErr(c, "获取产品列表有误")
		return
	}

	if len(packings) != len(data.Links) {
		e.ApiErr(c, "打包数据有误")
		return
	}

	for _, v := range packings {
		if v.ProjectId != project.Id {
			e.ApiErr(c, v.PackingName+" 不属于当前项目")
			return
		}
		if v.Status != 0 {
			e.ApiErr(c, v.PackingName+" 状态不可发货")
			return
		}
	}

	//处理发货
	cb, err := send_service.Add(data.Send, packings)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	e.ApiOk(c, "发货成功", cb)
}

//查看二維碼
func SendLookQrcode(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//查询打包
	send, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "打包信息不存在")
		return
	}

	company, _ := c.Get("company")
	if send.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	cb_url, err := send_service.QrcodeBuild(*send)

	if err != nil {
		e.ApiErr(c, "获取失败"+err.Error())
	} else {
		e.ApiOk(c, "获取成功", struct {
			Url string `json:"url"`
		}{
			Url: cb_url,
		})
	}
}
