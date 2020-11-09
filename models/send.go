package models

import "github.com/jinzhu/gorm"

// 发货表
type Send struct {
	Model
	SendNo            string   `json:"send_no"`            // 订单编号
	Count             float64  `json:"count"`              // 发货总数
	ReturnCount       float64  `json:"return_count"`       // 退货数量
	ActualReceiver    string   `json:"actual_receiver"`    // 签收人
	Address           string   `json:"address"`            // 收货地址
	ReceiveAttachment string   `json:"receive_attachment"` // 收货附件
	ReceiveDate       string   `json:"receive_date"`       // 收货时间
	ReceiveCount      float64  `json:"receive_count"`      // 收货总数量
	ReceiveRemark     string   `json:"receive_remark"`     // 收货备注
	Remark            string   `json:"remark"`             // 备注
	CompanyId         int64    `json:"company_id"`         //
	Company           Contract `gorm:"ForeignKey:CompanyId" json:"company"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	ExpressNo string `json:"express_no"`

	Status int `json:"status"` //0未签收 1签收

	Packing []Packing `gorm:"ForeignKey:SendId" json:"packing"`

	IsSync      int    `json:"is_sync"`      // 是否同步 如果platform存在就需要同步
	PlatformKey string `json:"platform_key"` // 平台key

	ReceiveMobile string `json:"receive_mobile"`
}

// 发货详情
func SendGetInfo(id int64) (*Send, error) {
	var d Send
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").Preload("Packing").Preload("Project").First(&d).Error
	if err != nil {
		return &Send{}, err
	}
	return &d, nil
}
func SendGetInfoT(id int64, t *gorm.DB) (*Send, error) {
	var d Send
	err := t.Where("id = ? AND flag = 1", id).Preload("Company").Preload("Packing").Preload("Project").First(&d).Error
	if err != nil {
		return &Send{}, err
	}
	return &d, nil
}

// 发货列表
func SendGetLists(pageNum int, pageSize int, maps interface{}) ([]*Send, error) {
	var d []*Send
	err := db.Model(&Send{}).Preload("Company").Preload("Project").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []*Send{}, err
	}
	return d, nil
}

// 查询仓库列表数量
func SendListsCount(maps interface{}) int {
	var d []*Send
	count := 0
	db.Preload("Company").Where(maps).Find(&d).Count(&count)
	return count
}

// 新增仓库
func SendAdd(d *Send) error {
	d.Flag = 1
	if err := db.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

func SendAddT(packing *Send, t *gorm.DB) error {
	packing.Flag = 1
	if err := t.Create(&packing).Error; err != nil {
		return err
	}
	return nil
}

// 编辑发货
func SendEdit(id int64, data interface{}) error {
	if err := db.Model(&Send{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 编辑发货
func SendEditT(id int64, data interface{}, t *gorm.DB) error {
	if err := t.Model(&Send{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
