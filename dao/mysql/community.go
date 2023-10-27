package mysql

import (
	"database/sql"
	"go_web_app/models"

	"go.uber.org/zap"
)

func GetCommunityList() ([]models.Community, error) {
	// 查询操作
	tx, err := sqlxdb.Beginx()
	if err != nil {
		zap.L().Error("begin transaction failed", zap.Error(err))
		return nil, err
	}
	sqlstr := "select id, community_id, community_name, community_intro, create_at, update_at from community"
	var communityList []models.Community
	// 这里使用Query 方法，当然也可以使用Select 方法： Select(&communityList, sqlstr)
	// Query 需要手动迭代， Select 不需要
	err = sqlxdb.Select(&communityList, sqlstr)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows { // 没有查到数据
			zap.L().Warn("there is no community in table community")
			return nil, nil
		} else {
			zap.L().Error("query community list failed", zap.Error(err))
			return nil, err
		}
	}
	// 提交事务
	tx.Commit()

	// 返回
	return communityList, nil
}

func GetCommunityDetailById(id int64) (communityDetail *models.Community, err error) {
	// 初始化communityDetail
	communityDetail = new(models.Community)
	// 开启事务
	tx, err := sqlxdb.Beginx()
	if err != nil {
		zap.L().Error("begin transaction failed", zap.Error(err))
		return nil, err
	}

	sqlstr := "select id, community_id, community_name, community_intro, create_at, update_at from community where id = ?"
	err = sqlxdb.Get(communityDetail, sqlstr, id)
	if err != nil {
		zap.L().Error("query community detail failed", zap.Error(err))
		err = ErrorInvalidCommunityId
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return
}

func CreatePost(post *models.Post) error {
	// 开启事务
	tx, err := sqlxdb.Beginx()
	if err != nil {
		zap.L().Error("begin transaction failed", zap.Error(err))
		return err
	}

	sqlstr := "insert into post(post_id, title, content, author_id, community_id) values(?, ?, ?, ?, ?)"
	_, err = sqlxdb.Exec(sqlstr, post.PostID, post.Title, post.Content, post.AuthorID, post.CommunityID)
	if err != nil {
		zap.L().Error("insert into post failed", zap.Error(err))
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}

func GetPostDetailById(postId int64) (data *models.Post, err error) {
	// 1. 查询帖子详情
	data = new(models.Post)
	sqlstr := `select title, content, author_id, community_id, create_at, update_at from post where post_id = ?`
	err = sqlxdb.Get(data, sqlstr, postId)
	if err != nil {
		zap.L().Error("sqlxdb.Get(data, sqlstr, postId) failed", zap.Error(err))
		return nil, err
	}
	return
}
