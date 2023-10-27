package mysql

import "errors"

var (
	ErrorUserExist          = errors.New("用户已经存在")
	ErrorUserNotExist       = errors.New("用户不存在")
	ErrorQueryFailed        = errors.New("查询出错")
	ErrorInvalidCommunityId = errors.New("无效的社区ID")
)
