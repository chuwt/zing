package zing

import (
	"errors"
	"fmt"
	_ "github.com/chuwt/zing/config"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/zing/receiver"
	"go.uber.org/zap"
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

	eventReceiver chan object.Event // 收到gateway推送的数据，然后推送到strategy的chan里

	position      map[object.VtCurrency]object.PositionData
	orderStrategy map[object.ClientOrderId]object.StrategyId
	order         map[object.ClientOrderId]object.OrderData
	trade         map[object.TradeId]object.TradeData
}

func NewUserRuntime(id object.UserId) *userRuntime {
	ur := &userRuntime{
		Id:            id,
		gateway:       make(map[object.Gateway]gateway.UserGateway),
		strategy:      make(map[object.StrategyId]*userStrategy),
		eventReceiver: make(chan object.Event, 1024),
		position:      make(map[object.VtCurrency]object.PositionData),
		orderStrategy: make(map[object.ClientOrderId]object.StrategyId),
		order:         make(map[object.ClientOrderId]object.OrderData),
		trade:         make(map[object.TradeId]object.TradeData),
	}
	go ur.Loop()
	return ur
}

func (ur *userRuntime) Loop() {
	for {
		select {
		case event := <-ur.eventReceiver:
			/*
				order订单通知后，存储order的记录
				trade订单通知后，存储trade的记录
			*/
			switch event.Type {
			case object.EventTypeConnection:
				// todo 新链接建立的时候，需要重新检查订单状态
			case object.EventTypePosition:
				// todo 用户当前金额
				position := event.Data.(object.PositionData)
				ur.position[position.VtCurrency()] = position
			case object.EventTypeOrder:
				// 用户订单
				order := event.Data.(object.OrderData)
				if existedOrder, ok := ur.order[object.ClientOrderId(order.ClientOrderId)]; ok {
					if existedOrder.Status == order.Status &&
						existedOrder.ExecAmt == order.ExecAmt {
						Log.Warn("existed order", zap.Any("order", order))
						break
					}
					// todo 更新订单并发送通知到策略
					// todo 此时的入库只是创建记录，策略消费后需要更新记录状态，标记消费了
				}

			case object.EventTypeTrade:
				// todo 用户成交单
				trade := event.Data.(object.TradeData)
				if _, ok := ur.trade[object.TradeId(trade.TradeId)]; ok {
					Log.Warn("existed tradeId", zap.Any("trade", trade))
					break
				}
				// todo 更新成交并发送通知到策略
				// todo 此时的入库只是创建记录，策略消费后需要更新记录状态，标记消费了
			default:
				continue
			}
		}
	}
}

// 添加用户的gateway
func (ur *userRuntime) AddUserGateway(userGateway gateway.UserGateway) error {
	if _, ok := ur.gateway[userGateway.Name()]; !ok {
		ur.gateway[userGateway.Name()] = userGateway
		// gateway启动
		go userGateway.Start(ur.eventReceiver)
		return nil
	}
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
	Id     object.StrategyId
	symbol object.VtSymbol

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

type runtime struct{}
