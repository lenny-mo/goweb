package models

// 请求的参数结构体，在logic层, controllers层都会用到，
// 所以放在models层
type SignupParam struct {
	Username string `json:"username" binding:"required" ` // binding:"required"表示必须要传
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Gender   uint8  `json:"gender" binding:"required"`
}

type LoginParam struct {
	Username string `json:"username" binding:"required" ` // binding:"required"表示必须要传
	Password string `json:"password" binding:"required"`
}
