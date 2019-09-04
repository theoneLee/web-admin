package member

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
)

type Member struct {
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

//添加会员代码
func (m *Member) AddMember() (err e.SelfError) {
	data := make(map[string]interface{})
	data["relation_id"] = m.RelationId
	data["name"] = m.Name
	data["sex"] = m.Sex
	data["id_card"] = m.IdCard
	data["birth"] = m.Birth
	data["phone"] = m.Phone
	data["spare_phone"] = m.SparePhone
	data["email"] = m.Email
	data["bank_card"] = m.BankCard
	data["bank"] = m.Bank
	data["status"] = m.Status
	data["level_id"] = m.LevelId

	res := models.AddMember(data)

	if !res { //添加会员失败
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}
