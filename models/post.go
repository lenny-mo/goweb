package models

import (
	"time"
)

type Post struct {
	ID          int64     `db:"id" json:"id,string"`
	PostID      int64     `db:"post_id" json:"post_id,string"`
	AuthorID    int64     `db:"author_id" json:"author_id, string"`
	CommunityID int64     `db:"community_id" json:"community_id, string" binding:"required"`
	Status      int8      `db:"status" json:"status"`
	Score       int64     `db:"score" json:"score"`
	Title       string    `db:"title" json:"title" binding:"required"`
	Content     string    `db:"content" json:"content" binding:"required"`
	CreateAt    time.Time `db:"create_at" json:"create_at"`
	UpdateAt    time.Time `db:"update_at" json:"update_at"`
}

type APIPostDetail struct {
	AuthorName    string        `json:"author_name"`
	CommunityName string        `json:"community_name"`
	TotalVote     int64         `json:"total_vote"`
	AgreeVote     int64         `json:"agree_vote"`
	DisagreeVote  int64         `json:"disagree_vote"`
	*Post         `json:"Post"` // 嵌套匿名结构体
}
