package models

type Member struct {
	Model
	RelationId int    `json:"relation_id"`
	Name       string `json:"name"`
	Sex        int    `json:"sex"`
	IdCard     string `json:"id_card"`
	Birth      string `json:"birth"`
	Phone      string `json:"phone"`
	SparePhone string `json:"spare_phone"`
	Email      string `json:"email"`
	BankCard   string `json:"bank_card"`
	Bank       string `json:"bank"`
	Status     int    `json:"status"`
	LevelId    int    `json:"level_id"`
}

func AddMember(data map[string]interface{}) bool {
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
		return false
	}

	return true
}
