package order

import (
	"encoding/json"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/order"
	"time"
)

type Order struct {
	Id          int
	MemberId    int
	Status      int
	Number      string
	Remark      string
	StartTime   string
	EndTime     string
	OrderField  string
	OrderSort   int
	Offset      int
	Limit       int
	RecommendId int
	MemberName  string
	BillNumber  string
	GoodsInfo   string
}

type GoodsInfo struct {
	Id     int
	Number int
}

//订单列表
func (o *Order) ListOrders() (orders []models.Order, err e.SelfError) {
	fields := "order.id,order.bill_number,order.ship_time,m.name as member_name, m.username as member_user_name,order.number,order.create_at, order.status," +
		"order.reference_price,order.actual_price,order.discount,order.integral," +
		"order.remark,m1.name as recommend_name,m1.username as recommend_user_name,m1.operate_address as recommend_address"
	orders, ordersErr := models.ListOrders(o.Offset, o.Limit, o.getMaps(), fields, o.Remark, o.StartTime, o.EndTime, o.OrderField, o.OrderSort)
	if ordersErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	//处理显示文案部分
	for key, value := range orders {
		orders[key].StatusDesc = order.GetStatus(value.Status)
		orders[key].CreateTimeDesc = time.Unix(int64(value.CreateAt), 0).Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
		if value.ShipTime == 0 {
			orders[key].ShipTimeDesc = "未发货"
		} else {
			orders[key].ShipTimeDesc = time.Unix(int64(value.ShipTime), 0).Format("2006-01-02 15:04:05")
		}
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

//获取文章（redis不存在读取数据库）
func (o *Order) DetailOrder() (rst map[string]interface{}, err e.SelfError) {
	rst = make(map[string]interface{}, 10)
	fields := "order.id,order.number,order.create_at,order.ship_time,order.remark,order.bill_number,order.status," +
		"m.name as member_name,m.username as member_username," +
		"m1.name as recommend_name,m1.username as recommend_user_name,m1.operate_address as recommend_address"
	orderDetail, orderErr := models.DetailOrder(o.Id, fields)
	if orderErr {
		err.Code = e.ERROR_SQL_FAIL
	}
	orderDetail.CreateTimeDesc = time.Unix(int64(orderDetail.CreateAt), 0).Format("2006-01-02 15:04:05") //2015-06-15 08:52:32

	if orderDetail.ShipTime == 0 {
		orderDetail.ShipTimeDesc = "暂未审核"
	} else {
		orderDetail.ShipTimeDesc = time.Unix(int64(orderDetail.ShipTime), 0).Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
	}
	orderDetail.StatusDesc = order.GetStatus(orderDetail.Status)

	orderGoodsFields := "order_goods.goods_name,order_goods.number," +
		"order_goods.specification,order_goods.member_price,order_goods.remark,order_goods.total_price"
	orderDetailGoods, orderDetailGoodsFlag := models.DetailOrderGoods(o.Id, orderGoodsFields)

	if orderDetailGoodsFlag {
		err.Code = e.ERROR_SQL_FAIL
	}
	rst["orderInfo"] = orderDetail
	rst["goodsInfo"] = orderDetailGoods

	return
}

//添加会员代码
func (o *Order) StatusChange() (err e.SelfError) {
	data := make(map[string]interface{})
	data["status"] = o.Status
	if o.Status == 1 {
		data["ship_time"] = int(time.Now().Unix())
	}
	maps := o.getMaps()
	maps["order.id"] = o.Id

	selectOrder, selectErr := models.DetailOrder(o.Id, "order.id,order.status") //获取订单详情

	if selectErr || selectOrder == nil {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	if selectOrder.Status == 1 && o.Status == -3 {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	if selectOrder.Status == -3 && o.Status == 1 {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	if selectOrder.Status == o.Status {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	updateMaps := o.getMaps()
	updateMaps["id"] = o.Id
	res := models.OrderStatusChange(updateMaps, data)

	if res { //会员状态变化失败
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

func (o *Order) AddOrder() (err e.SelfError) {
	data := make(map[string]interface{})
	goodsInfo := o.GoodsInfo
	remark := o.Remark
	member := models.CheckAuth(o.MemberName, 0)
	if member.ID == 0 {
		err.Code = e.ERROR_USER
		return
	}
	recommend := models.CheckUser(o.RecommendId)

	if recommend.IsOperate == 0  {
		err.Code = e.ERROR_OPERATE
		return
	}
	recommendId := o.RecommendId
	memberId := member.ID
	billNumber := o.BillNumber

	var totalPrice float64
	var sumPrice float64
	var totalIntegral int

	//开启事务
	tx := models.Db.Begin()

	data["recommend_id"] = recommendId
	data["member_id"] = memberId
	data["remark"] = remark
	data["bill_number"] = billNumber
	data["status"] = -1
	data["number"] = order.GenerateNumber()

	orderId, orderAddFlag := models.AddOrder(data, tx)
	if orderAddFlag { //json解析错误
		err.Code = e.ERROR_SQL_FAIL
		tx.Rollback()
		return
	}

	//获取商品信息
	goodsInfoFormat := make([]GoodsInfo, 5)
	jsonErr := json.Unmarshal([]byte(goodsInfo), &goodsInfoFormat)
	if jsonErr != nil { //json解析错误
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	//遍历处理商品信息
	for _, value := range goodsInfoFormat {
		goodsData := make(map[string]interface{})

		goods, _ := models.DetailGoods(value.Id, "*")
		sumPrice = sumPrice + float64(value.Number)*goods.MemberPrice
		totalPrice = totalPrice + float64(value.Number)*goods.Price
		totalIntegral = totalIntegral + value.Number*goods.Integral

		goodsData["order_id"] = orderId
		goodsData["goods_id"] = goods.ID
		goodsData["goods_name"] = goods.Name
		goodsData["goods_img"] = goods.Img
		goodsData["member_price"] = goods.MemberPrice
		goodsData["price"] = goods.Price
		goodsData["number"] = value.Number
		goodsData["total_price"] = goods.Price
		goodsData["sum_price"] = goods.MemberPrice
		goodsData["specification"] = goods.Specification
		goodsData["remark"] = goods.Remark
		OrderAddGoodsFlag := models.AddOrderGoods(goodsData, tx)

		if OrderAddGoodsFlag {
			err.Code = e.ERROR_SQL_FAIL
			tx.Rollback()
			return
		}

	}

	//更新订单表的数据
	updateOrderData := make(map[string]interface{})
	updateOrderData["reference_price"] = sumPrice
	updateOrderData["actual_price"] = totalPrice
	updateOrderData["discount"] = totalPrice - sumPrice
	updateOrderData["integral"] = totalIntegral

	maps := make(map[string]interface{})
	maps["delete_at"] = 0
	maps["id"] = orderId
	updateFlag := models.UpdateOrder(maps, updateOrderData, tx)
	if updateFlag {
		err.Code = e.ERROR_SQL_FAIL
		tx.Rollback()
		return
	}

	tx.Commit() //事务提交
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
