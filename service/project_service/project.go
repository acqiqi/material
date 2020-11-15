package project_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
	"strconv"
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
	Status            int     `json:"status"`             //状态 1已接收 如果是自建的会自动设置1
	PlatformKey       string  `json:"platform_key"`       // 平台key
	PlatformUid       string  `json:"platform_uid"`       // 平台用户id
	PlatformId        string  `json:"platform_id"`        // 平台用户id
	IsPlatform        int     `json:"is_platform"`        // 是否三方平台同步

	ReceiverAddress string `json:"receiver_address"` //收货地址
	SupplierId      string `json:"supplier_id"`
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

//新增项目
func Add(data *ProjectAdd) (*models.Project, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.Cuid, "Cuid").Message("CUID不能为空！")
	valid.Required(data.ProjectName, "ProjectName").Message("请输入项目名称")
	if data.PlatformKey != "" {
		valid.Required(data.PlatformId, "PlatformId").Message("请输入项目唯一id")
	} else {
		valid.Required(data.ContractMoney, "ContractMoney").Message("请输入正确合同金额")
	}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	company, err := models.CompanyUsersGetInfoOrCompanyId(data.CompanyId)
	if err != nil {
		return nil, errors.New("请选择正确的企业")
	}
	//检测Company是否存在
	if data.PlatformKey == "" {
		// 检测有用户权限
		if company.Cuid != data.Cuid {
			return nil, errors.New("非法请求")
		}
	}

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
	model.PlatformId = data.PlatformId
	model.PlatformKey = data.PlatformKey
	model.PlatformUid = data.PlatformUid
	model.IsPlatform = data.IsPlatform
	model.ReceiverAddress = data.ReceiverAddress
	model.SupplierId = data.SupplierId
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
	model, err := models.ProjectGetInfo(data.Id)
	if err != nil {
		return nil, err
	}

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

	log.Println(model)
	if err := models.ProjectEdit(model.Id, model); err != nil {
		return nil, err
	}

	return model, nil
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]map[string]interface{}, error) {
	offset := (page - 1) * limit
	list, err := models.ProjectGetLists(offset, limit, maps)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	cb := make([]map[string]interface{}, len(list))
	for i, v := range list {
		cb_item := map[string]interface{}{
			"id":           v.Id,
			"created_at":   v.CreatedAt,
			"updated_at":   v.UpdatedAt,
			"flag":         v.Flag,
			"project_name": v.ProjectName,
			"state":        v.State,
			"remark":       v.Remark,
			"cuid":         v.Cuid,
			"company_id":   v.CompanyId,
			"company": map[string]interface{}{
				"id":          v.Company.Id,
				"created_at":  v.Company.CreatedAt,
				"updated_at":  v.Company.UpdatedAt,
				"flag":        v.Company.Flag,
				"name":        v.Company.Name,
				"mobile":      v.Company.Mobile,
				"tel":         v.Company.Tel,
				"address":     v.Company.Address,
				"desc":        v.Company.Desc,
				"auth_pics":   v.Company.AuthPics,
				"company_key": v.Company.CompanyKey,
			},
			"append_attachment":  v.AppendAttachment,
			"receiver_members":   v.ReceiverMembers,
			"bind_state":         v.BindState,
			"bind_type":          v.BindType,
			"data_origin":        v.DataOrigin,
			"project_account":    v.ProjectAccount,
			"supplier_accountid": v.SupplierAccountid,
			"project_accountid":  v.ProjectAccountid,
			"contract_money":     v.ContractMoney,
			"received_money":     v.ReceivedMoney,
			"receipt_money":      v.ReceiptMoney,
			"status":             v.Status,
			"platform_key":       v.PlatformKey,
			"platform_uid":       v.PlatformUid,
			"platform_id":        v.PlatformId,
			"is_platform":        v.IsPlatform,
			"receive_time":       v.ReceiveTime,
		}
		cb[i] = cb_item
	}
	return cb, nil
}

type ProjectSelectData struct {
	Id          int64  `json:"id"`
	ProjectName string `json:"name"`
}

// 获取Select
func SelectLists(company_id int64) ([]ProjectSelectData, error) {
	lists, err := models.ProjectGetSelect("flag =1 AND company_id = " + strconv.Itoa(int(company_id)))
	if err != nil {
		return nil, err
	}
	cb := make([]ProjectSelectData, len(lists))
	for i := 0; i < len(lists); i++ {
		cb[i] = ProjectSelectData{
			Id:          lists[i].Id,
			ProjectName: lists[i].ProjectName,
		}
	}
	return cb, nil
}
