package models

import (
	"gitee.com/muzipp/Distribution/pkg/logging"
	"github.com/jinzhu/gorm"
)

type Goods struct {
	Model
	Name          string
	Stock         int
	Price         float64
	Integral      int
	Specification string
	Remark        string
	Status        int
	Description   string
	StatusDesc    string   `gorm:"-"`
	Images        []string `gorm:"-"`
	Img           string   `gorm:"-"`
}

func AddGoods(data map[string]interface{}, tx *gorm.DB) (id int, flag bool) {
	goods := &Goods{
		Name:          data["name"].(string),
		Stock:         data["stock"].(int),
		Price:         data["price"].(float64),
		Integral:      data["integral"].(int),
		Specification: data["specification"].(string),
		Remark:        data["remark"].(string),
		Status:        data["status"].(int),
		Description:   data["description"].(string),
	}
	err := tx.Create(goods).Error

	if err != nil { //添加商品失败
		flag = true
		logging.Info("添加商品错误", err) //记录错误日志
		return
	}

	return goods.ID, false
}

func ListGoods(pageNum int, pageSize int, maps interface{}, fields string) (goods []Goods, flag bool) {

	err := Db.Table("goods").
		Where(maps).
		Joins("left join `goods_img` as gi on gi.goods_id = goods.id").
		Offset(pageNum).
		Limit(pageSize).
		Select(fields).
		Scan(&goods).Error
	if err != nil {
		flag = true
		logging.Info("商品列表错误", err) //记录错误日志
		return
	}
	return
}

func CountGoods(maps interface{}) (count int, flag bool) {
	err := Db.Model(&Goods{}).Where(maps).Count(&count).Error
	if err != nil {
		flag = true
		logging.Info("商品数量错误", err) //记录错误日志
		return
	}
	return
}

//商品详情
func DetailGoods(id int, fields string) (*Goods, bool) {
	var goods Goods
	var flag bool
	err := Db.Table("goods").
		Where("goods.id = ? AND goods.delete_at = ? ", id, 0).
		Joins("left join `goods_img` as gi on gi.goods_id = goods.id").
		Select(fields).
		Find(&goods).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		flag = true
		logging.Info("商品详情错误", err) //记录错误日志
		return &goods, flag
	}

	if gorm.IsRecordNotFoundError(err) { //查询结果不存在的情况
		flag = true
		return nil, flag
	}
	return &goods, flag
}

func DeleteGoods(id int, tx *gorm.DB) (flag bool) {
	var goods Goods
	err := tx.Where("id = ?", id).Delete(&goods).Error
	if err != nil {
		flag = true
	}
	return flag
}

func EditGoods(data map[string]interface{}, tx *gorm.DB, id int) (flag bool) {
	goods := &Goods{
		Name:          data["name"].(string),
		Stock:         data["stock"].(int),
		Price:         data["price"].(float64),
		Integral:      data["integral"].(int),
		Specification: data["specification"].(string),
		Remark:        data["remark"].(string),
		Status:        data["status"].(int),
		Description:   data["description"].(string),
	}
	err := tx.Table("goods").Where("id = ? ", id).Update(goods).Error

	if err != nil { //添加商品失败
		flag = true
		logging.Info("编辑商品错误", err) //记录错误日志
		return
	}

	return
}
