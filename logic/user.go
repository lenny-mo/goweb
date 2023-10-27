package logic

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"go_web_app/dao/mysql"
	"go_web_app/models"
	"go_web_app/pkg"
	_ "go_web_app/pkg"
	"go_web_app/pkg/snowflake"
)

var (
	SALT = "web_app"
)

// Signup 如果用户不存在，就插入数据库
//
// 否则返回用户存在的错误
func Signup(params *models.SignupParam) error {
	// 1. 判断用户是否存在
	ok := mysql.CheckUserExist(params.Username)

	if ok {
		// 用户已经存在
		return errors.New("用户已经存在 or 查询出错")
	}
	// 构造一个User实例
	user := &models.User{
		UserID:   snowflake.GetId(), // 2. 生成UID
		Username: params.Username,
		PassWord: encryptPassword(params.Password), // 3. 加密密码
		Email:    params.Email,
		Gender:   params.Gender,
	}

	// 2. 保存User实例进数据库
	mysql.InsertUser(user)

	return nil
}

// Login 登录逻辑 成功则返回两个token
func Login(params *models.LoginParam) (string, string) {
	// 1. 从数据库中查找用户的密码信息并且进行比对
	user, err := mysql.GetUserByUsername(params.Username)
	if err == sql.ErrNoRows {
		return "", ""
	} else if err != nil {
		return "", ""
	}

	// 2. 比较数据库中的用户密码是否和用户输入的密码一致
	if user.PassWord != encryptPassword(params.Password) {
		return "", ""
	}

	// 3. 生成JWT
	return pkg.GenerateToken(user.Username, user.UserID)
}

// encryptPassword 对密码进行加密
func encryptPassword(password string) string {
	hash := sha1.New()
	// 对字符串进行加盐操作
	password += SALT

	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
