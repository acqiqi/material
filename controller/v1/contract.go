package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/contract_service"
)

// 获取合同Select
func ContractSelect(c *gin.Context) {
	company, _ := c.Get("company")
	lists, err := contract_service.SelectLists(company.(models.CompanyUsers).Company.Id)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "获取成功", lists)
	return
}

// 创建合同
func ContractCreate(c *gin.Context) {
	data := contract_service.ContractAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	log.Println(data)
	user_info, _ := c.Get("user_info")
	data.Cuid = int(user_info.(models.Users).Cuid)
	company, _ := c.Get("company")
	data.CompanyId = company.(models.CompanyUsers).Company.Id
	data.BindState = 1
	//手动创建默认绑定

	cb, err := contract_service.Add(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	p, _ := models.ProjectGetInfo(cb.Id)
	e.ApiOk(c, "创建成功", p)
}

// 编辑合同
func ContractEdit(c *gin.Context) {
	data := contract_service.ContractAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	// 检测是否有项目
	contract, err := models.ContractInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	if contract.IsPlatform == 1 {
		e.ApiErr(c, "三方平台合同无法编辑")
		return
	}

	company, _ := c.Get("company")
	if contract.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		log.Println(company.(models.CompanyUsers))
		return
	}

	cb, err := contract_service.Edit(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	p, _ := models.ProjectGetInfo(cb.Id)
	e.ApiOk(c, "编辑成功", p)
}

// 合同列表
func ContractList(c *gin.Context) {
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

	lists, _ := contract_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.ProjectGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

func ContractInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	contract, err := models.ContractInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "合同不存在")
		return
	}

	//查詢合同配置
	cc, err := models.ContractConfigGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "配置查詢有誤")
		return
	}
	e.ApiOk(c, "获取成功", struct {
		Info   models.Contract       `json:"info"`
		Config models.ContractConfig `json:"config"`
	}{
		Info:   *contract,
		Config: *cc,
	})
}
