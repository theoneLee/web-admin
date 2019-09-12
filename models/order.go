package models

import (
	"gitee.com/muzipp/Distribution/pkg/common"
	"gitee.com/muzipp/Distribution/pkg/logging"
)

type Order struct {
	Model
	Number         string
	MemberId       int
	MemberName     string
	StatusDesc     string
	TeamName       string
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
	CreateTime     string
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
		Offset(pageNum).
		Limit(pageSize).
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
