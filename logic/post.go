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

// GetPostListByCommunityId 根据社区id查询帖子列表
//
// 没有排行榜也没有投票数据
// TODO: 深分页查询优化
func GetPostListByCommunityId(id, offset, limit int64) (postlist []*models.APIPostDetail, err error) {
	// 1. 查询帖子列表, 这时候的post信息里面没有作者名字和社区名字
	postlist, err = mysql.GetPostListById(id, offset, limit)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(id) failed", zap.Error(err))
		return nil, err
	}

	return
}

func VoteForPost(p *models.VoteData, userid int64) error {

	err := redis.VoteForPost(p, userid)
	if err != nil {
		zap.L().Error("VoteForPost failed", zap.Error(err))
	}
	return err
}

// GetPostList 获取帖子列表, 根据排序列表进行排序
//
// post:time 根据时间进行排序； post:vote 根据投票分数进行排序
func GetSortedPost(p *models.PostListParam) ([]*models.APIPostDetail, error) {
	// 1. 根据用户请求参数offset, limit 去redis查询id列表
	// order表示查询的zset集合, 1表示post:time, 2表示post:score
	ids, err := redis.GetPostList(p)
	if err != nil {
		zap.L().Error("redis.GetPostList(p) failed", zap.Error(err))
		return nil, err
	}
	if len(ids) == 0 { // redis 中没有数据
		zap.L().Warn("redis.GetPostList(p) len(ids) == 0")
		return nil, nil
	}

	// 2. 查询redis 获取到 post对应的投票总数，赞成票数，反对票数, 格式使用slice
	totalVote, err1 := redis.GetPostVoteData(ids)
	agreeVote, err2 := redis.GetPostAgreeVoteData(ids)
	disagreeVote, err3 := redis.GetPostDisagreeVoteData(ids)

	if err1 != nil {
		zap.L().Error("GetPostVoteData(idlist) failed", zap.Error(err))
		return nil, err
	}
	if err2 != nil {
		zap.L().Error("GetPostAgreeVoteData(idlist) failed", zap.Error(err))
		return nil, err
	}
	if err3 != nil {
		zap.L().Error("GetPostDisagreeVoteData(idlist) failed", zap.Error(err))
		return nil, err
	}

	// 3. 根据id list 查询mysql得到帖子列表
	data, err := mysql.GetPostListByIds(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIds(ids) failed", zap.Error(err))
		return nil, err
	}

	// 把每个帖子的投票数据放到对应的帖子中
	for i := range data {
		data[i].TotalVote = totalVote[i]
		data[i].AgreeVote = agreeVote[i]
		data[i].DisagreeVote = disagreeVote[i]
	}

	// 4. 返回帖子列表
	return data, nil
}
