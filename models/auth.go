package models

type User struct {
	Model
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
}

func CheckAuth(username string) User {
	var user User

	/**
	根据用户名和密码查询对应的用户记录
	*/
	db.Select([]string{"id", "username", "password", "status"}).Where(User{Username: username, Status: 1}).First(&user)
	return user
}
