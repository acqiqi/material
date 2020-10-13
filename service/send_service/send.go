package send_service

import (
	uuid "github.com/satori/go.uuid"
	"material/lib/utils"
	"material/models"
)

// 发货表
type SendAdd struct {
	Id                int64   `json:"id"`
	SendNo            string  `json:"send_no"`            // 订单编号
	Count             float64 `json:"count"`              // 发货总数
	ActualReceiver    string  `json:"actual_receiver"`    // 签收人
	Address           string  `json:"address"`            // 收货地址
	RecevieAttachment string  `json:"recevie_attachment"` // 收货附件
	RecevieDate       string  `json:"recevie_date"`       // 收货时间
	RecevieCount      float64 `json:"recevie_count"`      // 收货总数量
	ReceiveRemark     string  `json:"receive_remark"`     // 收货备注
	Remark            string  `json:"remark"`             // 备注
	CompanyId         int64   `json:"company_id"`         //
	ProjectId         int64   `json:"project_id"`
	Express           string  `json:"express_no"`
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Send, error) {
	offset := (page - 1) * limit
	return models.SendGetLists(offset, limit, maps)
}

//新增
func Add(data SendAdd, links []*models.Packing) (*models.Send, error) {

	send_no := uuid.NewV4().String()

	//查询总打包量
	total_count := float64(0)
	for _, v := range links {
		total_count = total_count + v.Count
	}

	model := models.Send{
		SendNo:    send_no,
		Count:     total_count,
		Remark:    data.Remark,
		CompanyId: data.CompanyId,
		ProjectId: data.ProjectId,
		ExpressNo: data.Express,
	}
	//创建事务
	t := *models.NewTransaction()

	if err := models.SendAddT(&model, &t); err != nil {
		t.Rollback()
		return nil, err
	}

	//处理链接
	for _, v := range links {
		v.Status = 1
		v.SendId = model.Id
		if err := models.PackingEditT(v.Id, v, &t); err != nil {
			t.Rollback()
			return nil, err
		}
		//也对应查一下Product
		maps := utils.WhereToMap(nil)
		maps["flag"] = 1
		maps["packing_id"] = v.Id
		pp_list, err := models.PackingProductGetLists(0, 999, utils.BuildWhere(maps))
		if err != nil {
			t.Rollback()
			return nil, err
		}
		for _, v := range pp_list {
			v.Product.SendCount = v.Product.SendCount + v.Count
			status_save := make(map[string]interface{})
			status_save["status"] = 1
			models.PackingProductEditT(v.Id, status_save, &t)
			models.ProductEditT(v.Product.Id, v.Product, &t)
		}
	}
	t.Commit()
	return &model, nil
}
