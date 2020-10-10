package packing_service

import (
	"material/models"
)

type PackingProductAdd struct {
	Id            int64
	PackingId     int64   `orm:"packing_id"`     // 打包id
	CompanyId     int64   `orm:"company_id"`     // 企业id
	OrderReturnid int64   `orm:"order_returnid"` // 订单退货详情id
	ProductId     int64   `orm:"product_id"`
	MaterialId    int64   `orm:"material_id"`
	Count         float64 `orm:"count"`         // 打包数量
	ReturnCount   float64 `orm:"return_count"`  // 退货数量
	MaterialName  string  `orm:"material_name"` // 产品名称

	ContractId int64           `json:"contract_id"` //合同
	Contract   models.Contract `gorm:"ForeignKey:ContractId" json:"contract"`
}
