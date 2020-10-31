package models

type ProductLinkWire struct {
	Model
	MaterialStatus string `json:"material_status"` // 材质状态
	Weight         string `json:"weight"`          // 单重量
	TotalArea      string `json:"total_area"`      // 总面积(㎡)
	ProductId      int64  `json:"product_id"`
}

type ProductLinkAuxiliary struct {
	Model
	SurfaceTreatment string `json:"surface_treatment"` // 表面处理
	Color            string `json:"color"`             // 颜色
	Area             string `json:"area"`              // 单面积(㎡/片)
	TotalCount       string `json:"total_count"`       // 总面积(㎡)
	ProductId        int64  `json:"product_id"`
}

type ProductLinkSurface struct {
	Model
	W1               float64 `json:"w1"`
	W2               float64 `json:"w2"`
	W3               float64 `json:"w3"`
	W4               float64 `json:"w4"`
	W5               float64 `json:"w5"`
	W6               float64 `json:"w6"`
	W7               float64 `json:"w7"`
	W8               float64 `json:"w8"`
	W9               float64 `json:"w9"`
	H1               float64 `json:"h1"`
	H2               float64 `json:"h2"`
	H3               float64 `json:"h3"`
	H4               float64 `json:"h4"`
	H5               float64 `json:"h5"`
	H6               float64 `json:"h6"`
	H7               float64 `json:"h7"`
	H8               float64 `json:"h8"`
	H9               float64 `json:"h9"`
	L1               float64 `json:"l1"`
	L2               float64 `json:"l2"`
	L3               float64 `json:"l3"`
	L4               float64 `json:"l4"`
	L5               float64 `json:"l5"`
	L6               float64 `json:"l6"`
	L7               float64 `json:"l7"`
	L8               float64 `json:"l8"`
	L9               float64 `json:"l9"`
	WSize            int64   `json:"w_size"`            // 宽度数量
	HSize            int64   `json:"h_size"`            // 高度数量
	LSize            int64   `json:"l_size"`            // 长度数量
	SurfaceTreatment string  `json:"surface_treatment"` // 表面处理
	Color            string  `json:"color"`             // 颜色
	Area             string  `json:"area"`              // 单面积
	TotalCount       string  `json:"total_count"`       // 总面积
	ProductId        int64   `json:"product_id"`
}

const (
	内装材料 int = 0 + iota
	幕墙面材
	幕墙辅材
	幕墙线材
)

// 内装材料
type ProductLinkDefaultData Product

// 面材
type ProductLinkSurfaceData struct {
	Product
	ProductLinkSurface ProductLinkSurface `gorm:"ForeignKey:ProductId" json:"product_link_surface"`
}

// 辅材
type ProductLinkAuxiliaryData struct {
	Product
	ProductLinkAuxiliary ProductLinkAuxiliary `gorm:"ForeignKey:ProductId" json:"product_link_auxiliary"`
}

// 线材
type ProductLinkWireData struct {
	Product
	ProductLinkWire ProductLinkWire `gorm:"ForeignKey:ProductId" json:"product_link_wire"`
}

func ProductLinkGetInfo(product_type int, product_id string) (cb interface{}, err error) {
	switch product_type {
	case 内装材料:
		var product ProductLinkDefaultData
		if err := db.Where("id = ? AND flag =1", product_id).Table("vhake_product").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case 幕墙面材:
		var product ProductLinkSurfaceData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("ProductLinkSurface").Table("vhake_product").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case 幕墙辅材:
		var product ProductLinkAuxiliaryData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("ProductLinkAuxiliary").Table("vhake_product").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case 幕墙线材:
		var product ProductLinkWireData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("ProductLinkWire").Table("vhake_product").First(&product).Error; err == nil {
			return product, nil
		}
		break
	}
	return nil, err
}
