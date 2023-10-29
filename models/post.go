package models

import "time"

type Post struct {
	ID          int64     `db:"id" json:"id"`
	PostID      int64     `db:"post_id" json:"post_id"`
	AuthorID    int64     `db:"author_id" json:"author_id"`
	CommunityID int64     `db:"community_id" json:"community_id" binding:"required"`
	Status      int8      `db:"status" json:"status"`
	Title       string    `db:"title" json:"title" binding:"required"`
	Content     string    `db:"content" json:"content" binding:"required"`
	CreateAt    time.Time `db:"create_at" json:"create_at"`
	UpdateAt    time.Time `db:"update_at" json:"update_at"`
}

type APIPostDetail struct {
	AuthorName    string        `json:"author_name"`
	CommunityName string        `json:"community_name"`
	*Post         `json:"Post"` // 嵌套匿名结构体
}
