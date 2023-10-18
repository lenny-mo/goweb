package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node // 当前节点
)

func Init(startTime string, machineId int64) (err error) {
	// 1. 转化成Time类型
	// 2006-01-02 15:04:05 是go语言诞生的时间
	start, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return err
	}
	// 2. Epoch 是 snowflake 包的一个公共变量，用于定义雪花算法中的自定义纪元时间（起点时间）
	snowflake.Epoch = start.UnixNano() / 1000000 // 把纳秒转换为毫秒

	// 3. 创建一个节点实例
	node, err = snowflake.NewNode(machineId) // 传入当前节点的ID
	if err != nil {
		return err
	}
	return nil
}

func GetId() int64 {
	return node.Generate().Int64()
}
