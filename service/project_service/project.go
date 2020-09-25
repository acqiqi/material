package project_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

type ProjectAdd struct {
	Id                int64   `json:"id"`
	ProjectName       string  `json:"project_name"` // 项目名称
	State             int     `json:"state"`        // 0 在建中 1已完成
	Remark            string  `json:"remark"`       // 备注
	Cuid              int     `json:"cuid"`
	CompanyId         int64   `json:"company_id"`         //关联企业
	AppendAttachment  string  `json:"append_attachment"`  // 附加信息
	ReceiverMembers   string  `json:"receiver_members"`   // 收货人ids
	BindState         int     `json:"bind_state"`         // 绑定工程账号的状态,手动新建的项目默认已绑定.未绑定,待处理,已绑定
	BindType          int     `json:"bind_type"`          // 绑定的项目属于哪个系统(项目版，企业版，还是新建的) 新建,项目版,企业版,私有化定制 这里可以Platform控制
	DataOrigin        string  `json:"data_origin"`        // 绑定工程的账号名称的数据来源(例如:码里公装，中建深装)
	ProjectAccount    string  `json:"project_account"`    // 绑定工程的账号名称
	SupplierAccountid int     `json:"supplier_accountid"` // 绑定工程供应商id
	ProjectAccountid  int     `json:"project_accountid"`  // 绑定工程账号(项目)的id
	ContractMoney     float64 `json:"contract_money"`     // 合同总金额
	ReceivedMoney     float64 `json:"received_money"`     // 已回款总金额
	ReceiptMoney      float64 `json:"receipt_money"`      // 已开票总金额
}

var state = map[int]string{
	0: "在建中",
	1: "已完成",
}

var bindState = map[int]string{
	0: "未绑定",
	1: "待处理",
	2: "已绑定",
}

var bindType = map[int]string{
	0: "新建",
	1: "项目版",
	2: "企业版",
	3: "私有化定制",
}

//新增企业
func Add(data *ProjectAdd) (*models.Project, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.Cuid, "Cuid").Message("CUID不能为空！")
	valid.Required(data.ProjectName, "ProjectName").Message("请输入项目名称")
	valid.Required(data.ContractMoney, "ContractMoney").Message("请输入正确合同金额")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	//检测Company是否存在
	company, err := models.CompanyUsersGetInfoOrCompanyId(data.CompanyId)
	if err != nil {
		return nil, err
	}
	// 检测有用户权限
	if company.Cuid != data.Cuid {
		return nil, errors.New("非法请求")
	}

	log.Println("???")
	model := models.Project{}
	model.Cuid = data.Cuid
	model.ProjectName = data.ProjectName
	model.State = utils.CheckStatusIndex(state, data.State)
	model.Remark = data.Remark
	model.CompanyId = company.Company.Id
	model.AppendAttachment = data.AppendAttachment
	model.ReceiverMembers = data.ReceiverMembers
	model.BindState = utils.CheckStatusIndex(bindState, data.BindState)
	model.BindType = utils.CheckStatusIndex(bindType, data.BindType)
	model.DataOrigin = data.DataOrigin
	model.ProjectAccount = data.ProjectAccount
	model.SupplierAccountid = data.SupplierAccountid
	model.ProjectAccountid = data.ProjectAccountid
	model.ContractMoney = data.ContractMoney
	model.ReceivedMoney = 0
	model.ReceiptMoney = 0
	if err := models.ProjectAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 编辑项目
func Edit(data *ProjectAdd) (*models.Project, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.ProjectName, "ProjectName").Message("请输入项目名称")
	valid.Required(data.ContractMoney, "ContractMoney").Message("请输入正确合同金额")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	model := models.Project{}
	model.Id = data.Id
	model.ProjectName = data.ProjectName
	model.State = utils.CheckStatusIndex(state, data.State)
	model.Remark = data.Remark
	model.AppendAttachment = data.AppendAttachment
	model.ReceiverMembers = data.ReceiverMembers
	model.BindState = utils.CheckStatusIndex(bindState, data.BindState)
	model.BindType = utils.CheckStatusIndex(bindType, data.BindType)
	model.DataOrigin = data.DataOrigin
	model.ProjectAccount = data.ProjectAccount
	model.SupplierAccountid = data.SupplierAccountid
	model.ProjectAccountid = data.ProjectAccountid
	model.ContractMoney = data.ContractMoney

	if err := models.ProjectEdit(model.Id, model); err != nil {
		return nil, err
	}

	return &model, nil
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Project, error) {
	offset := (page - 1) * limit
	return models.ProjectGetLists(offset, limit, maps)
}
