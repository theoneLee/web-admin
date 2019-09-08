package member

import (
	"fmt"
	"gitee.com/muzipp/Distribution/models"
	"gitee.com/muzipp/Distribution/pkg/e"
	"gitee.com/muzipp/Distribution/pkg/member"
)

type Member struct {
	Id         int
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
	data["relation_id"] = m.RelationId

	res := models.AddMember(data)

	if res { //添加会员失败
		err.Code = e.ERROR_SQL_FAIL
	}

	return
}

func (m *Member) ListMembers() (members []models.Member, err e.SelfError) {
	fields := "member.id,member.name,member.status,member.sex,member.id_card," +
		"member.birth,member.phone,member.spare_phone,member.email,member.bank," +
		"member.bank_card,member.available_income,member.extract_income," +
		"l.name as level_name,m1.name as relation_name," +
		"count(o.id) as total_order_number,sum(o.reference_price) as total_order_income"
	members, memberErr := models.ListMembers(m.Offset, m.Limit, m.getMaps(), fields)
	fmt.Println(members)
	if memberErr {
		err.Code = e.ERROR_SQL_FAIL
	}

	for key, value := range members {
		members[key].StatusDesc = member.GetStatus(value.Status)
		members[key].SexDesc = member.GetSex(value.Sex)
		members[key].Age = member.GetAge(member.GetTimeFromStrDate(value.Birth))
	}

	return

}

//添加会员代码
func (m *Member) StatusChange() (err e.SelfError) {
	data := make(map[string]interface{})
	data["status"] = m.Status
	maps := m.getMaps()
	maps["id"] = m.Id

	selectMember, selectErr := models.DetailMember(m.Id, "id,status")//获取订单详情

	if selectErr || selectMember == nil {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	if selectMember.Status == m.Status {
		err.Code = e.ERROR_SQL_FAIL
		return
	}

	res := models.StatusChange(maps, data)

	if res { //会员状态变化失败
		err.Code = e.ERROR_SQL_FAIL
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
