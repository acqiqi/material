package models

import (
	"github.com/jinzhu/gorm"
)

// 项目
type Project struct {
	Model
	ProjectName       string  `json:"project_name"` // 项目名称
	State             int     `json:"state"`        // 0 在建中 1已完成
	Remark            string  `json:"remark"`       // 备注
	Cuid              int     `json:"cuid"`
	CompanyId         int64   `json:"company_id"`         //关联企业
	Company           Company `json:"company"`            //关联企业
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
	DeletedAt         string  `json:"deleted_at"`
}

func ProjectGetInfo(id int64) (*Project, error) {
	var project Project
	err := db.Where("id = ?", id).Preload("Company").First(&project).Error
	if err != nil {
		return &Project{}, err
	}
	return &project, nil
}

// 获取项目列表
func ProjectGetLists(pageNum int, pageSize int, maps interface{}) ([]*Project, error) {
	var projects []*Project
	err := db.Model(&Project{}).Preload("Company").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&projects).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return projects, nil
}

type ProjectSelectData struct {
	Id          int64  `json:"id"`
	ProjectName string `json:"name"`
}

func ProjectGetSelect(maps string) ([]*ProjectSelectData, error) {
	var projects []*ProjectSelectData
	err := db.Model(&Project{}).Where(maps).Order("id desc").Find(&projects).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return projects, nil
}

//查询项目总数
func ProjectGetListsCount(maps interface{}) int {
	var projects []*Project
	count := 0
	db.Preload("Company").Where(maps).Find(&projects).Count(&count)
	return count
}

// 新增项目
func ProjectAdd(project *Project) error {
	project.Flag = 1
	if err := db.Create(&project).Error; err != nil {
		return err
	}
	return nil
}

// 编辑项目
func ProjectEdit(id int64, data interface{}) error {
	if err := db.Model(&Project{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
