package logic

import (
	"encoding/json"
	"errors"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/models"
	"go_web_app/pkg/snowflake"
	"time"

	redisv9 "github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 1. 创建帖子
	post.PostID = snowflake.GetId() // 生成post_id
	// 2. 添加时间
	post.CreateAt, post.UpdateAt = time.Now(), time.Now()

	// cache aside pattern
	// 3. 插入mysql
	err = mysql.CreatePost(post)
	if err != nil {
		return err
	}

	// 4. 把帖子的数据插入到redis
	err = redis.CreatePost(post)
	if err != nil {
		zap.L().Error("redis.CreatePost(postdata) failed", zap.Error(err))
		return
	}
	return
}

func GetPostDetailById(postId int64) (data interface{}, err error) {
	// cache aside pattern
	// 先缓存再看数据库
	if bytes, err := redis.GetPostById(postId); err != nil {
		if !errors.Is(err, redisv9.Nil) {
			zap.L().Error(err.Error())
			return data, err
		}
	} else {
		// 反序列化数据并返回
		zap.L().Debug("走缓存了兄弟")
		p := new(models.Post)
		if err := json.Unmarshal(bytes, p); err != nil {
			zap.L().Error(err.Error())
			return data, err
		}
		return p, err
	}

	// 从mysql中找
	data, err = mysql.GetPostDetailById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(postId) failed", zap.Error(err))
		return nil, err
	}
	return
}

func VoteForPost(p *models.VoteData, userid int64) error {
	// write back 策略
	// 1. 更新缓存数据，标记该条数据为脏数据
	// redis中有个单独的set dirty set, 里面记录了脏数据的post_id
	// 有一个单独的goroutine，定时任务，把里面的脏数据post flush到数据库
	err := redis.VoteForPost(p, userid)
	if err != nil {
		zap.L().Error("VoteForPost failed", zap.Error(err))
	}

	return err
}
