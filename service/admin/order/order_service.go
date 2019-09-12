package order

import (
	"encoding/json"
	"fmt"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/order"
	"time"
)

type Order struct {
	Id            int
	MemberId      int
	Status        int
	Number        string
	Remark        string
	StartTime     string
	EndTime       string
	OrderField    string
	OrderSort     int
	Offset        int
	Limit         int
	RecommendName string
	MemberName    string
	BillNumber    string
	GoodsInfo     string
}

type GoodsInfo struct {
	Id     int
	Number int
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
		orders[key].CreateTimeDesc = timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32

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
	fields := "id,number,create_at,review_time,remark"
	orderDetail, orderErr := models.DetailOrder(o.Id, fields)
	if orderErr {
		err.Code = e.ERROR_SQL_FAIL
	}
	timeNow := time.Unix(int64(orderDetail.CreateAt), 0)
	orderDetail.CreateTimeDesc = timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32

	if orderDetail.ReviewTime == 0 {
		orderDetail.ReviewTimeDesc = "暂未审核"
	}
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
	maps := o.getMaps()
	maps["id"] = o.Id

	selectOrder, selectErr := models.DetailOrder(o.Id, "id,status") //获取订单详情
	fmt.Println(selectOrder)

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

	res := models.OrderStatusChange(maps, data)

	if res { //会员状态变化失败
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

func (o *Order) AddOrder() (err e.SelfError) {
	data := make(map[string]interface{})
	goodsInfo := o.GoodsInfo
	remark := o.Remark
	recommend := models.CheckAuth(o.RecommendName, 0)
	if recommend.ID == 0 {
		err.Code = e.ERROR_USER
		return
	}
	member := models.CheckAuth(o.MemberName, 0)
	if member.ID == 0 {
		err.Code = e.ERROR_USER
		return
	}
	recommendId := recommend.ID
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
