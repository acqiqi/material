package models

import (
	"github.com/jinzhu/gorm"
	"material/lib/utils"
)

type PackingProduct struct {
	Model
	PackingId int64   `json:"packing_id"` // 打包id
	Packing   Packing `gorm:"ForeignKey:PackingId" json:"packing"`

	CompanyId     int64   `json:"company_id"`     // 企业id
	OrderReturnid int64   `json:"order_returnid"` // 订单退货详情id
	ProductId     int64   `json:"product_id"`
	Product       Product `gorm:"ForeignKey:ProductId" json:"product"`

	MaterialId   int64   `json:"material_id"`
	Count        float64 `json:"count"`         // 打包数量
	ReturnCount  float64 `json:"return_count"`  // 退货数量
	ReceiveCount float64 `json:"receive_count"` //签收数量
	MaterialName string  `json:"material_name"` // 产品名称

	ContractId int64    `json:"contract_id"` //合同
	Contract   Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	ProjectId int64   `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectId" json:"project"`

	DepositoryId int64      `json:"depository_id"`
	Depository   Depository `gorm:"ForeignKey:DepositoryId" json:"depository"`

	Status      int        `json:"status"`       //0已打包 1已发货 4已收货 已验收
	ReceiveTime utils.Time `json:"receive_time"` //签收时间
	IsReceive   int        `json:"is_receive"`   //是否签收
}

// 获取详情
func PackingProductGetInfo(id int64) (*PackingProduct, error) {
	var d PackingProduct
	err := db.Where("id = ? AND flag = 1", id).Preload("Project").Preload("Product").Preload("Packing").First(&d).Error
	if err != nil {
		return &PackingProduct{}, err
	}
	return &d, nil
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

func PackingProductEdit(id int64, data interface{}) error {
	if err := db.Model(&PackingProduct{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
