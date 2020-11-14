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

func SendInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "下料单不存在")
		return
	}

	ps, err := models.PackingGetLists(0, 9999, utils.BuildWhere(map[string]interface{}{
		"send_id": info.Id,
		"flag":    1,
	}))

	e.ApiOk(c, "获取成功", struct {
		Info  models.Send `json:"info"`
		Table interface{} `json:"table"`
	}{
		Info:  *info,
		Table: ps,
	})
}

func SendSync(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "下料单不存在")
		return
	}

	if info.IsSync != 0 {
		e.ApiErr(c, "已经同步无法重复同步")
		return
	}

	if info.PlatformKey == "" {
		e.ApiErr(c, "无平台关联，无需同步")
		return
	}

	// 处理Callback
	if err := send_service.SyncCallback(*info, false); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "操作成功", e.GetEmptyStruct())
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

	if data.Send.Title == "" {
		e.ApiErr(c, "请输入发货标题")
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
	data.Send.PlatformKey = project.PlatformKey
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

func SendReturnInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.SendReturnGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info models.SendReturn `json:"info"`
	}{
		Info: *info,
	})
}

// 接收退貨
func SendReturnUse(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.SendReturnGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	//查询是否有自己权限
	company, _ := c.Get("company")
	if info.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}
	if info.Status == 1 {
		e.ApiErr(c, "已经接收，请勿重复接收")
		return
	}

	err = models.SendReturnEdit(info.Id, map[string]interface{}{
		"status": 1,
	})
	if err != nil {
		e.ApiErr(c, "操作失败")
		return
	}
	e.ApiOk(c, "接收成功", e.GetEmptyStruct())
}

// 补货
func SendReturnReplenish(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.SendReturnGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	//查询是否有自己权限
	company, _ := c.Get("company")
	if info.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}
	if info.Status != 1 {
		e.ApiErr(c, "未接收，无法补货")
		return
	}

	if info.IsReplenish == 1 {
		e.ApiErr(c, "已经补货，请勿重复补货")
		return
	}

	//创建补货
	m := models.Product{
		MaterialName:       info.Product.MaterialName + "- 补货",
		BlankingAttachment: info.Product.BlankingAttachment,
		Attachment:         info.Product.Attachment,
		InstallMap:         info.Product.InstallMap,
		Price:              info.Product.Price,
		Count:              info.Count, //补货数量
		ContractCount:      info.Count,
		PackCount:          0,
		SendCount:          0,
		ReturnCount:        0,
		ReceiveCount:       0,
		Unit:               info.Product.Unit,
		ProjectId:          info.Product.ProjectId,
		ProjectName:        info.Product.ProjectName,
		ReplenishmentFlag:  1, //是否补货
		ProductSubFlag:     0,
		ConfigData:         info.Product.ConfigData,
		AppendAttachment:   info.Product.AppendAttachment,
		ProjectAdditional:  info.Product.ProjectAdditional,
		Remark:             info.Product.Remark,
		Length:             info.Product.Length,
		Width:              info.Product.Width,
		Height:             info.Product.Height,
		Location:           info.Product.Location,
		Standard:           info.Product.Standard,
		ArriveDate:         info.Product.ArriveDate,
		Cuid:               info.Product.Cuid,
		CompanyId:          info.Product.CompanyId,
		Company:            info.Product.Company,
		SupplyCycle:        info.Product.SupplyCycle,
		MaterialId:         info.Product.MaterialId,
		PlatformKey:        info.Product.PlatformKey,
		PlatformUid:        info.Product.PlatformUid,
		PlatformId:         info.Product.PlatformId,
		ContractId:         info.Product.ContractId,
		SendReturnId:       info.Id,
	}
	models.ProductAdd(&m)
	err = models.SendReturnEdit(info.Id, map[string]interface{}{
		"status":       1,
		"replenish_id": m.Id,
		"is_replenish": 1,
	})
	if err != nil {
		e.ApiErr(c, "操作失败")
		return
	}
	e.ApiOk(c, "补货成功", e.GetEmptyStruct())
}
