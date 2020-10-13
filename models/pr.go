package models

import (
	"github.com/jinzhu/gorm"
	"material/lib/utils"
)

type Pr struct {
	Model

	CompanyId int64   `json:"company_id"`
	Company   Company `gorm:"ForeignKey:CompanyId" json:"company"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	Price        float64    `json:"price"`     // 实际金额
	TplPrice     float64    `json:"tpl_price"` // 模板金额 输入的
	Type         int        `json:"type"`      // 0进度款 1预付款 2结算款 3质保金
	Desc         string     `json:"desc"`      // 描述
	Count        float64    `json:"count"`     // 数量
	Cuid         int64      `json:"cuid"`
	Status       int        `json:"status"` // 0正常 1接收 2审批通过 -1未通过
	PlatformKey  string     `json:"platform_key"`
	PlatformUid  string     `json:"platform_uid"`
	PlatformId   string     `json:"platform_id"`
	IsPlatform   int        `json:"is_platform"`   // 是否平台审核
	PlatformMsg  string     `json:"platform_msg"`  // 平台审核消息
	PrNo         string     `json:"pr_no"`         // 请款唯一编号
	ApprovalTime utils.Time `json:"approval_time"` // 审批时间
}

// 请款详情
func PrGetInfo(id int64) (*Pr, error) {
	var pr Pr
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").First(&pr).Error
	if err != nil {
		return &Pr{}, err
	}
	return &pr, nil
}

// 请款列表
func PrGetLists(pageNum int, pageSize int, maps interface{}) ([]*Pr, error) {
	var pr []*Pr
	err := db.Model(&Pr{}).Preload("Company").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pr, nil
}

// 查询请款列表数量
func PrCount(maps interface{}) int {
	var pr []*Pr
	count := 0
	db.Preload("Company").Where(maps).Find(&pr).Count(&count)
	return count
}

// 新增请款
func PrAdd(pr *Pr) error {
	pr.Flag = 1
	if err := db.Create(&pr).Error; err != nil {
		return err
	}
	return nil
}

// 编辑请款
func PrEdit(id int64, data interface{}) error {
	if err := db.Model(&Pr{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func PrGetSelect(maps string) ([]*Pr, error) {
	var pr []*Pr
	err := db.Where(maps).Order("id desc").Find(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pr, nil
}
