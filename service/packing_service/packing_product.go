package packing_service

import (
	"material/lib/utils"
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

	Status         int   `json:"status"`           //0已打包 1已发货 4已收货 已验收
	MaterialLinkId int64 `json:"material_link_id"` //下料单链接id

}

func ApiListsPP(maps string) ([]*models.PackingProduct, error) {
	return models.PackingProductGetLists(0, 9999, maps)
}

// 获取同步列表
func SyncGetListPP(packing_id int64) ([]map[string]interface{}, error) {
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["packing_id"] = packing_id
	list, err := ApiListsPP(utils.BuildWhere(maps))
	if err != nil {
		return nil, err
	}
	cb_list := make([]map[string]interface{}, len(list))
	for i, v := range list {
		cb_list[i] = map[string]interface{}{
			"id":            v.Id,
			"product_id":    v.Product.Id,
			"platform_key":  v.Product.PlatformKey,  //平台key
			"platform_id":   v.Product.PlatformId,   //平台id
			"platform_uid":  v.Product.PlatformUid,  //平台uid
			"material_name": v.Product.MaterialName, //材料名称
			"standard":      v.Product.Standard,     //规格
			"count":         v.Count,                //打包数量
			"return_count":  v.ReturnCount,          //退货
			"receive_count": v.ReceiveCount,         //接收数量
		}
	}
	return cb_list, nil
}
