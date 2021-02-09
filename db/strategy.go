package db

import (
	"fmt"
	"vngo/object"
	"xorm.io/xorm"
)

type Strategy struct {
	Base           `xorm:"extends"`
	UserId         string
	StrategyId     int64 `xorm:"unique"`
	Symbol         string
	Gateway        string
	LoadBar        int
	Contract       bool
	RuntimeSetting string
	ClassName      string
	Status         int32 // 0 未启动 1 运行 -1 失败 2 暂停 3 停止
}

const (
	StrategyStatusNormal  = 0
	StrategyStatusRunning = 1
	StrategyStatusStopped = 2
	StrategyStatusExit    = 3
	StrategyStatusFailed  = -1
)

func (s *Strategy) StrategyKey() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%d", s.UserId, s.StrategyId))
}

func GetStrategyByKey(session xorm.Interface, userId object.UserId, strategyId object.StrategyId) (*Strategy, error) {
	st := new(Strategy)
	ok, err := session.Where("user_id=? and strategy_id=?", userId, strategyId).Get(st)
	if !ok {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return st, nil
}
