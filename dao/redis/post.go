package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_web_app/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
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
 * 投票功能说明：这个功能具有幂等性，用户的重复赞成和反对并不会影响结果
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
 *   1. 每个贴子自发表之日起两个星期之内允许用户投票，超过两个星期就不允许再投票了。
 *   2. 到期之后，将redis中保存的赞成票数及反对票数存储到mysql表中。
 *   3. 到期之后删除redis key: "KeyPostVotedZSetPF"（此key用于...）。
 *
 *
 */

func GetPostById(postid int64) ([]byte, error) {
	return redisClient.Get(strconv.FormatInt(postid, 10)).Bytes()
}

func VoteForPost(v *models.VoteData, userid int64) error {
	// 1. 判断投票时间是否在有效期内
	postIDstr := strconv.FormatInt(v.PostID, 10)
	userIDstr := strconv.FormatInt(userid, 10)

	postdata, _ := GetPostById(v.PostID)
	postTimeFloat := redisClient.ZScore(PostTimeZSetKey+":"+v.CommunityID, string(postdata)).Val()
	postTime := int64(postTimeFloat)
	fmt.Println("posttime", postTime)

	if postTime == 0 {
		// 没有记录
		return redis.Nil
	} else if time.Now().Unix()-postTime > TWOWEEKSECOND { // 超过两周，不允许再投票
		// 如果posttime大于0, 判断时间是否超过两周
		return VoteTimeExpireError
	}

	// 2. 获取当前用户给当前帖子的投票记录，根据历史操作来判断接下来的动作
	// 在PostVoted zset中查找key=userid 的分数, 如果用户没有投票，返回默认0值
	currentUserVoteFloat := redisClient.ZScore(PostVotedPrefix+postIDstr, userIDstr).Val()
	currentVote := int64(currentUserVoteFloat)

	// 3. 标记为脏数据
	cmd3 := redisClient.SAdd(Dirty, postIDstr)
	if cmd3.Err() != nil {
		zap.L().Error(cmd3.Err().Error())
		return cmd3.Err()
	}

	// 4. 判断投票方向
	switch currentVote {
	case 1: // 之前投过赞成票
		switch v.Vote {
		case 0: // 之前投过赞成票，现在要取消投票
			// 删除当前用户的投票记录, 从zset2中扣除赞成票分数
			cmd1 := redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey+":"+v.CommunityID, DISAGREEPOSTSCORE, string(postdata))
			if cmd1.Err() != nil {
				zap.L().Error("redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr) failed", zap.Error(cmd1.Err()))
				return cmd1.Err()
			}
			if cmd2.Err() != nil {
				zap.L().Error("redisClient.ZIncrBy(PostVoteZSetKey, DISAGREEPOSTSCORE, postIDstr) failed", zap.Error(cmd2.Err()))
				return cmd2.Err()
			}
		}
	case -1: // 之前投过反对票
		switch v.Vote {
		case 0: // 之前投过反对票，现在要取消投票
			// 删除当前用户的投票记录, 扣除反对票分数
			cmd1 := redisClient.ZRem(PostVotedPrefix+postIDstr, userIDstr)
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey+":"+v.CommunityID, DISAGREEPOSTSCORE, string(postdata))
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
			// 增加当前用户的投票记录, 增加赞成票分数
			cmd1 := redisClient.ZAdd(PostVotedPrefix+postIDstr, redis.Z{Score: float64(v.Vote), Member: userIDstr})
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey+":"+v.CommunityID, AGREEPOSTSCORE, string(postdata))
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
			cmd2 := redisClient.ZIncrBy(PostVoteZSetKey+":"+v.CommunityID, DISAGREEPOSTSCORE, string(postdata))
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
func CreatePost(post *models.Post) error {

	// 序列化
	bytes, err := json.Marshal(post)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	// 添加key value =  postid, post
	err = redisClient.Set(strconv.FormatInt(post.PostID, 10), bytes, 24*time.Hour).Err()
	if err != nil {
		zap.L().Error("Set failed", zap.Error(err))
		return err
	}

	// 1. 将帖子创建时间插入到redis的post:time zset中
	err = redisClient.ZAdd(PostTimeZSetKey+":"+strconv.FormatInt(post.CommunityID, 10), redis.Z{Score: float64(post.CreateAt.Unix()), Member: bytes}).Err()
	if err != nil {
		zap.L().Error("ZAdd (PostTimeZSet) failed", zap.Error(err))
		return err
	}

	// 2. 使用当前时间作为帖子的初始score插入到post:score zset中
	err = redisClient.ZAdd(PostVoteZSetKey+":"+strconv.FormatInt(post.CommunityID, 10), redis.Z{Score: float64(post.Score), Member: bytes}).Err()
	if err != nil {
		zap.L().Error("ZAdd (PostVoteZSet) failed", zap.Error(err))
		return err
	}

	// 3. 根据帖子的社区id，将帖子id插入到社区对应的set中
	err = redisClient.SAdd(CommunityPrefix+strconv.FormatInt(post.CommunityID, 10), bytes).Err()
	if err != nil {
		zap.L().Error("SAdd failed", zap.Error(err))
		return err
	}

	zap.L().Info("CreatePost success", zap.Int64("post_id", post.PostID))

	return nil
}

// CommunitySortedPost 根据社区id，排序类型，指定范围，获取对应的数据
func CommunitySortedPost(cid string, p *models.PostListParam) ([]string, error) {
	switch p.Order {
	case "time":
		return redisClient.ZRevRange(PostTimeZSetKey+":"+cid, p.Offset, p.Offset+p.Limit).Result()
	case "vote":
		return redisClient.ZRange(PostVoteZSetKey+":"+cid, p.Offset, p.Offset+p.Limit).Result()
	}
	return nil, errors.New("排序有问题")
}

func ZsetExist(cid string, p *models.PostListParam) bool {
	switch p.Order {
	case "time":
		return redisClient.Exists(PostTimeZSetKey+":"+cid).Val() == 1
	case "vote":
		return redisClient.Exists(PostVoteZSetKey+":"+cid).Val() == 1
	}
	return false
}
