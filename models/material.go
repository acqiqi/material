package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"material/lib/utils"
)

// 材料单
type Material struct {
	Model
	Name           string  `json:"name"`             // 材料单名称
	TotalAmount    float64 `json:"total_amount"`     // 下料总额（不含税）
	TotalTaxAmount float64 `json:"total_tax_amount"` // 下料总额（含税）
	DesignNo       string  `json:"design_no"`        // 设计订单号
	ApplyNo        string  `json:"apply_no"`         // 下料单号
	Remark         string  `json:"remark"`           // 备注
	CreateType     int     `json:"create_type"`      // 创建类型 新建,    采购计划生成
	Type           int     `json:"type"`             // 类型    内装材料,    幕墙面材,    幕墙辅材,    幕墙线材
	PlatformKey    string  `json:"platform_key"`     // 平台key
	PlatformUid    string  `json:"platform_uid"`     // 平台用户id
	PlatformId     string  `json:"platform_id"`      // 平台id  或者对照订单号

	ProjectId  int64    `json:"project_id"`
	Project    Project  `gorm:"ForeignKey:ProjectId" json:"project"`
	CompanyId  int64    `json:"company_id"`
	Company    Company  `gorm:"ForeignKey:CompanyId" json:"company"`
	ContractId int64    `json:"contract_id"` //合同
	Contract   Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	BeginTime utils.Time `json:"begin_time"` //同步開始時間
}

// 三方检测是否存在
func MaterialCheck(platform_id string, platform_key string, platform_uid string) (*Material, error) {
	log.Println(platform_id, platform_key, platform_uid)
	var material Material
	err := db.Where("platform_id = ? AND platform_key = ? AND platform_uid = ? AND flag =1",
		platform_id, platform_key, platform_uid).Preload("Company").First(&material).Error
	if err != nil {
		return &Material{}, err
	}
	return &material, nil
}

// 新增下料单带事物
func MaterialAddT(material *Material, t *gorm.DB) error {
	material.Flag = 1
	if err := t.Create(&material).Error; err != nil {
		return err
	}
	return nil
}
