package send_service

import (
	"errors"
	"material/lib/utils"
	"material/models"
)

type SendReturn struct {
	Id               int64    `json:"id"`
	SendId           int64    `json:"send_id"`
	PackingId        int64    `json:"packing_id"`
	ProductId        int64    `json:"product_id"`
	PackingProductId int64    `json:"packing_product_id"`
	ProjectId        int64    `json:"project_id"`
	Cuid             int64    `json:"cuid"`
	Count            float64  `json:"count"`       // 退货数量
	Reason           string   `json:"reason"`      // 退货原因
	Remark           string   `json:"remark"`      // 其他描述
	Attachment       []string `json:"attachment"`  // 照片
	RepairId         int64    `json:"repair_id"`   // 补货对应id 暂时不用
	Status           int      `json:"status"`      // 发起中 1接收
	UseCount         float64  `json:"use_count"`   // 接收数量
	ReturnType       string   `json:"return_type"` //退货类型

	CompanyId int64           `json:"company_id"`
	Company   models.Contract `gorm:"ForeignKey:CompanyId" json:"company"`
}

var ReturnType = []string{
	"已损坏", "质量有问题", "大小不符", "与原设计不符", "其他",
}

func init() {
}

// 新增退货
func AddReturn(data SendReturn) (models.SendReturn, error) {
	if data.Count <= 0 {
		return models.SendReturn{}, errors.New("请输入正确的退货数量")
	}
	if data.Reason == "" {
		return models.SendReturn{}, errors.New("请输入退货原因")
	}
	if data.ReturnType == "" {
		return models.SendReturn{}, errors.New("请选择退货类型")
	}
	if len(data.Attachment) > 0 {

	} else {
		return models.SendReturn{}, errors.New("请上传退货照片")
	}

	//其实只需要查ppid就可以直接拿到所有信息了
	pp, err := models.PackingProductGetInfo(data.PackingProductId)
	if err != nil {
		return models.SendReturn{}, errors.New("打包数据不存在")
	}

	//查询权限
	_, err = models.ReceiverUsersCheckAuth(data.Cuid, pp.ProjectId)
	if err != nil {
		return models.SendReturn{}, errors.New("非法请求")
	}

	// 查询是否有发货
	if pp.Packing.SendId == 0 {
		return models.SendReturn{}, errors.New("未发货无法退货")
	}
	send, err := models.SendGetInfo(pp.Packing.SendId)
	if err != nil {
		return models.SendReturn{}, errors.New("未发货无法退货 1")
	}
	if send.Status == 1 {
		return models.SendReturn{}, errors.New("当前状态无法退货")
	}

	if (pp.Count - pp.ReturnCount) < data.Count {
		return models.SendReturn{}, errors.New("退货数量已经超过可退数量")
	}

	model := models.SendReturn{
		SendId:           pp.Packing.SendId,
		PackingId:        pp.PackingId,
		ProductId:        pp.ProductId,
		PackingProductId: pp.Id,
		ProjectId:        pp.ProjectId,
		Cuid:             data.Cuid,
		Count:            data.Count,
		Reason:           data.Reason,
		Remark:           data.Remark,
		Attachment:       utils.JsonEncode(data.Attachment),
		RepairId:         data.RepairId,
		Status:           data.Status,
		UseCount:         data.UseCount,
		ReturnType:       data.ReturnType,
		CompanyId:        pp.CompanyId,
	}
	err = models.SendReturnAdd(&model)
	if err != nil {
		return models.SendReturn{}, errors.New("创建失败")
	}

	//修改当前数量
	pp.ReturnCount = pp.ReturnCount + data.Count
	models.PackingProductEdit(pp.Id, map[string]interface{}{
		"return_count": pp.ReturnCount,
	})
	pp.Packing.ReturnCount = pp.Packing.ReturnCount + data.Count
	models.PackingEdit(pp.Packing.Id, map[string]interface{}{
		"return_count": pp.Packing.ReturnCount,
	})
	//编辑产品
	product, _ := models.ProductGetInfo(pp.ProductId)
	product.ReturnCount = product.ReturnCount + data.Count
	models.ProductEdit(product.Id, map[string]interface{}{
		"return_count": product.ReturnCount,
	})
	// 编辑发货
	send.ReturnCount = send.ReturnCount + data.Count
	models.SendEdit(send.Id, map[string]interface{}{
		"return_count": send.ReturnCount,
	})

	cb, _ := models.SendReturnGetInfo(model.Id)
	return *cb, nil
}
