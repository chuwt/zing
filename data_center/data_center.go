package data_center

import (
	"context"
	"encoding/json"
	"github.com/chuwt/zing/client/redis"
	"github.com/chuwt/zing/object"
	pubsub "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sync"
)

var (
	Log = zap.L().With(zap.Namespace("dataTower"))
)

type DataTower interface {
	Init() error
	Pub(ctx context.Context, symbol object.VtSymbol) error
	//Sub(ctx context.Context, symbol object.VtSymbol) error
	RegisterEvent(ctx context.Context, symbol object.VtSymbol, strategyId object.StrategyId, event Event)
	//UnRegisterEvent(ctx context.Context, symbol object.VtSymbol)
}

type TowerManager struct {
	events sync.Map
}

func (tm *TowerManager) Init() error {
	return nil
}

func (tm *TowerManager) RegisterEvent(ctx context.Context, symbol object.VtSymbol, strategyId object.StrategyId, event Event) {
	if eventIn, ok := tm.events.LoadOrStore(symbol.String(), []Event{event}); ok {
		events := eventIn.([]Event)
		events = append(events, event)
		tm.events.Store(symbol.String(), events)
	}
}

func (tm *TowerManager) Pub(ctx context.Context, symbol object.VtSymbol) error {
	return nil
}

type RedisTower struct {
	TowerManager
	redis     *redis.Redis
	subSymbol map[object.VtSymbol]struct{} // 记录订阅的交易对
	mu        sync.Mutex
}

func NewRedisTower(cfg redis.Config) *RedisTower {
	return &RedisTower{
		TowerManager: TowerManager{
			events: sync.Map{},
		},
		redis:     redis.NewRedis(cfg),
		subSymbol: make(map[object.VtSymbol]struct{}),
		mu:        sync.Mutex{},
	}
}

func (rt *RedisTower) Pub(ctx context.Context, symbol object.VtSymbol) error {
	Log.Info("开始添加tick订阅", zap.String("vtSymbol", symbol.String()))

	rt.mu.Lock()
	defer rt.mu.Unlock()
	if _, ok := rt.subSymbol[symbol]; ok {
		Log.Warn("交易对已订阅", zap.String("vtSymbol", symbol.String()))
		return nil
	}
	rt.subSymbol[symbol] = struct{}{}
	if err := rt.redis.Publish(ctx, "subscribe_symbol", symbol.String()).Err(); err != nil {
		return err
	}
	go rt.Sub(ctx, symbol)
	return nil
}

func (rt *RedisTower) Sub(ctx context.Context, symbol object.VtSymbol) {
	var (
		err  error
		msg  *pubsub.Message
		tick *object.TickData
	)
retry:
	pubSub := rt.redis.Subscribe(ctx, symbol.String())
	for {
		msg, err = pubSub.ReceiveMessage(ctx)
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

		eventsIn, ok := rt.events.Load(symbol.String())
		if !ok {
			continue
		}

		j := 0
		events := eventsIn.([]Event)

		for i, event := range events {
			if !event.IsClosed() {
				event.OnReceive() <- *tick
			} else {
				events[j], events[i] = events[i], events[j]
				j += 1
			}
		}
		if j >= len(events) {
			rt.events.Store(symbol, events[0:0])
		}
		rt.events.Store(symbol, events[j:])
	}
}

type Event interface {
	OnReceive() chan object.TickData
	IsClosed() bool
}
