package models

import "github.com/jinzhu/gorm"

// 发货表
type Send struct {
	Model
	SendNo            string   `json:"send_no"`            // 订单编号
	Count             float64  `json:"count"`              // 发货总数
	ActualReceiver    string   `json:"actual_receiver"`    // 签收人
	Address           string   `json:"address"`            // 收货地址
	RecevieAttachment string   `json:"recevie_attachment"` // 收货附件
	RecevieDate       string   `json:"recevie_date"`       // 收货时间
	RecevieCount      float64  `json:"recevie_count"`      // 收货总数量
	ReceiveRemark     string   `json:"receive_remark"`     // 收货备注
	Remark            string   `json:"remark"`             // 备注
	CompanyId         int64    `json:"company_id"`         //
	Company           Contract `gorm:"ForeignKey:CompanyId" json:"company"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	ExpressNo string `json:"express_no"`
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

// 编辑仓库
func SendEdit(id int64, data interface{}) error {
	if err := db.Model(&Send{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
