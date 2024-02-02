package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_web_app/dao/mysql"
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
	err = redisClient.ZAdd(PostTimeZSetKey, redis.Z{Score: float64(post.CreateAt.Unix()), Member: bytes}).Err()
	if err != nil {
		zap.L().Error("ZAdd (PostTimeZSet) failed", zap.Error(err))
		return err
	}

	// 2. 使用当前时间作为帖子的初始score插入到post:score zset中
	err = redisClient.ZAdd(PostVoteZSetKey, redis.Z{Score: float64(post.Score), Member: bytes}).Err()
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

// GetPostList 获取帖子列表, 根据排序列表进行排序
func GetPostList(p *models.PostListParam) ([]string, error) {
	//1. 从p获取offset和limit
	start, end := p.Offset, p.Offset+p.Limit

	//2. 从p 获取order顺序从而决定使用time 还是score, 涉及到zset的key 拼接
	idlist, err := redisClient.ZRevRange(PostPrefix+p.Order, start, end).Result()

	if err != nil {
		zap.L().Error("redisClient.ZRevRange(PostPrefix+\":\"+p.Order, start, end).Result() failed", zap.Error(err))
		return nil, err
	}

	return idlist, nil
}

// 返回每个帖子的投票总数，包括反对票和赞成票
//
// 根据postid 拼接查询zset3, key=post:voted:postid, 获取帖子的所有投票记录
//
// zset3 post:voted:postid 存储了这个帖子的所有投票记录
func GetPostVoteData(idlist []string) ([]int64, error) {
	//1. 开启pipeline，减少rtt
	// TxPipeline 就是用于这个目的，它自动在底层处理 MULTI 和 EXEC
	data := make([]int64, 0, len(idlist))
	pipeline := redisClient.TxPipeline()
	for i := range idlist {
		// 统计赞同和反对票数
		pipeline.ZCount(PostVotedPrefix+idlist[i], "-1", "1")
	}
	cmds, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return nil, err
	}

	// 2. 遍历cmd并且断言intcmd, 存储了zcount的结果
	for _, cmd := range cmds {
		if value, ok := cmd.(*redis.IntCmd); ok {
			data = append(data, value.Val())
		} else {
			zap.L().Error("pipeline.Exec() failed", zap.Error(err))
			return nil, err
		}
	}

	return data, nil
}

// GetPostAgreeVoteData 获取每个帖子的赞成票数
//
// 遍历idlist 获取每个帖子对应的zset key, post:voted:postid, 获取帖子的赞成票数
func GetPostAgreeVoteData(idlist []string) ([]int64, error) {
	// 1. 开启pipeline，减少rtt
	data := make([]int64, 0, len(idlist))
	pipeline := redisClient.TxPipeline()

	// 2. 遍历所有的postid 获取对应的赞成票
	for i := range idlist {
		pipeline.ZCount(PostVotedPrefix+idlist[i], "1", "1").Val()
	}
	cmds, err := pipeline.Exec()

	// 3. 遍历cmd并且断言intcmd
	for _, cmd := range cmds {
		if value, ok := cmd.(*redis.IntCmd); ok {
			data = append(data, value.Val())
		} else {
			zap.L().Error("pipeline.Exec() failed", zap.Error(err))
			return nil, err
		}
	}

	return data, nil
}

// GetPostDisagreeVoteData 获取每个帖子的反对票数
//
// 遍历idlist 获取每个帖子对应的zset key, post:voted:postid, 获取帖子的反对票数
func GetPostDisagreeVoteData(idlist []string) ([]int64, error) {
	// 1. 开启pipeline，减少rtt
	data := make([]int64, 0, len(idlist))
	pipeline := redisClient.TxPipeline()

	// 2. 遍历所有的postid 获取对应的反对票
	for i := range idlist {
		pipeline.ZCount(PostVotedPrefix+idlist[i], "-1", "-1").Val()
	}
	cmds, err := pipeline.Exec()

	// 3. 遍历cmd并且断言intcmd
	for _, cmd := range cmds {
		if value, ok := cmd.(*redis.IntCmd); ok {
			data = append(data, value.Val())
		} else {
			zap.L().Error("pipeline.Exec() failed", zap.Error(err))
			return nil, err
		}
	}

	return data, nil
}

func GetCommunitySortedPostIds(key string, p *models.PostListParam) ([]string, error) {
	// 1. 拼接key
	desKey := key + ":" + PostPrefix + p.Order
	// 2. 判断这个zset是否存在，如果不存在则和p.order 做交集，产生一个新的zset
	// 如果这个集合存在，则直接返回这个集合的所有元素
	// 3. 查询redis，根据score返回这个zset的所有元素id

	if cmd := redisClient.Exists(desKey); cmd.Val() == 0 {
		pipeline := redisClient.TxPipeline()
		pipeline.ZInterStore(desKey, redis.ZStore{Aggregate: "MAX"}, PostPrefix+p.Order, CommunityPrefix+key)
		pipeline.Expire(desKey, 60*time.Second) // 设置交集的key的过期时间
		pipeline.Exec()
	}

	// 根据offset和limit获取交集的key的所有元素
	return redisClient.ZRevRange(desKey, p.Offset, p.Offset+p.Limit).Result()
}
