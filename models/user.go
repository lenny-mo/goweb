package models

// User 结构体用于表示数据库中的用户表
type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"name"`
	PassWord string `db:"password"`
	Email    string `db:"email"`
	Gender   uint8  `db:"gender"`
}
