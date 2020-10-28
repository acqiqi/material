package models

import "github.com/jinzhu/gorm"

type MaterialLink struct {
	Model
	MaterialId int64   `json:"material_id"` // 材料单id
	ProductId  int64   `json:"product_id"`
	Product    Product `gorm:"ForeignKey:ProductId" json:"product"`
	CompanyId  int64   `json:"company_id"`
	Count      float64 `json:"count"`
	ProjectId  int64   `json:"project_id"`
	Price      float64 `json:"price"`
}

// 新增下料单链接带事物
func MaterialLinkAddT(ml *MaterialLink, t *gorm.DB) error {
	ml.Flag = 1
	if err := t.Create(&ml).Error; err != nil {
		return err
	}
	return nil
}

// 获取下料单链接列表
func MaterialLinkGetAllLists(maps interface{}) ([]*MaterialLink, error) {
	var mls []*MaterialLink
	err := db.Model(&MaterialLink{}).Preload("Product").Where(maps).Order("id desc").Find(&mls).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return mls, nil
}
