package models

import "github.com/jinzhu/gorm"

type Packing struct {
	Model
	PackingName string   `json:"packing_name"` // 包装名称
	SerialNo    string   `json:"serial_no"`    // 包装编号
	Count       float64  `json:"count"`        // 产品总数
	ReturnCount float64  `json:"return_count"` // 包装下退货数量
	Remark      string   `json:"remark"`       // 描述
	CompanyId   int64    `json:"company_id"`
	Company     Contract `gorm:"ForeignKey:CompanyId" json:"company"`
	ProductId   int64    `json:"product_id"`
	MaterialId  int64    `json:"material_id"`
	Material    Material `gorm:"ForeignKey:MaterialId" json:"material"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	ContractId int64    `json:"contract_id"` //合同
	Contract   Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	DepositoryId int64      `json:"depository_id"`
	Depository   Depository `gorm:"ForeignKey:DepositoryId" json:"depository"`

	Status int `json:"status"` //0已打包 1已发货 4已收货 已验收

	SendId int64 `json:"send_id"` //发货id
}

// 产品类型详情
func PackingGetInfo(id int64) (*Packing, error) {
	var pc Packing
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").First(&pc).Error
	if err != nil {
		return &Packing{}, err
	}
	return &pc, nil
}

// 产品类型列表
func PackingGetLists(pageNum int, pageSize int, maps interface{}) ([]*Packing, error) {
	var pc []*Packing
	err := db.Model(&Packing{}).Preload("Material").Preload("Project").Preload("Depository").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}

// 查询产品类型列表数量
func PackingGetListsCount(maps interface{}) int {
	var packing []*Packing
	count := 0
	db.Preload("Company").Where(maps).Find(&packing).Count(&count)
	return count
}

// 新增打包 T
func PackingAddT(packing *Packing, t *gorm.DB) error {
	packing.Flag = 1
	if err := t.Create(&packing).Error; err != nil {
		return err
	}
	return nil
}

// 新增打包类型
func PackingAdd(pc *Packing) error {
	pc.Flag = 1
	if err := db.Create(&pc).Error; err != nil {
		return err
	}
	return nil
}

// 编辑打包类型
func PackingEdit(id int64, data interface{}) error {
	if err := db.Model(&Packing{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func PackingEditT(id int64, data interface{}, t *gorm.DB) error {
	if err := t.Model(&Packing{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 获取打包select
func PackingGetSelect(maps string) ([]*Packing, error) {
	var packing []*Packing
	err := db.Where(maps).Order("id asc").Preload("Contract").Preload("Project").Find(&packing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return packing, nil
}
