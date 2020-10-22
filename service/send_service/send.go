package send_service

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"material/lib/setting"
	"material/lib/utils"
	"material/models"
	"strconv"
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
