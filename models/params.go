package models

type SignupParam struct {
	Username string `json:"username" binding:"required" ` // binding:"required"表示必须要传
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Gender   uint8  `json:"gender" binding:"required"`
}
