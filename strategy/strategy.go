package strategy

import (
	"context"
	"encoding/json"
	"errors"
	pubsub "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"sync"
	"vngo/client/redis"
	"vngo/config"
	"vngo/db"
	"vngo/object"
	"vngo/python"
)

var Log = zap.L().With(zap.Namespace("strategy"))

type Strategy struct {
	once  singleflight.Group
	mu    sync.Mutex
	ctx   context.Context
	redis *redis.Redis

	userStrategy UserStrategy

	pyEngine python.PyEngine
	// todo 这里可以做个interface，支持更多的订阅发布形式，如本地ws和远程redis
	subMap map[object.VtSymbol]map[object.StrategyKey]*strategyEntity
}

func NewStrategy(redisCfg redis.Config, strategyCfg config.Strategy) Strategy {
	return Strategy{
		once:         singleflight.Group{},
		mu:           sync.Mutex{},
		ctx:          context.Background(),
		redis:        redis.NewRedis(redisCfg),
		userStrategy: NewUserStrategy(),
		pyEngine:     python.NewPyEngine(strategyCfg.Path, strategyCfg.PythonPath),
		subMap:       make(map[object.VtSymbol]map[object.StrategyKey]*strategyEntity),
	}
}

func (s *Strategy) Init() error {
	// todo
	// 读库
	// 将所有策略的交易对进行订阅
	if err := s.pyEngine.Init(); err != nil {
		return err
	}
	return nil
}

func (s *Strategy) pub(symbol object.VtSymbol) error {
	return s.redis.Publish(s.ctx, "subscribe_symbol", symbol.String()).Err()
}

// todo 订阅的取消
func (s *Strategy) Sub(symbol object.VtSymbol, entity *strategyEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subMap[symbol]; !ok {
		if err := s.pub(symbol); err != nil {
			return err
		}
		s.subMap[symbol] = make(map[object.StrategyKey]*strategyEntity)
		s.subMap[symbol][entity.Data.StrategyKey()] = entity
	} else if _, ok := s.subMap[symbol][entity.Data.StrategyKey()]; !ok {
		s.subMap[symbol][entity.Data.StrategyKey()] = entity
		return nil
	} else {
		return nil
	}
	go s.sub(symbol)
	return nil
}

func (s *Strategy) sub(symbol object.VtSymbol) {
	var (
		err  error
		msg  *pubsub.Message
		tick *object.TickData
	)
retry:
	pubSub := s.redis.Subscribe(s.ctx, symbol.String())
	for {
		msg, err = pubSub.ReceiveMessage(s.ctx)
		if err != nil {
			Log.Error("接收订阅消息失败", zap.Error(err))
			_ = pubSub.Close()
			goto retry
		}

		Log.Debug("接收订阅消息", zap.String("tick", msg.Payload))

		tick = new(object.TickData)

		err = json.Unmarshal([]byte(msg.Payload), tick)
		if err != nil {
			Log.Error("解析tick失败", zap.Error(err))
			continue
		}

		for _, userStrategy := range s.subMap[symbol] {
			if userStrategy.Data.Status == db.StrategyStatusRunning {
				// todo 考虑waitGroup+timeout的并发执行
				userStrategy.Runtime.OnTick(tick)
			}
		}
	}
}

func (s *Strategy) OnTick() {

}

type UserStrategy struct {
	strategies map[object.StrategyKey]*strategyEntity
}

type strategyEntity struct {
	Data    *db.Strategy // 数据
	Runtime runtime      `xorm:"-"` // 策略运行时 python 或 golang
}

type runtime interface {
	OnTick(*object.TickData) error
	OnBar(*object.BarData) error
	OnContract(*object.ContractData) error
}

func NewUserStrategy() UserStrategy {
	return UserStrategy{
		strategies: make(map[object.StrategyKey]*strategyEntity),
	}
}

func (us *UserStrategy) AddStrategy(strategy *strategyEntity) error {
	if _, ok := us.strategies[strategy.Data.StrategyKey()]; !ok {
		us.strategies[strategy.Data.StrategyKey()] = strategy
		return nil
	} else {
		return errors.New("strategy existed")
	}
}

func (us *UserStrategy) GetStrategyByKey(key object.StrategyKey) *strategyEntity {
	if s, ok := us.strategies[key]; ok {
		return s
	} else {
		return nil
	}
}

func (us *UserStrategy) RemoveStrategy(strategyKey object.StrategyKey) {
	delete(us.strategies, strategyKey)
}
