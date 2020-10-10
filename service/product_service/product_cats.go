package product_service

import (
	"material/lib/utils"
	"material/models"
)

// 获取Api列表
func CatsSelect() ([]*models.ProductCats, error) {
	maps := utils.WhereToMap(nil)
	maps["is_show"] = 1
	maps["flag"] = 1
	return models.ProductCatsGetSelect(utils.BuildWhere(maps))
}
