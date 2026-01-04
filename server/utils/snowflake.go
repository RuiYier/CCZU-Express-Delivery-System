package utils

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

// InitSnowflake 初始化全局 Snowflake 节点，nodeId 通常为节点编号（0-1023）
func InitSnowflake(nodeId int64) error {
	// 使用自定义纪元（毫秒）以保证 ID 空间和时间偏移
	snowflake.Epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6
	var err error
	Node, err = snowflake.NewNode(nodeId)
	return err
}

// GenerateID 生成一个新的 int64 类型的雪花 ID（线程安全）。
// 返回错误当全局节点未初始化时。
func GenerateID() (int64, error) {
	if Node == nil {
		return 0, errors.New("snowflake node not initialized")
	}
	id := Node.Generate()
	return id.Int64(), nil
}

// MustGenerateID 与 GenerateID 类似，但在未初始化时会 panic，适合在启动时确保已初始化的场景。
func MustGenerateID() int64 {
	if Node == nil {
		panic("snowflake node not initialized")
	}
	return Node.Generate().Int64()
}
