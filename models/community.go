package models

import "time"

type Community struct {
	ID          int64     `db:"id"`
	CommunityID int64     `db:"community_id"`
	Name        string    `db:"community_name"`
	Intro       string    `db:"community_intro"`
	CreateAt    time.Time `db:"create_at"`
	UpdateAt    time.Time `db:"update_at"`
}
