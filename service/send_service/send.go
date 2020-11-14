package send_service

import (
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"material/lib/e"
	"material/lib/setting"
	"material/lib/utils"
	"material/models"
	"material/service/packing_service"
	"strconv"
	"strings"
)

// 发货表
type SendAdd struct {
	Id                int64    `json:"id"`
	SendNo            string   `json:"send_no"`            // 订单编号
	Count             float64  `json:"count"`              // 发货总数
	ActualReceiver    string   `json:"actual_receiver"`    // 签收人
	Address           string   `json:"address"`            // 收货地址
	ReceiveAttachment []string `json:"receive_attachment"` // 收货附件
	ReceiveDate       string   `json:"receive_date"`       // 收货时间
	ReceiveCount      float64  `json:"receive_count"`      // 收货总数量
	ReceiveRemark     string   `json:"receive_remark"`     // 收货备注
	Remark            string   `json:"remark"`             // 备注
	CompanyId         int64    `json:"company_id"`         //
	ProjectId         int64    `json:"project_id"`
	Express           string   `json:"express_no"`

	Status int `json:"status"` //0未签收 1签收

	IsSync      int    `json:"is_sync"`      // 是否同步 如果platform存在就需要同步
	PlatformKey string `json:"platform_key"` // 平台key

	Title string `json:"title"`
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
		SendNo:      send_no,
		Count:       total_count,
		Remark:      data.Remark,
		CompanyId:   data.CompanyId,
		ProjectId:   data.ProjectId,
		ExpressNo:   data.Express,
		PlatformKey: data.PlatformKey,
		IsSync:      0,
		Title:       data.Title,
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
			status_save := make(map[string]interface{})
			status_save["status"] = 1
			models.PackingProductEditT(v.Id, status_save, &t)
			product, _ := models.ProductGetInfoT(v.Product.Id, &t)
			product.SendCount = product.SendCount + v.Count
			models.ProductEditT(v.Product.Id, map[string]interface{}{
				"send_count": product.SendCount,
			}, &t)
		}
	}
	t.Commit()

	s_data, _ := models.SendGetInfoT(model.Id, &t)
	//log.Println(err)
	SyncCallback(*s_data, false)
	return &model, nil
}

// 生成二维码
func QrcodeBuild(send models.Send) (string, error) {
	c := gofpdf.InitType{
		Size: gofpdf.SizeType{
			Wd: 160,
			Ht: 80,
		},
	}
	// 获取二维码
	wechatUtils := new(utils.WechatUtils)
	wechatUtils.SmallQrcodeData.Page = "pages/send/send"
	wechatUtils.SmallQrcodeData.Width = 430
	wechatUtils.SmallQrcodeData.Scene = strconv.FormatInt(send.Id, 10)
	wechatUtils.Init(setting.WechatSetting.SmallAppID, setting.WechatSetting.AppSecret)
	wechatUtils.GetAccessToken()
	b, err := wechatUtils.GetSmallQrcode()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	ioutil.WriteFile("static/qrcode/send-"+strconv.FormatInt(send.Id, 10)+".jpg", b, 0755)

	pdf := gofpdf.NewCustom(&c)
	pdf.AddUTF8Font("u8font", "", "ttf/pdf.ttf")
	//pdf.SetCellMargin(0)
	//pdf.SetLeftMargin(0)
	//pdf.SetTopMargin(0)
	//pdf.SetTopMargin(0)
	pdf.AddPage()
	pdf.Image("static/qrcode/send-"+strconv.FormatInt(send.Id, 10)+".jpg", 10, 10, 60, 60, false, "", 0, "")
	//pdf.Text(60, 10, "所属项目：哈哈哈哈哈哈哈哈哈哈")
	//pdf.SetY(pdf.SetY(0)            //先要设置 Y，然后再设置 X。否则，会导致 X 失效
	pdf.SetY(10) //水平居中的算法0)            //先要设置 Y，然后再设置 X。否则，会导致 X 失效
	pdf.SetX(70) //水平居中的算法
	pdf.SetFont("u8font", "", 20)
	pdf.MultiCell(90, 9,
		fmt.Sprintf("项目名称：%s \n\r发货号：%s \n\r",
			send.Project.ProjectName,
			send.SendNo), "", "Left", false)
	fileStr := setting.AppSetting.StaticUrl + "pdf/send-" + strconv.FormatInt(send.Id, 10) + ".pdf"
	err = pdf.OutputFileAndClose("static/pdf/send-" + strconv.FormatInt(send.Id, 10) + ".pdf")
	if err != nil {
		return "", err
	}
	return fileStr, nil
}

// 同步发货
func SyncCallback(send models.Send, is_r bool) error {
	//查询平台是否存在
	platform, err := models.PlatformGetInfoOrKey(send.PlatformKey)
	if err != nil {
		return errors.New("平台有误")
	}

	//创建Send信息
	send_data := make(map[string]interface{})
	send_data["id"] = send.Id
	send_data["send_no"] = send.SendNo                       // 订单编号
	send_data["count"] = send.Count                          // 发货总数
	send_data["remark"] = send.Remark                        // 发货备注
	send_data["express_no"] = send.ExpressNo                 // 快递单号
	send_data["return_count"] = send.ReturnCount             // 退货数量
	send_data["receive_count"] = send.ReceiveCount           // 接收数量
	send_data["actual_receiver"] = send.ActualReceiver       // 签收人
	send_data["receive_remark"] = send.ReceiveRemark         //收货备注
	send_data["receive_attachment"] = send.ReceiveAttachment //收件人附件
	send_data["receive_mobile"] = send.ReceiveMobile         //手机号
	var ar []string
	if err := utils.JsonDecode(send.ReceiveAttachment, &ar); err == nil {
		send_data["receive_attachment"] = strings.Join(ar, ",") // 附件
	}

	//处理打包信息
	var product_list []map[string]interface{}
	//for i, v := range send.Packing {
	//	packing_data[i] = map[string]interface{}{
	//		"id":           v.Id,          //包装id
	//		"packing_name": v.PackingName, //包装名
	//		"serial_no":    v.SerialNo,    //包装编号
	//		"count":        v.Count,       //产品总数
	//		"remark":       v.Remark,      //描述
	//		"product_list": packing_service.SyncGetListPP(v.Id),
	//	}
	//}
	for _, v := range send.Packing {
		p_list, err := packing_service.SyncGetListPP(v.Id)
		if err == nil {
			for _, val := range p_list {
				product_list = append(product_list, val)
			}
		}
	}
	send_data["product_list"] = utils.JsonEncode(product_list)

	action := ""
	if is_r {
		action = e.PLATFORT_ACTION_SEND_RECEIVER
	} else {
		action = e.PLATFORT_ACTION_SEND
	}

	callback := e.HttpCallbackData{
		Code:        0,
		Msg:         "Send Success",
		Action:      action,
		CallbackUrl: platform.MessageCallbackUrl,
		Data:        send_data,
	}
	c_data := new(e.HttpCallbackData)

	str := utils.JsonEncode(callback)
	log.Println(str)
	if err := callback.RequestCallback(&c_data); err != nil {
		return err
	}
	if c_data.Code == 0 {
		//修改同步状态
		up_send := map[string]interface{}{
			"is_sync": 1,
		}
		models.SendEdit(send.Id, up_send)
	}
	return nil
}
