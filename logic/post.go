package logic

import (
	"go_web_app/dao/redis"
	"go_web_app/models"

	"go.uber.org/zap"
)

func VoteForPost(p *models.VoteData, userid int64) error {

	err := redis.VoteForPost(p, userid)
	if err != nil {
		zap.L().Error("VoteForPost failed", zap.Error(err))
	}
	return err
}
