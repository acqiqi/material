package models

import "github.com/jinzhu/gorm"

type PackingProduct struct {
	Model
	PackingId     int64   `orm:"packing_id"`     // 打包id
	CompanyId     int64   `orm:"company_id"`     // 企业id
	OrderReturnid int64   `orm:"order_returnid"` // 订单退货详情id
	ProductId     int64   `orm:"product_id"`
	Product       Product `gorm:"ForeignKey:ProductId" json:"product"`

	MaterialId   int64   `orm:"material_id"`
	Count        float64 `orm:"count"`         // 打包数量
	ReturnCount  float64 `orm:"return_count"`  // 退货数量
	MaterialName string  `orm:"material_name"` // 产品名称

	ContractId int64    `json:"contract_id"` //合同
	Contract   Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	DepositoryId int64      `json:"depository_id"`
	Depository   Depository `gorm:"ForeignKey:DepositoryId" json:"depository"`

	Status int `json:"status"` //0已打包 1已发货 4已收货 已验收
}

// 新增打包 T
func PackingProductAddT(pp *PackingProduct, t *gorm.DB) error {
	pp.Flag = 1
	if err := t.Create(&pp).Error; err != nil {
		return err
	}
	return nil
}

// 获取打包产品列表
func PackingProductGetLists(pageNum int, pageSize int, maps interface{}) ([]*PackingProduct, error) {
	var pp []*PackingProduct
	err := db.Model(&PackingProduct{}).Where(maps).Preload("Product").Offset(pageNum).Limit(pageSize).Order("id desc").Find(&pp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pp, nil
}

// 编辑打包产品
func PackingProductEditT(id int64, data interface{}, t *gorm.DB) error {
	if err := t.Model(&PackingProduct{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
