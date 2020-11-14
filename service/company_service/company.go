package company_service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/validation"
	uuid "github.com/satori/go.uuid"
	"log"
	"material/lib/app"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
)

type CompanyAdd struct {
	Cuid       int      `json:"cuid"`         // 注册人cuid
	Name       string   `json:"name" `        // 企业名称
	Mobile     string   `json:"mobile" `      // 企业手机号
	Tel        string   `json:"tel" `         // 电话号
	Address    string   `json:"address" `     // 企业地址
	Desc       string   `json:"desc"`         // 描述
	AuthPics   []string `json:"auth_pics"`    // 资质 多图
	VipLevel   int      `json:"vip_level"`    // 企业购买等级
	VipEndTime int      `json:"vip_end_time"` // 到期时间
	Status     int      `json:"status"`       // 状态 0 停业  1营业 -1停用
	DeletedAt  string   `json:"deleted_at"`
	CompanyKey string   `json:"company_key"` //企业Key
	BindState  int      `json:"bind_state"`  //是否绑定 0否 1是
	Ak         string   `json:"ak"`
	Sk         string   `json:"sk"`
}

//新增企业
func Add(data *CompanyAdd) (*models.Company, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.Cuid, "Cuid").Message("CUID不能为空！")
	valid.Required(data.Name, "Name").Message("请输入公司名称")
	valid.Mobile(data.Mobile, "Mobile").Message("请输入正确的手机号")
	valid.Required(data.Address, "Address").Message("请输入公司地址")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	log.Println("???")
	model := models.Company{}
	model.Cuid = data.Cuid
	model.Name = data.Name
	model.Mobile = data.Mobile
	model.Tel = data.Tel
	model.Address = data.Address
	model.Desc = data.Desc
	model.AuthPics = utils.JsonEncode(data.AuthPics)
	model.VipLevel = 0
	model.VipEndTime = 0
	model.Status = 1
	model.CompanyKey = models.CompanyGetKey()

	h := md5.New()
	h.Write([]byte(uuid.NewV4().String()))
	ak := hex.EncodeToString(h.Sum(nil))
	model.Ak = ak
	sk := uuid.NewV4().String()
	model.Sk = sk

	if err := models.CompanyAdd(&model); err != nil {
		return nil, err
	}
	//创建连接
	cu_model := models.CompanyUsers{
		Cuid:      data.Cuid,
		CompanyId: int(model.Id),
		IsMain:    1,
		RoleId:    0,
		RuleData:  utils.JsonEncode(e.GetEmptyStruct()),
		Status:    1,
	}
	if err := models.CompanyUsersAdd(&cu_model); err != nil {
		return nil, err
	}
	return &model, nil
}
