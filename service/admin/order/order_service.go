package order

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
)

type Order struct {
	MemberId   int
	Status     int
	Number     string
	Remark     string
	StartTime  string
	EndTime    string
	OrderField string
	OrderSort  int
	Offset     int
	Limit      int
}

//商品列表
func (o *Order) ListOrders() (goods []models.Goods, err e.SelfError) {
	fields := "*"
	goods, goodsErr := models.ListOrders(o.Offset, o.Limit, o.getMaps(), fields)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return

}

//商品数量
//func (o *Order) CountOrders() (count int, err e.SelfError) {
//	count, goodsErr := models.CountOrders(o.getMaps())
//	if goodsErr {
//		err.Code = e.ERROR_SQL_FAIL
//	}
//
//	return
//}

//封装搜索条件
//func (o *Order) getMaps() map[string]interface{} {
//	maps := make(map[string]interface{})
//	maps["order.delete_at"] = 0
//
//	if o.MemberId != 0 {//拼接会员ID
//		maps["order.member_id"] = o.MemberId
//	}
//
//	if o.Number != "" {//拼接订单编号
//		maps["order.number"] = o.Number
//	}
//
//	if o.Remark != "" {//拼接备注
//		maps["order.remark"] = o.Remark
//	}
//
//	if o.StartTime!="" {
//		maps["create_at"] =
//	}
//
//	if o.EndTime!="" {
//
//	}
//	return maps
//}
