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

// GetPostDetailById 通过post id 查询帖子详情, 取中包含作者名字和社区名字
func GetPostDetailById(postId int64) (data *models.APIPostDetail, err error) {
	// 1. 查询帖子详情
	data = new(models.APIPostDetail)
	postData := new(models.Post)
	sqlstr := `select title, content, author_id, community_id, create_at, update_at from post where post_id = ?`
	err = sqlxdb.Get(postData, sqlstr, postId)
	if err != nil {
		zap.L().Error("sqlxdb.Get(data, sqlstr, postId) failed", zap.Error(err))
		return nil, err
	}
	data.Post = postData

	// 2. 根据作者id 查询作者信息
	userData := new(models.User)
	sqlstr = `select name, email, gender from user where user_id = ?`
	err = sqlxdb.Get(userData, sqlstr, postData.AuthorID)
	if err != nil {
		zap.L().Error("sqlxdb.Get(userData, sqlstr, postData.AuthorID) failed", zap.Int64("authorId", postData.AuthorID), zap.Error(err))
		return nil, err
	}
	data.AuthorName = userData.Username
	// 3. 根据社区id 查询社区信息
	communityData := new(models.Community)
	sqlstr = `select community_name, community_intro from community where community_id = ?`
	err = sqlxdb.Get(communityData, sqlstr, postData.CommunityID)
	if err != nil {
		zap.L().Error("sqlxdb.Get(communityData, sqlstr, postData.CommunityID) failed", zap.Int64("communityId", postData.CommunityID), zap.Error(err))
		return nil, err
	}
	data.CommunityName = communityData.Name
	return
}

// GetPostListByCommunityId 通过社区id 查询帖子列表
func GetPostListById(id, offset, limit int64) ([]*models.APIPostDetail, error) {
	postlist := []*models.Post{}
	sqlstr := "select " +
		"post_id, " +
		"title, " +
		"content, " +
		"author_id, " +
		"community_id, " +
		"create_at, " +
		"update_at " +
		"from post where community_id = ? " +
		"order by create_at desc" +
		"limit ?, ?"
	err := sqlxdb.Select(&postlist, sqlstr, id, offset, limit) // 跳过前面offset 条数据，取limit 条数据

	if err != nil {
		zap.L().Error("sqlxdb.Select(&postlist, sqlstr, id) failed", zap.Error(err))
		return nil, err
	}

	// 遍历postlist，查询每个帖子对应的作者信息和社区信息
	apiPostList, err := AddAuthorandCommunityName(postlist)

	return apiPostList, err
}

func AddAuthorandCommunityName(postlist []*models.Post) (apiPostList []*models.APIPostDetail, err error) {
	// 遍历postlist，查询每个帖子对应的作者信息和社区信息
	apiPostList = make([]*models.APIPostDetail, len(postlist), len(postlist))
	for i := range postlist {
		// 查询作者信息
		userData := new(models.User)
		usersql := "select name from user where user_id = ?"
		err = sqlxdb.Get(userData, usersql, postlist[i].AuthorID)
		if err != nil {
			zap.L().Error("sqlxdb.Get(userData, usersql, postlist[i].AuthorID) failed", zap.Error(err))
		}
		// 查询社区信息
		communityData := new(models.Community)
		communitysql := "select community_name from community where community_id = ?"
		err = sqlxdb.Get(communityData, communitysql, postlist[i].CommunityID)
		if err != nil {
			zap.L().Error("sqlxdb.Get(communityData, communitysql, postlist[i].CommunityID) failed", zap.Error(err))
		}
		// 把作者信息和社区信息和postlist 组合到一起
		apiPostList[i] = &models.APIPostDetail{
			AuthorName:    userData.Username,
			CommunityName: communityData.Name,
			Post:          postlist[i],
		}
	}

	return
}
