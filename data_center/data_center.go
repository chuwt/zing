package data_center

import (
	"context"
	"github.com/chuwt/zing/object"
	"go.uber.org/zap"
	"sync"
)

var (
	Log = zap.L().With(zap.Namespace("dataTower"))
)

type DataTower interface {
	Init(ctx context.Context) error
	// 注册
	AddReceiver(ctx context.Context, symbol object.VtSymbol, receiver Receiver) error
	// 订阅symbol
	pub(ctx context.Context, symbol object.VtSymbol) error
	// 接收tick
	sub(ctx context.Context, symbol object.VtSymbol)
}

type Receiver interface {
	Id() object.StrategyId
	OnReceive() chan object.TickData
	IsClosed() bool
	IsFull() bool
}

type TowerManager struct {
	mu sync.Mutex
	// symbol的列表，用于数据推送
	symbolLink map[object.VtSymbol]*ReceiverLink
	// receiver的map，用于快速查询receiver是否已存在
	receiverMap map[object.StrategyId]*Receiver
}

type ReceiverLink struct {
	Receiver Receiver
	Next     *ReceiverLink
}

func (tm *TowerManager) Init(ctx context.Context) error {
	return nil
}

func (tm *TowerManager) AddReceiver(ctx context.Context, symbol object.VtSymbol, receiver Receiver) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if _, ok := tm.receiverMap[receiver.Id()]; ok {
		Log.Warn("receiver已存在", zap.Int64("id", int64(receiver.Id())))
		return nil
	}
	if lastReceiver, ok := tm.symbolLink[symbol]; !ok {
		tm.symbolLink[symbol] = &ReceiverLink{
			Receiver: receiver,
			Next:     nil,
		}
	} else {
		tm.symbolLink[symbol] = &ReceiverLink{
			Receiver: receiver,
			Next:     lastReceiver,
		}
	}
	return nil
}

func (tm *TowerManager) pub(ctx context.Context, symbol object.VtSymbol) error { return nil }

func (tm *TowerManager) sub(ctx context.Context, symbol object.VtSymbol) {}
