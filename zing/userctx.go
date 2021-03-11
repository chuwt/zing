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
	Id            object.UserId
	strategy      map[object.StrategyId]*userStrategy // 用户策略管理
	gateway       map[object.Gateway]*runtimeGateway  // 用户gateway列表，当前每个gateway的api只允许存在一个
	eventReceiver chan object.Event                   // 收到gateway推送的数据，然后推送到strategy的chan里
}

type runtimeGateway struct {
	g           gateway.UserGateway
	position    map[object.Currency]object.PositionData
	order       map[object.ClientOrderId]*orderStrategy
	orderClient map[object.OrderId]object.ClientOrderId
	trade       map[object.TradeId]object.TradeData
}

func (rg *runtimeGateway) SetPosition(currency object.Currency, position object.PositionData) {
	rg.position[currency] = position
}

// 根据 clientOrderId 获取 orderStrategy
func (rg *runtimeGateway) GetOrderStrategy(id object.ClientOrderId) *orderStrategy {
	order, ok := rg.order[id]
	if !ok {
		return nil
	}
	return order
}

// 设置 orderId 和 clientOrderId 的对应
func (rg *runtimeGateway) SetOrderClient(id object.OrderId, orderId object.ClientOrderId) {
	rg.orderClient[id] = orderId
}

// 根据 orderId 获取 clientOrderId
func (rg *runtimeGateway) GetClientOrderId(id object.OrderId) object.ClientOrderId {
	clientOrderId, ok := rg.orderClient[id]
	if !ok {
		return ""
	}
	return clientOrderId
}

// 通过 tradeId 检查trade是否已存在
func (rg *runtimeGateway) ExistedTrade(id object.TradeId) bool {
	_, ok := rg.trade[id]
	if !ok {
		return false
	}
	return true
}

type orderStrategy struct {
	StrategyId object.StrategyId
	OrderData  object.OrderData
}

func NewRuntimeGateway(gateway gateway.UserGateway) *runtimeGateway {
	return &runtimeGateway{
		g:        gateway,
		position: make(map[object.Currency]object.PositionData),
		order:    make(map[object.ClientOrderId]*orderStrategy),
		trade:    make(map[object.TradeId]object.TradeData),
	}
}

func NewUserRuntime(id object.UserId) *userRuntime {
	ur := &userRuntime{
		Id:            id,
		gateway:       make(map[object.Gateway]*runtimeGateway),
		strategy:      make(map[object.StrategyId]*userStrategy),
		eventReceiver: make(chan object.Event, 1024),
		//position:      make(map[object.VtCurrency]object.PositionData),
		//orderStrategy: make(map[object.ClientOrderId]object.StrategyId),
		//order:         make(map[object.ClientOrderId]object.OrderData),
		//trade:         make(map[object.TradeId]object.TradeData),
	}
	go ur.Loop()
	return ur
}

func (ur *userRuntime) GetRuntimeGateway(gateway object.Gateway) *runtimeGateway {
	if g, ok := ur.gateway[gateway]; ok {
		return g
	}
	return nil
}

func (ur *userRuntime) Loop() {
	for {
		select {
		case event := <-ur.eventReceiver:
			/*
				order订单通知后，存储order的记录
				trade订单通知后，存储trade的记录
			*/

			g := ur.GetRuntimeGateway(event.Gateway)
			if g == nil {
				Log.Error("gateway not exited", zap.Any("data", event))
				continue
			}

			switch event.Type {
			case object.EventTypeConnection:
				// todo 新链接建立的时候，需要重新检查订单状态
			case object.EventTypePosition:
				// 用户当前金额
				position := event.Data.(object.PositionData)
				g.SetPosition(object.Currency(position.Currency), position)
			case object.EventTypeOrder:
				// 用户订单
				order := event.Data.(object.OrderData)
				existedOrder := g.GetOrderStrategy(object.ClientOrderId(order.ClientOrderId))
				if existedOrder == nil {
					Log.Error("order not existed", zap.Any("order", order))
					break
				}
				if existedOrder.OrderData.Status == order.Status &&
					existedOrder.OrderData.ExecAmt == order.ExecAmt {
					Log.Warn("existed order", zap.Any("order", order))
					break
				}
				existedOrder.OrderData = order
				if order.OrderId != 0 {
					g.SetOrderClient(object.OrderId(order.OrderId), object.ClientOrderId(order.ClientOrderId))
				}
				strategy, ok := ur.strategy[existedOrder.StrategyId]
				if !ok {
					Log.Error("strategy not existed", zap.Any("order", existedOrder))
					break
				}

				// 更新订单并发送通知到策略
				strategy.order <- order
				//strategy.OnOrder(order)
				// todo 入库
				// todo 此时的入库只是创建记录，策略消费后需要更新记录状态，标记消费了

			case object.EventTypeTrade:
				// todo 用户成交单
				trade := event.Data.(object.TradeData)
				clientOrderId := g.GetClientOrderId(object.OrderId(trade.OrderId))
				if clientOrderId == "" {
					Log.Warn("未找到成交单的订单Id", zap.Any("trade", trade))
					break
				}
				// todo 此时检测到只能说明记录了，是否被消费还得等待策略标记
				if g.ExistedTrade(object.TradeId(trade.TradeId)) {
					// 成交已存在
					Log.Warn("existed tradeId", zap.Any("trade", trade))
					break
				}

				existedOrder := g.GetOrderStrategy(object.ClientOrderId(clientOrderId))
				if existedOrder == nil {
					Log.Error("order not existed", zap.Any("trade", trade))
					break
				}

				strategy, ok := ur.strategy[existedOrder.StrategyId]
				if !ok {
					Log.Error("strategy not existed", zap.Any("order", existedOrder))
					break
				}
				// 更新成交并发送通知到策略
				//strategy.OnTrade(trade)
				strategy.trade <- trade
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
		ur.gateway[userGateway.Name()] = NewRuntimeGateway(userGateway)
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
	// todo 将策略的订单倒入
	rg, ok := ur.gateway[symbol.GatewayName]
	if !ok {
		return errors.New("gateway not exist")
	}
	ur.strategy[strategyId] = NewUserStrategy(strategyId, symbol, rg)
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
	gateway *runtimeGateway

	tick    *receiver.Tick // 获取数据
	order   chan object.OrderData
	trade   chan object.TradeData
	runtime runtime // 用来调用策略实例
}

func NewUserStrategy(strategyId object.StrategyId, symbol object.VtSymbol, gateway *runtimeGateway) *userStrategy {
	return &userStrategy{
		Id:      strategyId,
		symbol:  symbol,
		gateway: gateway,
		tick:    receiver.NewTickReceiver(strategyId, 1024),
		runtime: nil,
	}
}

func (us *userStrategy) Run() {
	var tick object.TickData
	for {
		select {
		case tick = <-us.tick.Get():
			// 这一步是同步的
			//us.runtime.OnTick(tick)
			fmt.Println(tick)
		case order := <-us.order:
			fmt.Println(order)
		case trade := <-us.trade:
			fmt.Println(trade)
		}
	}
}

func (us *userStrategy) Init() {
	// 初始化
}

type runtime interface {
	OnTick(object.TickData)
	OnOrder(object.OrderData)
	OnTrade(object.TradeData)
}
