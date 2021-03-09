package zing

import (
	"errors"
	"fmt"
	_ "github.com/chuwt/zing/config"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/zing/receiver"
)

//
//// 管理所有用户
//type UserCtx struct {
//	ctx context.Context
//
//	userManager map[object.UserId]*userRuntime // 用户管理
//}
//
//func NewUserCtx(tower dataCenter.DataTower) UserCtx {
//	return UserCtx{
//		userManager: make(map[object.UserId]*userRuntime),
//		ctx:         context.Background(),
//		dataTower:   tower,
//	}
//}
//
//func (uc *UserCtx) Run() {
//	// 从数据库加载用户和策略
//	// 启动数据监控
//}
//
//// 订阅交易对
//func (uc *UserCtx) pub(symbol object.VtSymbol, strategyId object.StrategyId, event dataCenter.Event) error {
//	// 注册消息接收
//	uc.dataTower.RegisterEvent(uc.ctx, symbol, strategyId, event)
//	// 发布订阅，当收到消息时会推送到event中
//	return uc.dataTower.Pub(uc.ctx, symbol)
//}
//
//// 获取指定用户运行时
//func (uc *UserCtx) GetUser(userId object.UserId) *userRuntime {
//	if userRuntime, ok := uc.userManager[userId]; ok {
//		return userRuntime
//	}
//	return nil
//}
//
///*
//创建用户
//   |
//   |
//   v
//创建策略（初始化资源）——————>终止策略（释放资源）
//   |                        ^
//   |						|
//   |——————>启动策略—————————>|
//               |		    |
//               |———————> 暂停策略
//
//*/
//
//// 创建用户运行时
//func (uc *UserCtx) AddUser(userId object.UserId) error {
//	if user := uc.GetUser(userId); user == nil {
//		uc.userManager[userId] = NewUserRuntime(userId)
//		return nil
//	}
//	return errors.New("user existed")
//}
//
//// 添加用户gateway
//func (uc *UserCtx) AddUserGateway(userId object.UserId) error {
//	user := uc.GetUser(userId)
//	if user == nil {
//		return errors.New("user not existed")
//	}
//
//	//user.AddUserGateway()
//
//	return errors.New("user existed")
//}
//
//// 添加策略
//func (uc *UserCtx) AddUserStrategy(userId object.UserId, strategyId object.StrategyId, symbol object.VtSymbol, setting object.StrategySetting) error {
//	var err error
//
//	user := uc.GetUser(userId)
//	if user == nil {
//		return errors.New("user not existed")
//	}
//	if err = user.AddUserStrategy(strategyId, symbol, setting); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// 开启策略
//func (uc *UserCtx) StartUserStrategy(userId object.UserId, strategyId object.StrategyId) error {
//	// 开始订阅
//	user := uc.GetUser(userId)
//	if user == nil {
//		return errors.New("user not existed")
//	}
//	// 获取策略
//	strategy := user.GetUserStrategy(strategyId)
//	if strategy == nil {
//		return errors.New("user strategy not existed")
//	}
//
//	// 启动订阅
//	if err := uc.pub(strategy.symbol, strategyId, strategy.tick); err != nil {
//		return err
//	}
//	// 策略启动
//	go strategy.Run()
//	return nil
//}
//
//// 暂停策略
//func (uc *UserCtx) StopUserStrategy() {}
//
//// 移除策略
//func (uc *UserCtx) RemoveStrategy() {}
//
// 用户运行时 包括gateway管理和策略管理
type userRuntime struct {
	Id       object.UserId
	gateway  map[object.Gateway]gateway.UserGateway // 用户gateway列表，当前每个gateway的api只允许存在一个
	strategy map[object.StrategyId]*userStrategy    // 用户策略管理

	//event Event // 收到gateway推送的数据，然后推送到strategy的chan里
}

func NewUserRuntime(id object.UserId) *userRuntime {
	return &userRuntime{
		Id:       id,
		gateway:  make(map[object.Gateway]gateway.UserGateway),
		strategy: make(map[object.StrategyId]*userStrategy),
	}
}

// 添加用户的gateway
func (ur *userRuntime) AddUserGateway(userGateway gateway.UserGateway) error {
	if _, ok := ur.gateway[userGateway.Name()]; !ok {
		ur.gateway[userGateway.Name()] = userGateway
		return nil
	}
	// todo gateway的初始化
	return errors.New("gateway existed, remove before add")
}

// 添加策略
func (ur *userRuntime) AddUserStrategy(strategyId object.StrategyId, symbol object.VtSymbol, setting object.StrategySetting) error {
	if _, ok := ur.strategy[strategyId]; ok {
		return errors.New("strategy existed")
	}
	ur.strategy[strategyId] = NewUserStrategy(strategyId, symbol)
	return nil
}

// 获取策略
func (ur *userRuntime) GetUserStrategy(strategyId object.StrategyId) *userStrategy {
	if st, ok := ur.strategy[strategyId]; ok {
		return st
	}
	return nil
}

// 用户策略
type userStrategy struct {
	Id      object.StrategyId
	symbol  object.VtSymbol
	order   orderManager   // 用来管理订单，订单恢复等
	tick    *receiver.Tick // 获取数据
	runtime runtime        // 用来调用策略实例
}

type Sender interface {
	MarketSend()
	LimitSend()
	Cancel()
}

func NewUserStrategy(strategyId object.StrategyId, symbol object.VtSymbol) *userStrategy {
	return &userStrategy{
		symbol:  symbol,
		order:   orderManager{},
		tick:    receiver.NewTickReceiver(strategyId, 1024),
		runtime: runtime{},
	}
}

func (us *userStrategy) Run() {
	var tick object.TickData
	var counter int
	for {
		select {
		case tick = <-us.tick.Get():
			counter += 1
			if counter == 10 {
				fmt.Println("close!")
				us.tick.Close()
			}
			fmt.Println("接收到tick推送", tick, counter)
		}
	}
}

type orderManager struct {
	sender   Sender
	order    []OrderData
	trade    []TradeData
	position []Position
}

type OrderData struct{}
type TradeData struct{}
type Position struct{}

type runtime struct{}
