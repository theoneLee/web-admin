package models

import "gitee.com/muzipp/Distribution/pkg/logging"

type Member struct {
	Model
	RelationId int
	Name       string
	Sex        int
	IdCard     string
	Birth      string
	Phone      string
	SparePhone string
	Email      string
	BankCard   string
	Bank       string
	Status     int
	LevelId    int
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

func ListMembers(pageNum int, pageSize int, maps interface{}) (members []Member, flag bool) {

	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&members).Error
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
