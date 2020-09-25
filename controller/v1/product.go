package v1

//
//// 创建产品
//func ProductCreate(c *gin.Context) {
//	data := product_service.ProductAdd{}
//	if err := c.BindJSON(&data); err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	user_info, _ := c.Get("user_info")
//	data.Cuid = int(user_info.(models.Users).Id)
//	company, _ := c.Get("company")
//	data.CompanyId = company.(models.CompanyUsers).Company.Id
//	//检测项目是否存在
//
//	cb, err := project_service.Add(&data)
//	if err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//
//	p, _ := models.ProjectGetInfo(cb.Id)
//
//	e.ApiOk(c, "创建成功", p)
//}
//
//// 编辑产品
//func ProductEdit(c *gin.Context) {
//	data := project_service.ProjectAdd{}
//	if err := c.BindJSON(&data); err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	// 检测是否有项目
//	project, err := models.ProjectGetInfo(data.Id)
//	if err != nil {
//		e.ApiErr(c, "项目不存在")
//		return
//	}
//	company, _ := c.Get("company")
//	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
//		e.ApiErr(c, "非法请求")
//		return
//	}
//
//	cb, err := project_service.Edit(&data)
//	if err != nil {
//		e.ApiErr(c, err.Error())
//		return
//	}
//	p, _ := models.ProjectGetInfo(cb.Id)
//	e.ApiOk(c, "编辑成功", p)
//}
