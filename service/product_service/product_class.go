package product_service

import (
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

//材料类型表
type ProductClassAdd struct {
	Id          int64              `json:"id"`
	ClassName   string             `json:"class_name"` // 材料类型名称
	Desc        string             `json:"desc"`       // 描述
	CatsId      int                `json:"cats_id"`
	ProductCats models.ProductCats `gorm:"ForeignKey:CatsId" json:"product_cats"`
	CompanyId   int64              `json:"company_id"`
	Contract    models.Contract    `gorm:"ForeignKey:ContractId" json:"contract"`
	Cuid        int                `json:"cuid"`
}

// 获取Api列表
func ApiListsClass(page int, limit int, maps string) ([]*models.ProductClass, error) {
	offset := (page - 1) * limit
	return models.ProductClassGetLists(offset, limit, maps)
}

//新增
func AddClass(data *ProductClassAdd) (*models.ProductClass, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.ClassName, "ClassName").Message("请输入名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	log.Println("???")
	model := models.ProductClass{
		ClassName: data.ClassName,
		Desc:      data.Desc,
		CatsId:    0,
		CompanyId: data.CompanyId,
	}
	if err := models.ProductClassAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 编辑项目
func EditCalss(data *ProductClassAdd) error {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.ClassName, "ClassName").Message("请输入名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	c, err := models.ProductClassGetInfo(data.Id)
	if err != nil {
		return err
	}
	//model.Id = data.CompanyId
	model := make(map[string]interface{})
	model["ClassName"] = data.ClassName
	model["Desc"] = data.Desc

	log.Println(model)
	if err := models.ProductClassEdit(c.Id, model); err != nil {
		return err
	}
	return nil
}
