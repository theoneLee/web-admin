package models

import (
	"gitee.com/muzipp/Distribution/pkg/common"
	"gitee.com/muzipp/Distribution/pkg/logging"
	"github.com/jinzhu/gorm"
)

type Order struct {
	Model
	Number         string
	MemberId       int
	MemberName     string `gorm:"-"`
	StatusDesc     string `gorm:"-"`
	TeamName       string	`gorm:"-"`
	ReferencePrice float64
	ActualPrice    float64
	Discount       float64
	OtherIncome    float64
	Commission     float64
	Status         int
	Remark         string
	ShipTime       int
	ReviewTime     int
	Integral       int
	CreateTimeDesc string	`gorm:"-"`
	ReviewTimeDesc string	`gorm:"-"`
	RecommendId    int
	BillNumber     string
}

type OrderGoodsDetail struct {
	GoodsName     string
	Number        int
	MemberPrice   float64
	TotalPrice    float64
	Remark        string
	Specification string
}

func ListOrders(pageNum int, pageSize int, maps interface{}, fields string, remark string, startTime string, endTime string, orderField string, orderSort int) (orders []Order, flag bool) {

	query := Db.Table("order").
		Where(maps)

	if remark != "" { //拼接备注搜索
		query = query.Where("order.remark LIKE ?", "%"+remark+"%")
	}

	if startTime != "" {
		query = query.Where("order.create_at >= ?", common.GetTime(startTime))
	}

	if endTime != "" {
		query = query.Where("order.create_at <= ?", common.GetTime(endTime))
	}

	query = query.
		Joins("left join `member` as m on m.id = order.member_id").
		Joins("left join `member` as m1 on m1.id = m.relation_id").
		//Offset(pageNum).
		//Limit(pageSize).
		Select(fields)

	if orderField == "" { //没有自定义排序的情况，默认id倒叙
		query = query.Order("order.id desc")
	} else {
		if orderSort == 1 { //升序
			query = query.Order("order." + orderField + " asc")
		} else { //降序
			query = query.Order("order." + orderField + " desc")
		}
	}

	err := query.Scan(&orders).Error
	if err != nil {
		flag = true
		logging.Info("订单列表错误", err) //记录错误日志
		return
	}

	return
}

func CountOrders(maps interface{}, remark string, startTime string, endTime string) (count int, flag bool) {
	query := Db.Table("order").
		Where(maps)

	if remark != "" { //拼接备注搜索
		query = query.Where("order.remark LIKE ?", "%"+remark+"%")
	}

	if startTime != "" {
		query = query.Where("order.create_at >= ?", common.GetTime(startTime))
	}

	if endTime != "" {
		query = query.Where("order.create_at <= ?", common.GetTime(endTime))
	}

	err := query.Joins("left join `member` as m on m.id = order.member_id").
		Joins("left join `member` as m1 on m1.id = m.relation_id").
		Count(&count).Error
	if err != nil {
		flag = true
		logging.Info("订单量错误", err) //记录错误日志
		return
	}
	return
}

//商品详情
func DetailOrder(id int, fields string) (*Order, bool) {
	var order Order
	var flag bool
	err := Db.Table("order").
		Where("id = ? AND delete_at = ? ", id, 0).
		Select(fields).
		Find(&order).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		flag = true
		logging.Info("订单详情错误", err) //记录错误日志
		return &order, flag
	}

	if gorm.IsRecordNotFoundError(err) { //查询结果不存在的情况
		flag = true
		return nil, flag
	}
	return &order, flag
}

func DetailOrderGoods(orderId int, fields string) (orderGoods []OrderGoodsDetail, goodsFlag bool) {
	var flag bool
	err := Db.Table("order_goods").
		Where("order_goods.order_id = ? AND order_goods.delete_at = ? ", orderId, 0).
		Select(fields).
		Find(&orderGoods).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		flag = true
		logging.Info("订单商品详情错误", err) //记录错误日志
		return
	}

	if gorm.IsRecordNotFoundError(err) { //查询结果不存在的情况
		flag = true
		return nil, flag
	}
	return
}

func OrderStatusChange(maps interface{}, data map[string]interface{}) (flag bool) {
	err := Db.Debug().Model(Order{}).Where(maps).Update(data).Error

	if err != nil { //会员状态变化
		flag = true
		logging.Info("状态变化失败", err) //记录错误日志
		return
	}

	return
}

func AddOrder(data map[string]interface{}, tx *gorm.DB) (id int, flag bool) {
	order := &Order{
		RecommendId: data["recommend_id"].(int),
		MemberId:    data["member_id"].(int),
		Remark:      data["remark"].(string),
		BillNumber:  data["bill_number"].(string),
		Status:      data["status"].(int),
	}
	err := tx.Create(order).Error

	if err != nil { //添加商品失败
		flag = true
		logging.Info("添加订单错误", err) //记录错误日志
		return
	}

	return order.ID, false
}


func UpdateOrder(maps interface{}, data map[string]interface{}, tx *gorm.DB) (flag bool) {
	err := tx.Debug().Model(Order{}).Where(maps).Update(data).Error

	if err != nil { //会员状态变化
		flag = true
		logging.Info("更新订单失败", err) //记录错误日志
		return
	}

	return
}
