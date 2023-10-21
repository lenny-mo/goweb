package models

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"name"`
	PassWord string `db:"password"`
	Email    string `db:"email"`
	Gender   uint8  `db:"gender"`
}
