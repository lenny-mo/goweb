package mysql

import (
	"go_web_app/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetPostListByIds 根据id列表查询帖子列表
func GetPostListByIds(idlist []string) ([]*models.APIPostDetail, error) {
	// 1. 查询帖子详情
	postlist := []*models.Post{}
	sqlstr := "select " +
		"post_id, " +
		"title, " +
		"content, " +
		"author_id, " +
		"community_id, " +
		"create_at, " +
		"update_at " +
		"from post where post_id in (?) " +
		"order by find_in_set(post_id, ?)"
	// using In: query is now: "SELECT name FROM users WHERE id IN (?, ?, ?) order by find_in_set(id, ?, ?, ?)"
	// args is now: [1, 2, 3]
	query, args, err := sqlx.In(sqlstr, idlist, strings.Join(idlist, ",")) // 跳过前面offset 条数据，取limit 条数据
	if err != nil {
		return nil, err
	}

	query = sqlxdb.Rebind(query) //会根据当前的数据库驱动（在此例中是 PostgreSQL）自动替换占位符
	err = sqlxdb.Select(&postlist, query, args...)
	if err != nil {
		return nil, err
	}

	// 3. 拼接数据, 添加作者名字和社区名字
	apiPostList, err := AddAuthorandCommunityName(postlist)
	if err != nil {
		return nil, err
	}
	// 4. 返回
	return apiPostList, nil
}
