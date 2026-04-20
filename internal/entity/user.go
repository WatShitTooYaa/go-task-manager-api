package entity

type User struct {
	Id       uint16 `json:"id" validate:"required,min=3,max=200"`
	Username string `json:"username" validate:"required,min=3,max=200"`
	Password string `json:"password" validate:"required,min=3,max=200"`
}

type UserParam struct {
	Username string `json:"username" validate:"required,min=8,max=13"`
	Password string `json:"password" validate:"required,min=8,max=13"`
}
