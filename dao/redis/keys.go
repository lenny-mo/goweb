package redis

// redis key 使用命名空间方式进行区分
const (
	PostTimeZSetKey = "post:time"   // zset; 发帖时间作为分数
	PostVoteZSetKey = "post:vote"   // zset; 投票作为分数
	PostVotedPrefix = "post:voted:" // set; 记录用户及投票类型, 需要通过拼接postid使用
)
