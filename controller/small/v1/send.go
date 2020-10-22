package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/send_service"
	"time"
)

func SendInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	send, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "发货信息不存在")
		return
	}

	//查询对应打包信息
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["send_id"] = send.Id
	ps, err := models.PackingGetLists(0, 9999, utils.BuildWhere(maps))
	if err != nil {
		e.ApiErr(c, "导出列表失败")
		return
	}

	//检测当前用户是否有权限
	isAuth := false
	user_info, _ := c.Get("user_info")
	if _, err := models.ReceiverUsersCheckAuth(user_info.(models.Users).Cuid, send.ProjectId); err == nil {
		isAuth = true
	}

	e.ApiOk(c, "获取成功", struct {
		Send     models.Send       `json:"send"`
		Packings []*models.Packing `json:"packings"`
		IsAuth   bool              `json:"is_auth"`
	}{
		Send:     *send,
		Packings: ps,
		IsAuth:   isAuth,
	})
}

// 收货页面详情 - 只有收货人才有权限
func SendReceiverInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	send, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "发货信息不存在")
		return
	}

	//检测当前用户是否有权限
	user_info, _ := c.Get("user_info")
	ru, err := models.ReceiverUsersCheckAuth(user_info.(models.Users).Cuid, send.ProjectId)
	if err != nil {
		e.ApiErr(c, "非法请求")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Info  models.Send          `json:"info"`
		RUser models.ReceiverUsers `json:"r_user"`
	}{
		Info:  *send,
		RUser: *ru,
	})
}

// 退货
func SendReturn(c *gin.Context) {
	data := send_service.SendReturn{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//检测当前用户是否有权限
	user_info, _ := c.Get("user_info")
	data.Cuid = user_info.(models.Users).Cuid

	cb, err := send_service.AddReturn(data)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	e.ApiOk(c, "操作成功", cb)
}

// 确认收货
func SendReceiver(c *gin.Context) {
	data := send_service.SendAdd{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	if data.ReceiveRemark == "" {
		e.ApiErr(c, "请输入收货信息")
		return
	}
	if len(data.ReceiveAttachment) == 0 {
		e.ApiErr(c, "请上传收货附件")
		return
	}

	send, err := models.SendGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "发货信息不存在")
		return
	}

	if send.Status == 1 {
		e.ApiErr(c, "已收货，无法重复收货")
		return
	}

	//检测当前用户是否有权限
	user_info, _ := c.Get("user_info")
	ru, err := models.ReceiverUsersCheckAuth(user_info.(models.Users).Cuid, send.ProjectId)
	if err != nil {
		e.ApiErr(c, "非法请求")
		return
	}

	if len(send.Packing) == 0 {
		e.ApiErr(c, "打包信息有误")
		return
	}
	// 创建事务
	t := *models.NewTransaction()
	total_send_count := 0.0
	for _, v := range send.Packing {
		// 查询所有打包信息
		maps := utils.WhereToMap(nil)
		maps["flag"] = 1
		maps["packing_id"] = v.Id
		pp_list, err := models.PackingProductGetLists(0, 9999, utils.BuildWhere(maps))
		if err != nil {
			e.ApiErr(c, "打包列表有误")
			t.Rollback()
			return
		}
		tpl_packing_num := 0.0
		for _, pp := range pp_list {
			//操作全部数量
			tpl_pp_num := pp.Count - pp.ReturnCount
			if tpl_pp_num > 0 {
				// 编辑pp
				pp.ReceiveCount = tpl_pp_num
				tpl_packing_num = tpl_packing_num + tpl_pp_num
				models.PackingProductEditT(pp.Id, map[string]interface{}{
					"receive_count": tpl_pp_num,
					"is_receive":    1,
					"receive_time":  time.Now(),
				}, &t)
				// 编辑产品
				product, err := models.ProductGetInfoT(pp.Product.Id, &t)
				if err != nil {
					e.ApiErr(c, "产品获取有误")
					t.Rollback()
					return
				}
				models.ProductEditT(pp.Product.Id, map[string]interface{}{
					"receive_count": product.ReceiveCount + tpl_pp_num,
				}, &t)
			}
		}
		total_send_count = total_send_count + tpl_packing_num
		models.PackingEditT(v.Id, map[string]interface{}{
			"receive_count": tpl_packing_num,
			"status":        4,
		}, &t)
	}
	models.SendEditT(send.Id, map[string]interface{}{
		"receive_count":      total_send_count,
		"actual_receiver":    ru.Nickname,
		"receive_attachment": utils.JsonEncode(data.ReceiveAttachment),
		"receive_remark":     data.ReceiveRemark,
		"status":             1,
	}, &t)

	t.Commit()
	e.ApiOk(c, "操作成功", e.GetEmptyStruct())
}
