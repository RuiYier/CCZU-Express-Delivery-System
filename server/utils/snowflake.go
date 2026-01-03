package utils

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func InitSnowflake(nodeId int64) error {
	snowflake.Epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6
	var err error
	Node, err = snowflake.NewNode(nodeId)
	return err
}
