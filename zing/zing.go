package zing

import (
	"context"
	"errors"
	dataCenter "github.com/chuwt/zing/data_center"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/object"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"sync"
)

var (
	Log = zap.L().With(zap.Namespace("zing"))
)

type Zing struct {
	single singleflight.Group
	mu     sync.Mutex

	// 策略引擎，可以使用python或者go
	strategyEngine StrategyEngine

	// 用户管理
	userRuntime map[object.UserId]*userRuntime

	// tick数据中心
	dataTower dataCenter.DataTower
}

func NewZing(strategyEngine StrategyEngine, dataTower dataCenter.DataTower) Zing {
	return Zing{
		single:         singleflight.Group{},
		mu:             sync.Mutex{},
		strategyEngine: strategyEngine,
		userRuntime:    make(map[object.UserId]*userRuntime),
		dataTower:      dataTower,
	}
}

// 添加用户
func (z *Zing) AddUser(userId object.UserId) error {
	if _, ok := z.userRuntime[userId]; ok {
		Log.Error("user existed", zap.String("userId", string(userId)))
		return errors.New("user existed")
	}
	z.userRuntime[userId] = NewUserRuntime(userId)
	return nil
}

func (z *Zing) GetUser(userId object.UserId) *userRuntime {
	if ur, ok := z.userRuntime[userId]; ok {
		return ur
	}
	return nil
}

// 连接用户交易所信息
func (z *Zing) AddUserGateway(userId object.UserId, userGateway gateway.UserGateway) error {
	user := z.GetUser(userId)
	if user == nil {
		return errors.New("user not existed")
	}

	return user.AddUserGateway(userGateway)
}

// 添加策略
func (z *Zing) AdduserStrategy(userId object.UserId, strategyId object.StrategyId, symbol object.VtSymbol, setting object.StrategySetting) error {
	user := z.GetUser(userId)
	if user == nil {
		return errors.New("user not existed")
	}
	return user.AddUserStrategy(strategyId, symbol, setting)
}

// 启动策略
func (z *Zing) RunUserStrategy(userId object.UserId, strategyId object.StrategyId) error {
	user := z.GetUser(userId)
	if user == nil {
		return errors.New("user not existed")
	}
	strategy := user.GetUserStrategy(strategyId)
	if strategy == nil {
		return errors.New("strategy not existed")
	}
	z.dataTower.AddReceiver(context.Background(), strategy.symbol, strategy.tick)
	go strategy.Run()
	return nil
}

type StrategyEngine interface{}
