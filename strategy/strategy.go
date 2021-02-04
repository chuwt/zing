package strategy

import "vngo/object"

type Strategy struct {
	userCtx map[object.UserId]*userCtx
}

func (s *Strategy) Init() {
	// 订阅发布连接
}

func (s *Strategy) AddStrategy(userId object.UserId, strategyId object.StrategyId) {
	// 订阅交易对
}

func (s *Strategy) OnTick() {

}

type userCtx struct {
	userId     object.UserId
	strategies map[object.StrategyId]*UserStrategy
}

type UserStrategy struct {
	VtSymbol object.VtSymbol
	Setting  string
}
