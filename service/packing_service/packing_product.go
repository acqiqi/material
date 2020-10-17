package packing_service

import (
	"material/models"
)

type PackingProductAdd struct {
	Id            int64
	PackingId     int64   `json:"packing_id"`     // 打包id
	CompanyId     int64   `json:"company_id"`     // 企业id
	OrderReturnid int64   `json:"order_returnid"` // 订单退货详情id
	ProductId     int64   `json:"product_id"`
	MaterialId    int64   `json:"material_id"`
	Count         float64 `json:"count"`         // 打包数量
	ReturnCount   float64 `json:"return_count"`  // 退货数量
	ReceiveCount  float64 `json:"receive_count"` //签收数量
	MaterialName  string  `json:"material_name"` // 产品名称

	ContractId int64           `json:"contract_id"` //合同
	Contract   models.Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	ProjectId int64          `json:"project_id"`
	Project   models.Project `gorm:"ForeignKey:ProjectId" json:"project"`

	DepositoryId int64             `json:"depository_id"`
	Depository   models.Depository `gorm:"ForeignKey:DepositoryId" json:"depository"`

	Status int `json:"status"` //0已打包 1已发货 4已收货 已验收
}

func ApiListsPP(maps string) ([]*models.PackingProduct, error) {
	return models.PackingProductGetLists(0, 9999, maps)
}
