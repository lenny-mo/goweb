package logic

import (
	"encoding/json"
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/models"

	"go.uber.org/zap"
)

// CommunitySortedPost 根据community id 返回对应的帖子排行榜
func CommunitySortedPost(cid string, p *models.PostListParam) ([]*models.Post, error) {
	// cache aside

	// 1. 查询redis
	// 1.1 判断排行榜是否在redis中存在, 如果存在则返回结果，否则查阅mysql
	if redis.ZsetExist(cid, p) {
		fmt.Println("走的是redis")
		if list, err := redis.CommunitySortedPost(cid, p); err != nil {
			zap.L().Error(err.Error())
			return nil, err
		} else {
			// 反序列化，并返回
			res := make([]*models.Post, 0, len(list))
			for _, s := range list {
				temp := new(models.Post)
				json.Unmarshal([]byte(s), temp)
				res = append(res, temp)
			}
			return res, err
		}
	}

	// 2. 如果redis没有则查阅mysql
	fmt.Println("走的是mysql")
	return mysql.CommunitySortedPost(cid, p)
}
