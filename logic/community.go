package logic

import (
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/models"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (data any, err error) {
	// 1. 查询所有社区的信息
	data, err = mysql.GetCommunityList()
	if err != nil {
		return nil, err
	}

	return
}

func GetCommunityDetailById(id int64) (data any, err error) {
	// 1. 查询社区详情
	data, err = mysql.GetCommunityDetailById(id)
	if err != nil {
		return nil, err
	}
	return
}

func CommunitySortedPost(cid string, p *models.PostListParam) ([]*models.APIPostDetail, error) {
	// 2. 判断这个zset是否存在，如果不存在则和p.order 做交集，产生一个新的zset
	// 如果这个集合存在，则直接返回这个集合的所有元素
	// 3. 查询redis，根据score返回这个zset的所有元素id
	ids, err := redis.GetCommunitySortedPostIds(cid, p)
	if err != nil {
		zap.L().Error("redis.GetCommunityPostIds(key, p) failed", zap.Error(err))
		return nil, err
	}

	// 4. 根据id list 查询mysql，返回帖子列表
	data, err := mysql.GetPostListByIds(ids)

	return data, nil
}
