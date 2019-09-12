package order

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/order"
	"time"
)

type Order struct {
	models.Model
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

//订单列表
func (o *Order) ListOrders() (orders []models.Order, err e.SelfError) {
	fields := "order.id,m.name as member_name, order.number,order.create_at, order.status," +
		"order.reference_price,order.actual_price,order.discount,order.integral," +
		"order.remark,m1.name as team_name"
	orders, ordersErr := models.ListOrders(o.Offset, o.Limit, o.getMaps(), fields, o.Remark, o.StartTime, o.EndTime, o.OrderField, o.OrderSort)
	if ordersErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	//处理显示文案部分
	for key, value := range orders {
		orders[key].StatusDesc = order.GetStatus(value.Status)
		timeNow := time.Unix(int64(value.CreateAt), 0)
		orders[key].CreateTime = timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32

	}

	return

}

//订单数量
func (o *Order) CountOrders() (count int, err e.SelfError) {
	count, goodsErr := models.CountOrders(o.getMaps(), o.Remark, o.StartTime, o.EndTime)
	if goodsErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

//封装搜索条件
func (o *Order) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_at"] = 0

	if o.MemberId != 0 { //拼接会员ID
		maps["member_id"] = o.MemberId
	}

	if o.Number != "" { //拼接订单编号
		maps["number"] = o.Number
	}

	return maps
}
