package db

import (
	"fmt"
	"vngo/object"
)

type Strategy struct {
	Base       `xorm:"extends"`
	UserId     string
	StrategyId int64 `xorm:"unique"`
	Symbol     string
	Gateway    string
	Status     int32 // 0 未启动 1 运行 -1 失败 2 暂停 3 停止
}

func (s *Strategy) StrategyKey() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%d", s.UserId, s.StrategyId))
}
