package logic

import (
	"go_web_app/dao/mysql"
	"go_web_app/models"
	"go_web_app/pkg/snowflake"
)

func Signup(params *models.SignupParam) error {
	// 1. 判断用户是否存在
	ok := mysql.QueryUserByName()
	if ok {
		// 用户已经存在
		return nil
	}
	// 2. 生成UID
	snowflake.GetId()

	// 3. 密码加密

	// 2. 保存进数据库
	mysql.InsertUser()

	return nil
}
