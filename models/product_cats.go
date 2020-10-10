package models

import "github.com/jinzhu/gorm"

// 材料分类表
type ProductCats struct {
	Id         int64  `json:"id"`
	CatName    string `json:"cat_name"`    // 大类名
	ModelName  string `json:"model_name"`  // 数据模型名称
	ModelTable string `json:"model_table"` // 数据表名 全名
	Desc       string `json:"desc"`
	IsShow     int    `json:"is_show"`
}

// 获取材料分类
func ProductCatsGetSelect(maps string) ([]*ProductCats, error) {
	var pc []*ProductCats
	err := db.Where(maps).Order("id asc").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}
