package models

import "gitee.com/muzipp/Distribution/pkg/logging"

type Member struct {
	Model
	RelationId       int    //上级ID
	RelationName     string `gorm:"-"` //上级名称
	Name             string //姓名
	Sex              int    //性别
	IdCard           string //身份证
	Birth            string //生日
	Age              int    `gorm:"-"` //年龄
	Phone            string //手机号
	SparePhone       string //备用电话
	Email            string //EMAIL
	BankCard         string //账户
	Bank             string //银行
	LevelId          int    //等级ID
	Level            Level
	LevelName        string  `gorm:"-"` //等级文案
	Status           int     //状态
	StatusDesc       string  `gorm:"-"` //状态对应的文案
	AvailableIncome  float64 //可提取佣金
	ExtractIncome    float64 //已提取佣金
	TotalOrderIncome float64 `gorm:"-"` //订单总额
	TotalOrderNumber int     `gorm:"-"` //订单数量
	ExpiredTime      int     //到期时间
}

func AddMember(data map[string]interface{}) (flag bool) {
	err := db.Create(&Member{
		RelationId: data["relation_id"].(int),
		Name:       data["name"].(string),
		Sex:        data["sex"].(int),
		IdCard:     data["id_card"].(string),
		Birth:      data["birth"].(string),
		Phone:      data["phone"].(string),
		SparePhone: data["spare_phone"].(string),
		Email:      data["email"].(string),
		BankCard:   data["bank_card"].(string),
		Bank:       data["bank"].(string),
		Status:     data["status"].(int),
		LevelId:    data["level_id"].(int),
	}).Error

	if err != nil { //添加会员失败
		flag = true
		logging.Info("添加会员错误", err) //记录错误日志
		return
	}

	return
}

func ListMembers(pageNum int, pageSize int, maps interface{}, fields []string) (members []Member, flag bool) {

	err := db.Preload("Level").Select(fields).Where(maps).Offset(pageNum).Limit(pageSize).Find(&members).Error
	if err != nil {
		flag = true
		logging.Info("会员列表错误", err) //记录错误日志
		return
	}
	return
}

func CountMembers(maps interface{}) (count int, flag bool) {
	err := db.Model(&Member{}).Where(maps).Count(&count).Error
	if err != nil {
		flag = true
		logging.Info("会员人数错误", err) //记录错误日志
		return
	}
	return
}
