package models

import (
	"github.com/jinzhu/gorm"
	"material/lib/utils"
)

type Product struct {
	Model
	MaterialName       string  `json:"material_name"`       // 材料名称
	BlankingAttachment string  `json:"blanking_attachment"` // 下料附件信息(与码里公装关联)
	Attachment         string  `json:"attachment"`          // 附件
	InstallMap         string  `json:"install_map"`         // 安装示意图
	Price              float64 `json:"price"`               // 价格
	Count              float64 `json:"count"`               // 数量
	ContractCount      float64 `json:"contract_count"`      // 与供应商签的合同数量(来源码里公装)
	PackCount          float64 `json:"pack_count"`          // 打包数量
	SendCount          float64 `json:"send_count"`          // 发货数量
	ReturnCount        float64 `json:"return_count"`        // 退货数量
	ReceiveCount       float64 `json:"receive_count"`       // 签收数量
	Unit               string  `json:"unit"`                // 单位
	ProjectId          int64   `json:"project_id"`
	Project            Project `gorm:"ForeignKey:ProjectId" json:"project"`

	ProjectName         string     `json:"project_name"`
	ReplenishmentFlag   int        `json:"replenishment_flag"`    // 是否补货产品
	ProductSubFlag      int        `json:"product_sub_flag"`      // 是否有子部件
	ConfigData          string     `json:"config_data"`           // 自定义字段信息
	AppendAttachment    string     `json:"append_attachment"`     // 附加的资源库信息
	ProjectMaterialId   int        `json:"project_material_id"`   // 码里公装对应下料材料id
	AdminMaterialInfoId int        `json:"admin_materialInfo_id"` // 码里公装对应合同材料id，统计累计数量需要
	ProjectAdditional   string     `json:"project_additional"`    // 项目补充信息
	Remark              string     `json:"remark"`                // 备注
	Length              float64    `json:"length"`                // 长
	Width               float64    `json:"width"`                 // 宽
	Height              float64    `json:"height"`                // 厚
	Location            string     `json:"location"`              // 安装位置
	Standard            string     `json:"standard"`              // 规格
	ArriveDate          utils.Time `json:"arrive_date"`           // 到货时间
	Cuid                int        `json:"cuid"`
	CompanyId           int64      `json:"company_id"`
	Company             Company    `gorm:"ForeignKey:CompanyId" json:"company"`
	SupplyCycle         int        `json:"supply_cycle"` // 供货周期
	MaterialId          int64      `json:"material_id"`  // 材料单id
	PlatformKey         string     `json:"platform_key"` //平台key
	PlatformUid         string     `json:"platform_uid"` //平台uid
	PlatformId          string     `json:"platform_id"`  //平台id
}

// 新增单个产品 带事物
func ProductAddT(product *Product, t *gorm.DB) error {
	product.Flag = 1
	if err := t.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

// 新增多个产品 带事物
//func ProductAddAllT(products []*Product, t gorm.DB) error {
//	for i, _ := range products {
//		products[i].Flag = 1
//	}
//	if err := t.Create(&products).Error; err != nil {
//		return err
//	}
//	return nil
//}

// 新增产品
func ProductAdd(product *Product) error {
	product.Flag = 1
	if err := db.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

// 编辑产品
func ProductEdit(id int64, data interface{}) error {
	if err := db.Model(&Product{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 获取产品详情
func ProductGetInfo(id int64) (*Product, error) {
	var project Product
	err := db.Where("id = ? AND flag =1", id).Preload("Company").First(&project).Error
	if err != nil {
		return &Product{}, err
	}
	return &project, nil
}
