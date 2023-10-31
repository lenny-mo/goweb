package redis

import (
	"errors"
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/models"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

const (
	TWOWEEKSECOND = 60 * 60 * 24 * 7 * 2 // 两周的秒数
	// 赞成票的分数，反对票的分数
	AGREEPOSTSCORE    = 432
	DISAGREEPOSTSCORE = -432
)

var (
	VoteTimeExpireError = errors.New("投票时间已过")
)

/**
 * 投票功能说明：
 *
 * - `vote` 表示投票方向：
 *   - 1: 投赞成票
 *       1. 之前没有投过票，现在投赞成票
 *       2. 之前投反对票，现在改投赞成票
 *   - 0: 取消投票
 *       1. 之前投过赞成票，现在要取消投票
 *       2. 之前投过反对票，现在要取消投票
 *   - -1: 投反对票
 *       1. 之前没有投过票，现在投反对票
 *       2. 之前投赞成票，现在改投反对票
 *
 * 投票限制：
 *   1. 每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
 *   2. 到期之后，将redis中保存的赞成票数及反对票数存储到mysql表中。
 *   3. 到期之后删除redis key: "KeyPostVotedZSetPF"（此key用于...）。
 *
 *
 */

func VoteForPost(v *models.VoteData, userid int64) error {
	// 1. 判断投票时间是否在有效期内
	// zset1 存储所有帖子的发帖时间
	// 从zset1中获取post发帖时间，使用formatint 10进制转换
	postIDstr := strconv.FormatInt(v.PostID, 10)
	userIDstr := strconv.FormatInt(userid, 10)
	postTimeFloat := redisClient.ZScore(PostTimeZSetKey, postIDstr).Val()
	postTime := int64(postTimeFloat)
	fmt.Println("posttime", postTime)

	if postTime == 0 {
		// 没有记录，说明帖子记录没有插入到zset1中, 从mysql中获取帖子信息并且插入到zset1中
		postdata, err := mysql.GetPostDetailById(v.PostID)
		if err != nil {
			zap.L().Error("mysql.GetPostDetailById(v.PostID) failed", zap.Error(err))
			return err
		}
		postTime = postdata.CreateAt.Unix()
		redisClient.ZAdd(PostTimeZSetKey, redis.Z{Score: float64(postTime), Member: postIDstr})
	} else if time.Now().Unix()-postTime > TWOWEEKSECOND { // 超过两周，不允许再投票
		// 如果posttime大于0, 判断时间是否超过两周
		return VoteTimeExpireError
	}

	// 2. 获取当前用户给当前帖子的投票记录
	// zset3 存储当前帖子的所有投票记录
	currentUserVoteFloat := redisClient.ZScore(PostVotedPrefix+postIDstr, userIDstr).Val()
	currentVote := int64(currentUserVoteFloat)

	// 3. 判断投票方向
	// TODO: 需要优化一下错误处理
	switch currentVote {
	case 1: // 之前投过赞成票
		switch v.Vote {
		case 0: // 之前投过赞成票，现在要取消投票
			// 从zset3 中删除当前用户的投票记录, 从zset2中扣除赞成票分数
			cmd1 := redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		case -1: // 之前投过赞成票，现在要投反对票, 先更新之前的投票记录，再增加反对票分数
			cmd1 := redisClient.ZIncrBy(PostVotedPrefix+postIDstr, 2*float64(v.Vote), userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, 2*DISAGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVotedPrefix+postIDstr, 2*float64(v.Vote), userIDstr) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, 2*DISAGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		}
	case -1: // 之前投过反对票
		switch v.Vote {
		case 1: // 之前投过反对票，现在要投赞成票
			// 先修改zset3中的投票记录，再增加zset2赞成票分数
			cmd1 := redisClient.ZIncrBy(PostVotedPrefix+postIDstr, 2*float64(v.Vote), userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, 2*AGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVotedPrefix+postIDstr, 2*float64(v.Vote), userIDstr) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, 2*AGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		case 0: // 之前投过反对票，现在要取消投票
			// 从zset3 中删除当前用户的投票记录, 从zset2中扣除反对票分数
			cmd1 := redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		}
	case 0: // 之前没有投过票
		switch v.Vote {
		case 1: // 之前没有投过票，现在要投赞成票
			// 在zset3中增加当前用户的投票记录, 在zset2中增加赞成票分数
			cmd1 := redisClient.ZAdd(PostVotedPrefix+postIDstr, redis.Z{Score: float64(v.Vote), Member: userIDstr})
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, AGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZAdd(PostVotedPrefix+postIDstr, redis.Z{Score: float64(v.Vote), Member: userIDstr}) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, AGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		case -1: // 之前没有投过票，现在要投反对票
			// 在zset3中增加当前用户的投票记录, 在zset2中增加反对票分数
			cmd1 := redisClient.ZAdd(PostVotedPrefix+postIDstr, redis.Z{Score: float64(v.Vote), Member: userIDstr})
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr)
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZAdd(PostVotedPrefix+postIDstr, redis.Z{Score: float64(v.Vote), Member: userIDstr}) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		}
	}

	zap.L().Info("Update post vote success", zap.Int64("post_id", v.PostID), zap.Int8("vote", v.Vote), zap.Int64("user_id", userid))
	return nil
}

// CreatePost 将帖子创建时间插入到redis中
func CreatePost(post *models.APIPostDetail) error {
	pipeline := redisClient.TxPipeline()

	// 1. 将帖子创建时间插入到redis的posttimezset中
	pipeline.ZAdd(PostTimeZSetKey, redis.Z{Score: float64(post.Post.CreateAt.Unix()), Member: strconv.FormatInt(post.Post.PostID, 10)})

	// 2. 使用当前时间作为帖子的初始score
	pipeline.ZAdd(PostVoteZSetKey, redis.Z{Score: float64(post.Post.CreateAt.Unix()), Member: strconv.FormatInt(post.Post.PostID, 10)})

	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return err
	}
	fmt.Println("创建post并且更新redis")
	zap.L().Info("CreatePost success", zap.Int64("post_id", post.Post.PostID))

	// todo: 尝试从redis中取出刚刚存储的值

	return nil
}
