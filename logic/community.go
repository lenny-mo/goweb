package logic

import (
	"go_web_app/dao/mysql"
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
	err = mysql.CreatePost(post)
	if err != nil {
		return err
	}
	return
}

func GetPostDetailById(postId int64) (data *models.Post, err error) {
	// 1. 查询帖子详情
	data, err = mysql.GetPostDetailById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(postId) failed", zap.Error(err))
		return nil, err
	}
	return
}
