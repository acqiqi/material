package depository_service

import (
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

type DepositoryAdd struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`       // 仓库名称
	Desc      string `json:"desc"`       // 描述
	Address   string `json:"address"`    // 仓库地址
	CompanyId int64  `json:"company_id"` // 企业id
	Status    int    `json:"status"`     // 状态 0停用 1正常
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Depository, error) {
	offset := (page - 1) * limit
	return models.DepositoryGetLists(offset, limit, maps)
}

//新增
func Add(data *DepositoryAdd) (*models.Depository, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.Name, "Name").Message("请输入仓库名")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	log.Println("???")
	model := models.Depository{}
	model.Name = data.Name
	model.Desc = data.Desc
	model.Address = data.Address
	model.Status = data.Status
	model.CompanyId = data.CompanyId
	if err := models.DepositoryAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 编辑项目
func Edit(data *DepositoryAdd) error {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.Name, "Name").Message("仓库名")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	c, err := models.DepositoryGetInfo(data.Id)
	if err != nil {
		return err
	}
	//model.Id = data.CompanyId
	model := make(map[string]interface{})
	model["Name"] = data.Name
	model["Address"] = data.Address
	model["Desc"] = data.Desc
	model["Status"] = data.Status

	log.Println(model)
	if err := models.DepositoryEdit(c.Id, model); err != nil {
		return err
	}
	return nil
}
