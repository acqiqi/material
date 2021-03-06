package models

import "github.com/jinzhu/gorm"

type ProductLinkWire struct {
	Model
	SurfaceTreatment string `json:"surface_treatment"` // 表面处理
	Color            string `json:"color"`             // 颜色
	Area             string `json:"area"`              // 单面积(㎡/片)
	TotalCount       string `json:"total_count"`       // 总面积(㎡)
	ProductId        int64  `json:"product_id"`
}

func ProductLinkWireAddT(product *ProductLinkWire, t *gorm.DB) error {
	product.Flag = 1
	if err := t.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

type ProductLinkAuxiliary struct {
	Model
	MaterialStatus string `json:"material_status"` // 材质状态
	Weight         string `json:"weight"`          // 单重量
	TotalArea      string `json:"total_area"`      // 总面积(㎡)
	ProductId      int64  `json:"product_id"`
}

func ProductLinkAuxiliaryAddT(product *ProductLinkAuxiliary, t *gorm.DB) error {
	product.Flag = 1
	if err := t.Create(&product).Error; err != nil {
		return err
	}
	return nil
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

func ProductLinkSurfaceAddT(product *ProductLinkSurface, t *gorm.DB) error {
	product.Flag = 1
	if err := t.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

const (
	P内装材料 int = 0 + iota
	P幕墙面材
	P幕墙辅材
	P幕墙线材
)

var ProductType = map[int]string{
	P内装材料: "内装材料",
	P幕墙面材: "幕墙面材",
	P幕墙辅材: "幕墙辅材",
	P幕墙线材: "幕墙线材",
}

// 内装材料
type ProductLinkDefaultData MaterialLink

// 面材
type ProductLinkSurfaceData struct {
	MaterialLink
	ProductLinkSurface ProductLinkSurface `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId"  json:"link"`
}

// 辅材
type ProductLinkAuxiliaryData struct {
	MaterialLink
	ProductLinkAuxiliary ProductLinkAuxiliary `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"link"`
}

// 线材
type ProductLinkWireData struct {
	MaterialLink
	ProductLinkWire ProductLinkWire `gorm:"ForeignKey:ProductId;AssociationForeignKey:ProductId" json:"link"`
}

func ProductLinkGetInfo(product_type int, product_id string) (cb interface{}, err error) {
	switch product_type {
	case P内装材料:
		var product ProductLinkDefaultData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("Product").Table("vhake_material_link").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case P幕墙面材:
		var product ProductLinkSurfaceData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("Product").Preload("ProductLinkSurface").Table("vhake_material_link").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case P幕墙辅材:
		var product ProductLinkAuxiliaryData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("Product").Preload("ProductLinkAuxiliary").Table("vhake_material_link").First(&product).Error; err == nil {
			return product, nil
		}
		break
	case P幕墙线材:
		var product ProductLinkWireData
		if err := db.Where("id = ? AND flag =1", product_id).Preload("Product").Preload("ProductLinkWire").Table("vhake_material_link").First(&product).Error; err == nil {
			return product, nil
		}
		break
	}
	return nil, err
}

func ProductLinkGetLists(product_type, pageNum, pageSize int, maps interface{}) (cb interface{}, err error) {
	switch product_type {
	case P内装材料:
		var products []*ProductLinkDefaultData
		if err := db.Model(&Product{}).Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Preload("Product").Table("vhake_material_link").Find(&products).Error; err == nil {
			return products, nil
		} else {
			return []ProductLinkDefaultData{}, err
		}
		break
	case P幕墙面材:
		var products []*ProductLinkSurfaceData
		if err := db.Model(&Product{}).Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Preload("Product").Preload("ProductLinkSurface").Table("vhake_material_link").Find(&products).Error; err == nil {
			return products, nil
		} else {
			return []ProductLinkSurfaceData{}, err
		}
		break
	case P幕墙辅材:
		var products []*ProductLinkAuxiliaryData
		if err := db.Model(&Product{}).Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Preload("Product").Preload("ProductLinkAuxiliary").Table("vhake_material_link").Find(&products).Error; err == nil {
			return products, nil
		} else {
			return []ProductLinkAuxiliaryData{}, err
		}
		break
	case P幕墙线材:
		var products []*ProductLinkWireData
		if err := db.Model(&Product{}).Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Preload("Product").Preload("ProductLinkWire").Table("vhake_material_link").Find(&products).Error; err == nil {
			return products, nil
		} else {
			return []ProductLinkWireData{}, err
		}
		break
	}
	return nil, err
}
