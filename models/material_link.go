package models

import (
	"github.com/jinzhu/gorm"
	"material/lib/utils"
)

type MaterialLink struct {
	Model
	MaterialId int64    `json:"material_id"` // 材料单id
	Material   Material `gorm:"ForeignKey:MaterialId" json:"material"`

	ProductId   int64      `json:"product_id"`
	Product     Product    `gorm:"ForeignKey:ProductId" json:"product"`
	CompanyId   int64      `json:"company_id"`
	Count       float64    `json:"count"`
	ProjectId   int64      `json:"project_id"`
	Price       float64    `json:"price"`
	SupplyCycle int64      `json:"supply_cycle"` // 供货周期
	ReceiveTime utils.Time `json:"receive_time"` //接收时间
	Status      int64      `json:"status"`       //状态 0正常 -1逾期
	IsReceive   int64      `json:"is_receive"`   //是否接收

}

// 新增下料单链接带事物
func MaterialLinkAddT(ml *MaterialLink, t *gorm.DB) error {
	ml.Flag = 1
	if err := t.Create(&ml).Error; err != nil {
		return err
	}
	return nil
}

func MaterialLinkEditT(id int64, data interface{}, t *gorm.DB) error {
	if err := t.Model(&MaterialLink{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 获取下料单链接列表
func MaterialLinkGetAllLists(maps interface{}) ([]*MaterialLink, error) {
	var mls []*MaterialLink
	err := db.Model(&MaterialLink{}).Preload("Product").Preload("Material").Where(maps).Order("id desc").Find(&mls).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return mls, nil
}

// 获取逾期数据
//func MaterialLinkOverdueQueue() (*MaterialLink, error) {
//	//var ml MaterialLink
//	//err := db.Where("is_receive = 1 AND flag =1", id).Preload("Project").First(&project).Error
//	//if err != nil {
//	//	return &Material{}, err
//	//}
//	//return &project, nil
//}
