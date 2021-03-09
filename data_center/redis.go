package data_center

import (
	"context"
	"github.com/chuwt/zing/client/redis"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	pubsub "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sync"
)

type RedisTower struct {
	TowerManager
	redis     *redis.Redis
	subSymbol map[object.VtSymbol]struct{} // 记录订阅的交易对
}

func NewRedisTower(cfg redis.Config) *RedisTower {
	return &RedisTower{
		TowerManager: TowerManager{
			mu:          sync.Mutex{},
			symbolLink:  make(map[object.VtSymbol]*ReceiverLink),
			receiverMap: make(map[object.StrategyId]*Receiver),
		},
		redis:     redis.NewRedis(cfg),
		subSymbol: make(map[object.VtSymbol]struct{}),
	}
}

func (rt *RedisTower) pub(ctx context.Context, symbol object.VtSymbol) error {
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
	go rt.sub(ctx, symbol)
	return nil
}

func (rt *RedisTower) AddReceiver(ctx context.Context, symbol object.VtSymbol, receiver Receiver) error {
	_ = rt.TowerManager.AddReceiver(ctx, symbol, receiver)
	return rt.pub(ctx, symbol)
}

func (rt *RedisTower) sub(ctx context.Context, symbol object.VtSymbol) {
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

		err = json.Json.Unmarshal([]byte(msg.Payload), tick)
		if err != nil {
			Log.Error("解析tick失败", zap.Error(err))
			continue
		}

		link, ok := rt.symbolLink[symbol]
		if !ok {
			Log.Error("symbol 未订阅", zap.String("symbol", symbol.String()))
			continue
		}

		var prev *ReceiverLink
		for {
			if link == nil {
				break
			}
			if link.Receiver.IsClosed() {
				// 实时删除需要锁的支持，后面可以考虑异步删除
				rt.mu.Lock()
				// 已关闭，删除
				// 截断链表
				if prev == nil {
					rt.symbolLink[symbol] = link.Next
				} else {
					prev.Next = link.Next
				}

				// 删除map
				delete(rt.receiverMap, link.Receiver.Id())
				rt.mu.Unlock()
			} else {
				// 在receiver的chan满时不推送，防止阻塞
				if link.Receiver.IsFull() {
					Log.Warn("receiver full", zap.Int64("id", int64(link.Receiver.Id())))
				} else {
					link.Receiver.OnReceive() <- *tick
				}
			}
			prev = link
			link = link.Next
		}
	}
}
