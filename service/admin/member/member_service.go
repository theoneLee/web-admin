package member

import (
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/member"
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
	Offset     int
	Limit      int
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

	if res { //添加会员失败
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

func (m *Member) ListMembers() (members []map[string]interface{}, err e.SelfError) {
	memberRst, memberErr := models.ListMembers(m.Offset, m.Limit, m.getMaps())
	if memberErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	tempRst := make(map[string]interface{})
	for _, value := range memberRst {
		y, _, _ := member.GetTimeFromStrDate(value.Birth)
		tempRst["Id"] = value.ID
		tempRst["name"] = value.Name
		tempRst["IdCard"] = value.IdCard
		tempRst["age"] = member.GetAge(y)
		tempRst["phone"] = value.Phone
		tempRst["email"] = value.Email
		tempRst["bankCard"] = value.BankCard
		tempRst["bank"] = value.Bank
		tempRst["status"] = member.GetStatus(value.Status)
		tempRst["sex"] = member.GetSex(value.Sex)
		tempRst["sparePhone"] = value.SparePhone
		members = append(members, tempRst)
	}

	return
}

func (m *Member) CountMembers() (count int, err e.SelfError) {
	count, memberErr := models.CountMembers(m.getMaps())
	if memberErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

//封装搜索条件
func (m *Member) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_at"] = 0
	return maps
}
