package mysql

import "go_web_app/models"

// InsertUser 函数用于在数据库中插入一条新的用户记录。
func InsertUser(user *models.User) error {
	sql := "insert into user(user_id, name,  password, email, gender)" +
		"values(?,?,?,?,?)"
	_, err := sqlxdb.Exec(sql, user.UserID, user.Username, user.PassWord, user.Email, user.Gender)

	return err
}

func QueryUserByName(username string) bool {
	sql := "select count(user_id) from user where name = ?"

	var count int
	err := sqlxdb.Get(&count, sql, username)
	if err != nil {
		return true
	}

	return count > 0
}

func CheckUserExist(username string) bool {

	return QueryUserByName(username)
}

// GetUserByUsername 函数用于根据用户名从数据库中获取一条只包括密码的用户记录。
func GetUserByUsername(username string) (*models.User, error) {
	sql := "select password, user_id, name from user where name = ?"
	user := &models.User{}
	err := sqlxdb.Get(user, sql, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
