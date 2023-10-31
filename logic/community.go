package logic

import (
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/models"
	"go_web_app/pkg/snowflake"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (data any, err error) {
	// logic层会非常复杂
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

func CreatePost(post *models.Post) (err error) {
	// 1. 创建帖子
	post.PostID = snowflake.GetId() // 生成post_id
	// 插入mysql
	err = mysql.CreatePost(post)
	if err != nil {
		return err
	}
	// 获取帖子创建的时间，插入到redis中
	postdata, err := mysql.GetPostDetailById(post.PostID)
	if err != nil {
		zap.L().Error("GetPostDetailById(post.PostID) failed", zap.Error(err))
		return
	}
	err = redis.CreatePost(postdata)
	if err != nil {
		zap.L().Error("redis.CreatePost(postdata) failed", zap.Error(err))
		return
	}
	return
}

func GetPostDetailById(postId int64) (data *models.APIPostDetail, err error) {
	// 1. 查询帖子详情
	data, err = mysql.GetPostDetailById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(postId) failed", zap.Error(err))
		return nil, err
	}
	return
}

func GetPostListByCommunityId(id, offset, limit int64) (postlist []*models.APIPostDetail, err error) {
	// 1. 查询帖子列表, 这时候的post信息里面没有作者名字和社区名字
	postlist, err = mysql.GetPostListById(id, offset, limit)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(id) failed", zap.Error(err))
		return nil, err
	}

	return
}
