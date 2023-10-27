package pkg

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
)

func Init() {
	// 动态获取配置文件中的token过期时间
	TokenExpireDuration = time.Hour * time.Duration(viper.GetInt("auth.jwtexpire"))
}

// MyClaims 自定义声明, 用于生成token的结构体
type MyClaims struct {
	Username             string `json:"username"`
	UserID               int64  `json:"user_id"`
	jwt.RegisteredClaims        // 嵌套jwt.RegisteredClaims，包含了jwt的默认字段以及它的方法
}

const (
	SymetricKey = "go_web_app"
	Issuer      = "go_web_app"
)

var (
	// 默认值是24小时
	TokenExpireDuration time.Duration
)

// GenerateToken 生成JWT 返回accessToken和refreshToken
func GenerateToken(username string, userID int64) (accessToken, refreshToken string) {
	// 如果没有设置过期时间，就设置一个默认值	// 如果没有设置过期时间，就设置一个默认值
	if TokenExpireDuration == 0 {
		TokenExpireDuration = time.Hour * 24 // 默认值是24小时
	}
	// 1. 创建一个我们自己的声明
	c := MyClaims{
		username,
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    Issuer,
		},
	}

	// 2. 使用指定的签名方法创建一个access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	accessToken, err := token.SignedString([]byte(SymetricKey))
	if err != nil {
		accessToken = ""
	}

	// 3. 使用指定的签名方法创建一个refresh token, 有效期是access token的30倍
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration * 30)),
			Issuer:    Issuer,
		},
	})
	refreshToken, err = token.SignedString([]byte(SymetricKey))
	if err != nil {
		refreshToken = ""
	}

	return accessToken, refreshToken
}

// ParseToken 解析accessToken
func ParseToken(tokenString string) (*MyClaims, error) {
	// 1. 解析token
	claims := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 返回token对应的密钥
		return []byte(SymetricKey), nil
	})

	// 2. 校验token
	if err != nil {
		zap.L().Error("ParseToken failed", zap.Error(err))
		return nil, err
	}

	// 3. 返回claims
	if token.Valid {
		zap.L().Debug("ParseToken success")
		return claims, nil
	}

	zap.L().Error("ParseToken failed", zap.Error(err))
	return nil, err
}

// ParseRefreshToken 解析refreshToken
func ParseRefreshToken(aToken, rToken string) (newAToken, newrToken string) {
	// 1. 判断rToken 是否过期
	if _, err := ParseToken(rToken); err != nil {
		zap.L().Error("refreshToken过期了", zap.Error(err))
		return "", ""
	}

	// 2. 从aToken中解析出claims
	claims := new(MyClaims)
	_, err := jwt.ParseWithClaims(aToken, claims, func(token *jwt.Token) (interface{}, error) {
		// 返回token对应的密钥
		return []byte(SymetricKey), nil
	})
	// 3. 如果accessToken过期了，就用refreshToken生成一个新的accessToken
	if errors.Is(err, jwt.ErrTokenExpired) {
		return GenerateToken(claims.Username, claims.UserID)
	}

	zap.L().Error("accessToken过期了", zap.Error(err))
	return "", ""
}
