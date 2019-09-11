package goods

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	seleGoods "gitee.com/muzipp/Distribution/pkg/goods"
	"strings"
)

type Goods struct {
	Id            int
	Name          string
	Stock         int
	Price         float64
	Remark        string
	Image         string
	Integral      int
	Description   string
	Specification string
	Status        int
	Images        []string
	Offset        int
	Limit         int
}

func (g *Goods) AddGoods() (err e.SelfError) {
	data := make(map[string]interface{})
	data["name"] = g.Name
	data["remark"] = g.Remark
	data["price"] = g.Price
	data["stock"] = g.Stock
	data["status"] = g.Status
	data["description"] = g.Description
	data["integral"] = g.Integral
	data["status"] = g.Status
	data["specification"] = g.Specification

	//开启事务
	tx := models.Db.Begin()
	id, res := models.AddGoods(data, tx) //添加商品

	if res { //添加商品失败
		tx.Rollback() //事务回滚
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	imageUrl := strings.Split(g.Image, ",")

	//遍历图片
	dataImage := make(map[string]interface{})
	for _, value := range imageUrl {
		dataImage["goods_id"] = id
		dataImage["img"] = value
		res := models.AddGoodsImg(dataImage, tx)
		if res { //添加商品图片失败
			tx.Rollback() //事务回滚
			err.Code = e.ERROR_SQL_FAIL
			return
		}
	}

	tx.Commit() //事务提交
	return
}

//商品列表
func (g *Goods) ListGoods() (goods []models.Goods, err e.SelfError) {
	fields := "*"
	goods, goodsErr := models.ListGoods(g.Offset, g.Limit, g.getMaps(), fields)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	//切割字符串图片
	for key, value := range goods {
		goods[key].Images = strings.Fields(value.Img)
		goods[key].StatusDesc = seleGoods.GetStatus(value.Status)
	}

	return

}

//商品数量
func (g *Goods) CountGoods() (count int, err e.SelfError) {
	count, goodsErr := models.CountGoods(g.getMaps())
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

//获取文章（redis不存在读取数据库）
func (g *Goods) DetailGoods() (goods *models.Goods, err e.SelfError) {
	fields := "*"
	goods, goodsErr := models.DetailGoods(g.Id, fields)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	if goods != nil {
		goods.Images = strings.Fields(goods.Img)
		goods.StatusDesc = seleGoods.GetStatus(goods.Status)
	}
	return
}

//判断商品是否存在
func (g *Goods) DeleteGoods() (flag bool) {

	tx := models.Db.Begin()

	goodsErr := models.DeleteGoods(g.Id, tx)
	goodsImgErr := models.DeleteGoodsImg(g.Id, tx)

	if goodsErr || goodsImgErr {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return goodsErr || goodsImgErr
}

func (g *Goods) EditGoods(id int) (err e.SelfError) {
	data := make(map[string]interface{})
	data["name"] = g.Name
	data["remark"] = g.Remark
	data["price"] = g.Price
	data["stock"] = g.Stock
	data["status"] = g.Status
	data["description"] = g.Description
	data["integral"] = g.Integral
	data["status"] = g.Status
	data["specification"] = g.Specification

	//开启事务
	tx := models.Db.Begin()
	res := models.EditGoods(data, tx, id) //添加商品

	if res { //添加商品失败
		tx.Rollback() //事务回滚
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	deleteImgErr := models.DeleteGoodsImg(id, tx)
	if deleteImgErr { //添加商品失败
		tx.Rollback() //事务回滚
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	imageUrl := strings.Split(g.Image, ",")

	//遍历图片
	dataImage := make(map[string]interface{})
	for _, value := range imageUrl {
		dataImage["goods_id"] = id
		dataImage["img"] = value
		res := models.AddGoodsImg(dataImage, tx)
		if res { //添加商品图片失败
			tx.Rollback() //事务回滚
			err.Code = e.ERROR_SQL_FAIL
			return
		}
	}

	tx.Commit() //事务提交
	return
}

//封装搜索条件
func (g *Goods) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_at"] = 0
	return maps
}
