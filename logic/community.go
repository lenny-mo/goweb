package logic

import (
	"go_web_app/dao/mysql"
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
