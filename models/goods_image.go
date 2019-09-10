package models

import (
	"gitee.com/muzipp/Distribution/pkg/logging"
	"github.com/jinzhu/gorm"
)

type GoodsImg struct {
	Model
	GoodsId int
	Img     string
}

func AddGoodsImg(data map[string]interface{}, tx *gorm.DB) (flag bool) {
	goodsImg := &GoodsImg{
		GoodsId: data["goods_id"].(int),
		Img:     data["img"].(string),
	}
	err := tx.Create(goodsImg).Error

	if err != nil { //添加商品图片失败
		flag = true
		logging.Info("添加商品图片错误", err) //记录错误日志
		return
	}

	return false
}

//删除订单对应的图片
func DeleteGoodsImg(id int, tx *gorm.DB) (flag bool) {
	var goodsImg GoodsImg
	err := tx.Debug().Where("goods_id = ?", id).Delete(goodsImg).Error
	if err != nil {
		flag = true
	}
	return flag
}
