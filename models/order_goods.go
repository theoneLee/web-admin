package models

import (
	"gitee.com/muzipp/Distribution/pkg/logging"
	"github.com/jinzhu/gorm"
)

type OrderGoods struct {
	Model
	OrderId       int
	GoodsId       int
	GoodsName     string
	GoodsImg      string
	MemberPrice   float64
	Price         float64
	Number        int
	TotalPrice    float64
	SumPrice      float64
	Remark        string
	Specification string
}

func AddOrderGoods(data map[string]interface{}, tx *gorm.DB) (flag bool) {
	orderGoods := &OrderGoods{
		OrderId:       data["order_id"].(int),
		GoodsId:       data["goods_id"].(int),
		GoodsName:     data["goods_name"].(string),
		GoodsImg:      data["goods_img"].(string),
		MemberPrice:   data["member_price"].(float64),
		Price:         data["price"].(float64),
		Number:        data["number"].(int),
		TotalPrice:    data["total_price"].(float64),
		SumPrice:      data["sum_price"].(float64),
		Remark:        data["remark"].(string),
		Specification: data["specification"].(string),
	}
	err := tx.Create(orderGoods).Error

	if err != nil { //添加商品失败
		flag = true
		logging.Info("添加订单错误", err) //记录错误日志
		return
	}

	return  false
}
