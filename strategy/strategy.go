package strategy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pubsub "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"sync"
	"vngo/client/redis"
	"vngo/db"
	"vngo/object"
)

var Log = zap.L().With(zap.Namespace("strategy"))

type Strategy struct {
	ctx   context.Context
	redis *redis.Redis
	once  singleflight.Group
	mu    sync.Mutex

	userStrategy UserStrategy
	subMap       map[object.VtSymbol]map[object.StrategyId]*st
}

type tickSub struct {
	vtSymbol object.VtSymbol
}

type st struct {
}

func (*st) OnTick(tick *object.TickData) {}

func NewStrategy(cfg redis.Config) Strategy {
	return Strategy{
		ctx:          context.Background(),
		redis:        redis.NewRedis(cfg),
		userStrategy: NewUserStrategy(),
		once:         singleflight.Group{},
		mu:           sync.Mutex{},
		subMap:       make(map[object.VtSymbol]map[object.StrategyId]*st),
	}
}

func (s *Strategy) Init() {
	// todo
	// 读库
	// 将所有策略的交易对进行订阅
}

type AddStrategyReq struct {
	UserId     object.UserId
	StrategyId object.StrategyId
	VtSymbol   object.VtSymbol
	Setting    string
}

func (as *AddStrategyReq) Key() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%d", as.UserId, as.StrategyId))
}

func (s *Strategy) AddStrategy(strategy AddStrategyReq) error {

	var err error
	key := strategy.Key()

	Log.Info("添加策略", zap.String("key", string(key)))

	// 同一用户的同一个策略做过滤
	_, err, _ = s.once.Do(string(key), func() (interface{}, error) {
		dbSt := &db.Strategy{
			UserId:     string(strategy.UserId),
			StrategyId: int64(strategy.StrategyId),
			Symbol:     strategy.VtSymbol.Symbol,
			Gateway:    string(strategy.VtSymbol.GatewayName),
			Status:     0,
		}

		// 判断本地是否存在
		if err = s.userStrategy.AddStrategy(dbSt); err != nil {
			Log.Error("策略已存在", zap.String("key", string(key)))
			return nil, err
		}

		// 入库
		duplicate, err := db.CreateDupEntry(db.GetEngine(), dbSt)
		if err != nil {
			Log.Error("策略入库失败", zap.String("key", string(key)), zap.Error(err))
			s.userStrategy.RemoveStrategy(key)
			return nil, err
		} else if duplicate {
			Log.Warn("策略已存在", zap.String("key", string(key)))
		}

		// 订阅交易对信息
		return nil, s.Sub(strategy.VtSymbol, strategy.StrategyId)
	})

	Log.Info("添加策略成功", zap.String("key", string(key)))

	return err
}

func (s *Strategy) RemoveStrategy(userId object.UserId, strategyId object.StrategyId) {

}

func (s *Strategy) pub(symbol object.VtSymbol) error {
	return s.redis.Publish(s.ctx, "subscribe_symbol", symbol.String()).Err()
}

func (s *Strategy) Sub(symbol object.VtSymbol, strategyId object.StrategyId) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subMap[symbol]; !ok {
		if err := s.pub(symbol); err != nil {
			return err
		}
		s.subMap[symbol] = make(map[object.StrategyId]*st)
		s.subMap[symbol][strategyId] = &st{}
	} else if _, ok := s.subMap[symbol][strategyId]; !ok {
		s.subMap[symbol][strategyId] = &st{}
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
			userStrategy.OnTick(tick)
		}
	}
}

func (s *Strategy) OnTick() {

}

type UserStrategy struct {
	strategies map[object.StrategyKey]*db.Strategy
}

func NewUserStrategy() UserStrategy {
	return UserStrategy{
		strategies: make(map[object.StrategyKey]*db.Strategy),
	}
}

func (us *UserStrategy) AddStrategy(strategy *db.Strategy) error {
	if _, ok := us.strategies[strategy.StrategyKey()]; !ok {
		us.strategies[strategy.StrategyKey()] = strategy
		return nil
	} else {
		return errors.New("strategy existed")
	}
}

func (us *UserStrategy) RemoveStrategy(strategyKey object.StrategyKey) {
	delete(us.strategies, strategyKey)
}
