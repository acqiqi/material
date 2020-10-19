package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/product_service"
	"material/service/send_service"
)

// 产品 材料列表
func ProductList(c *gin.Context) {
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

	if data.OptParm["replenishment_flag"] != nil {
		maps["replenishment_flag"] = 1
	} else {
		maps["replenishment_flag"] = 0
	}

	lists, _ := product_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.ProductGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

// 退货列表
func ReturnProductList(c *gin.Context) {
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

	lists, _ := send_service.ApiListsReturn(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.SendReturnGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

func ProductInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.ProductGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "产品不存在")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info models.Product `json:"info"`
	}{
		Info: *info,
	})
}

//// 创建产品
//func ProductCreate(c *gin.Context) {
//	data := product_service.ProductAdd{}
//	if err := c.BindJSON(&data); err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	user_info, _ := c.Get("user_info")
//	data.Cuid = int(user_info.(models.Users).Id)
//	company, _ := c.Get("company")
//	data.CompanyId = company.(models.CompanyUsers).Company.Id
//	//检测项目是否存在
//
//	cb, err := project_service.Add(&data)
//	if err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//
//	p, _ := models.ProjectGetInfo(cb.Id)
//
//	e.ApiOk(c, "创建成功", p)
//}
//
//// 编辑产品
//func ProductEdit(c *gin.Context) {
//	data := project_service.ProjectAdd{}
//	if err := c.BindJSON(&data); err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	// 检测是否有项目
//	project, err := models.ProjectGetInfo(data.Id)
//	if err != nil {
//		e.ApiErr(c, "项目不存在")
//		return
//	}
//	company, _ := c.Get("company")
//	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
//		e.ApiErr(c, "非法请求")
//		return
//	}
//
//	cb, err := project_service.Edit(&data)
//	if err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	p, _ := models.ProjectGetInfo(cb.Id)
//	e.ApiOk(c, "编辑成功", p)
//}

// 材料类型列表
func ProductClassList(c *gin.Context) {
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

	lists, _ := product_service.ApiListsClass(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.ProductClassGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

// 创建材料类型
func ProductClassCreate(c *gin.Context) {
	data := product_service.ProductClassAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	company, _ := c.Get("company")
	data.CompanyId = company.(models.CompanyUsers).Company.Id
	cb, err := product_service.AddClass(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.ApiOk(c, "创建成功", cb)
}

// 编辑材料类型
func ProductClassEdit(c *gin.Context) {
	data := product_service.ProductClassAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	// 检测是否有项目
	d, err := models.ProductClassGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	company, _ := c.Get("company")
	if d.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	err = product_service.EditCalss(&data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	p, _ := models.ProductClassGetInfo(data.Id)
	e.ApiOk(c, "编辑成功", p)
}

// 刪除Class
func ProductClassDelete(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	//检测是否存在
	info, err := models.ProductClassGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "数据不存在")
		return
	}
	company, _ := c.Get("company")
	if info.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}
	info.Flag = -1
	models.ProductClassEdit(data.Id, info)
	e.ApiOk(c, "操作成功", e.GetEmptyStruct())
}

// 材料表格
func ProductTable(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	company, _ := c.Get("company")
	list, err := product_service.Tables(data.Id, company.(models.CompanyUsers).Company.Id)
	if err != nil {
		e.ApiErr(c, "非法请求")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Table interface{} `json:"table"`
	}{
		Table: list,
	})
}
