package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"material/lib/utils"
	"time"
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
	Status            int     `json:"status"`             //状态 1已接收 如果是自建的会自动设置1
	PlatformKey       string  `json:"platform_key"`       // 平台key
	PlatformUid       string  `json:"platform_uid"`       // 平台用户id
	PlatformId        string  `json:"platform_id"`        // 平台用户id
	IsPlatform        int     `json:"is_platform"`        // 是否三方平台同步

	ReceiveTime   utils.Time      `json:"receive_time"` //接收时间
	ReceiverUsers []ReceiverUsers `gorm:"ForeignKey:ProjectId" json:"receiver_users"`

	ReceiverAddress string `json:"receiver_address"` //收货地址

}

func ProjectGetInfo(id int64) (*Project, error) {
	var project Project
	err := db.Where("id = ? AND flag =1", id).Preload("Company").First(&project).Error
	if err != nil {
		return &Project{}, err
	}
	return &project, nil
}

// 三方检测是否存在
func ProjectCheck(platform_id string, platform_key string, platform_uid string) (*Project, error) {
	log.Println(platform_id, platform_key, platform_uid)
	var pc Project
	err := db.Where("platform_id = ? AND platform_key = ? AND platform_uid = ? AND flag =1",
		platform_id, platform_key, platform_uid).Preload("Company").First(&pc).Error
	if err != nil {
		return &Project{}, err
	}
	return &pc, nil
}

// 获取项目列表
func ProjectGetLists(pageNum int, pageSize int, maps interface{}) ([]*Project, error) {
	var projects []*Project
	err := db.Model(&Project{}).Preload("Company").Preload("ReceiverUsers").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&projects).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return projects, nil
}

func ProjectGetSelect(maps string) ([]*Project, error) {
	var projects []*Project
	err := db.Where(maps).Order("id desc").Find(&projects).Error
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

// 查询项目总数
func ProjectGetCount(company_id int64, begin_time, end_time time.Time) int {
	var projects []*Project
	count := 0
	db.Preload("Company").Where("company_id = ? AND created_at BETWEEN ? AND ?",
		company_id, begin_time, end_time).Find(&projects).Count(&count)
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
