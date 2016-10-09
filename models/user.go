package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
