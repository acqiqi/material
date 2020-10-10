package packing_service

import (
	"github.com/astaxie/beego/validation"
	uuid "github.com/satori/go.uuid"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

//打包
type PackingAdd struct {
	Id          int64  `json:"id"`
	PackingName string `orm:"packing_name"` // 包装名称
	SerialNo    string `orm:"serial_no"`    // 包装编号
	Count       int    `orm:"count"`        // 产品总数
	ReturnCount int    `orm:"return_count"` // 包装下退货数量
	Remark      string `orm:"remark"`       // 描述
	CompanyId   int64  `orm:"company_id"`
	ProductId   int64  `orm:"product_id"`
	MaterialId  int64  `orm:"material_id"`

	ContractId int64           `json:"contract_id"` //合同
	Contract   models.Contract `gorm:"ForeignKey:ContractId" json:"contract"`
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Packing, error) {
	offset := (page - 1) * limit
	return models.PackingGetLists(offset, limit, maps)
}

//新增
func Add(data PackingAdd, links []PackingProductAdd) (*models.Packing, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.PackingName, "PackingName").Message("请输入打包名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	log.Println("???")

	serial_no := uuid.NewV4().String()

	model := models.Packing{
		PackingName: data.PackingName,
		SerialNo:    serial_no,
		Count:       data.Count,
		ReturnCount: data.ReturnCount,
		Remark:      data.Remark,
		CompanyId:   data.CompanyId,
		ProductId:   data.ProductId,
		MaterialId:  data.MaterialId,
	}
	//创建事务
	t := models.NewTransaction()

	if err := models.PackingAddT(&model, t); err != nil {
		return nil, err
	}

	//处理链接
	for _, v := range links {
		v.PackingId = model.Id

		link_model := models.PackingProduct{
			PackingId:     v.PackingId,
			CompanyId:     v.CompanyId,
			OrderReturnid: v.OrderReturnid,
			ProductId:     v.ProductId,
			MaterialId:    v.MaterialId,
			Count:         v.Count,
			ReturnCount:   0,
			MaterialName:  v.MaterialName,
			ContractId:    v.ContractId,
		}
		models.PackingProductAddT(&link_model, t)
	}
	t.Commit()
	return &model, nil
}

// 编辑项目
func Edit(data *PackingAdd) error {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.PackingName, "PackingName").Message("请输入打包名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	c, err := models.PackingGetInfo(data.Id)
	if err != nil {
		return err
	}
	//model.Id = data.CompanyId
	model := make(map[string]interface{})
	model["PackingName"] = data.PackingName
	//model["SerialNo"] = data.SerialNo
	model["Count"] = data.Count
	model["ReturnCount"] = data.ReturnCount
	model["Remark"] = data.Remark
	model["CompanyId"] = data.CompanyId
	model["ProductId"] = data.ProductId
	model["MaterialId"] = data.MaterialId

	log.Println(model)
	if err := models.PackingEdit(c.Id, model); err != nil {
		return err
	}
	return nil
}
