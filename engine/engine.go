package engine

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"sync"
	"github.com/chuwt/zing/config"
	"github.com/chuwt/zing/db"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/python"
)

var (
	Log = zap.L().With(zap.Namespace("engine"))
)

type Engine struct {
	single   singleflight.Group
	mu       sync.Mutex
	ctx      context.Context
	pyEngine python.PyEngine

	gatewayFactory gateway.Factory
	userCtx        userCtx
	pubSub         PubSub
}

func NewEngine(gf gateway.Factory, strategyCfg config.Strategy) Engine {
	return Engine{
		single:         singleflight.Group{},
		mu:             sync.Mutex{},
		ctx:            context.Background(),
		pyEngine:       python.NewPyEngine(strategyCfg.Path, strategyCfg.PythonPath),
		gatewayFactory: gf,
		userCtx:        make(userCtx),
		pubSub:         nil,
	}
}

func (e *Engine) Init() error {
	// todo
	// 读库
	// 将所有策略的交易对进行订阅
	if err := e.pyEngine.Init(); err != nil {
		return err
	}
	return nil
}

/*
创建用户
   |
   |
   v
创建策略（初始化资源）——————>终止策略（释放资源）
   |                        ^
   |						|
   |——————>启动策略—————————>|
               |		    |
               |———————> 暂停策略

*/

// 创建用户
func (e *Engine) InitUser(id object.UserId) error {

	key := fmt.Sprintf("init user.%s", id)
	Log.Info("init user", zap.String("key", key))

	_, err, _ := e.single.Do(key, func() (interface{}, error) {
		if _, ok := e.userCtx[id]; ok {
			Log.Error("user existed", zap.String("key", key))
			return nil, errors.New("user existed")
		}
		// 创建用户上下文
		e.userCtx[id] = NewUserContext(id)
		Log.Info("init user success", zap.String("key", key))
		return nil, nil
	})
	return err
}

// 添加用户gateway
// todo 目前不支持一个用户有多个相同的gateway，后面可以通过hash做多个支持
func (e *Engine) AddUserGateway(req AddUserGatewayReq) error {
	var (
		err         error
		duplicate   bool
		key         = req.Key()
		userGateway gateway.UserGateway
	)

	Log.Info("add user userGateway", zap.String("key", key.String()))

	userCtx := e.userCtx.GetUserContext(req.UserId)
	if userCtx == nil {
		err = errors.New("user not existed")
		Log.Error("add user userGateway error",
			zap.String("key", key.String()),
			zap.Error(err))
		return err
	}

	_, err, _ = e.single.Do(key.String(), func() (interface{}, error) {
		// 入库
		session := db.GetSession()
		_ = session.Begin()
		defer session.Rollback()
		duplicate, err = db.CreateDupEntry(session, db.ApiAuth{
			UserId:      string(req.UserId),
			Gateway:     string(req.Gateway),
			ApiAuthJson: string(req.Auth),
		})
		if duplicate {
			Log.Warn("user api existed", zap.String("key", key.String()))
		} else if err != nil {
			Log.Error("add userGateway to db error",
				zap.String("key", key.String()),
				zap.Error(err))
			return nil, err
		}

		// userGateway 是否已存在
		userGateway = userCtx.GetUserGateway(req.Gateway)
		if userGateway != nil {
			err = errors.New("userGateway existed")
			Log.Error("add user userGateway error",
				zap.String("key", key.String()),
				zap.Error(err))
			return nil, err
		}

		// 新建gateway
		userGateway, err = e.gatewayFactory.NewUserGateway(req.UserId, req.Gateway, req.Auth)
		if err != nil {
			Log.Error(
				"new userGateway error",
				zap.String("key", key.String()),
				zap.Error(err))
			return nil, err
		}

		_ = userGateway.Init()

		// 添加到用户ctx
		userCtx.AddUserGateway(userGateway)

		Log.Info("add user userGateway success", zap.String("key", key.String()))

		_ = session.Commit()

		return nil, nil
	})
	return err
}

// 添加策略
func (e *Engine) AddStrategy(req AddStrategyReq) error {
	var (
		err error
		key = req.Key()
	)

	Log.Info("开始添加策略", zap.String("key", string(key)))

	_, err, _ = e.single.Do(key.String(), func() (interface{}, error) {

		//dbSt := &db.Strategy{
		//	UserId:         string(req.UserId),
		//	StrategyId:     int64(req.StrategyId),
		//	Symbol:         req.VtSymbol.Symbol,
		//	Gateway:        string(req.VtSymbol.GatewayName),
		//	LoadBar:        req.Setting.LoadBar,
		//	Contract:       req.Setting.Contract,
		//	RuntimeSetting: req.Setting.RuntimeSetting,
		//	ClassName:      req.StrategyClassName,
		//	Status:         0,
		//}
		//
		//// 入库
		//duplicate, err := db.CreateDupEntry(db.GetEngine(), dbSt)
		//if err != nil {
		//	Log.Error("策略入库失败", zap.String("key", string(key)), zap.Error(err))
		//	return nil, err
		//} else if duplicate {
		//	Log.Warn("策略已存在", zap.String("key", string(key)))
		//}
		//
		//// 创建策略运行时
		//runtime, err := s.pyEngine.NewStrategyInstance2(
		//	req.StrategyClassName,
		//	req.StrategyId,
		//	req.VtSymbol,
		//	req.Setting.RuntimeSetting)
		//if err != nil {
		//	Log.Error("添加策略失败: 创建策略运行时失败",
		//		zap.String("key", string(key)),
		//		zap.Error(err))
		//	return nil, nil
		//}
		//
		//entity := &strategyEntity{
		//	Data:    dbSt,
		//	Runtime: runtime,
		//}
		//
		//// 当前内存是否在
		//if err = s.userStrategy.AddStrategy(entity); err != nil {
		//	Log.Error("策略已存在", zap.String("key", string(key)))
		//	return nil, err
		//}
		//defer func() {
		//	if err != nil {
		//		s.userStrategy.RemoveStrategy(key)
		//	}
		//}()
		//
		//// todo 策略类的初始化
		//// 	包括 1. 是否加载历史数据
		//// 		2. 一些需要提前准备的数据
		//
		//if req.Setting.Contract {
		//	contract := gateway.Factor.GetContract(req.VtSymbol)
		//	if err = runtime.OnContract(contract.ContractData()); err != nil {
		//		Log.Error("添加策略失败: 初始化策略的contract失败",
		//			zap.String("key", string(key)),
		//			zap.Error(err))
		//		return nil, err
		//	}
		//}
		//
		//if err = runtime.OnInit(); err != nil {
		//	Log.Error("添加策略失败: 初始化策略on_init失败",
		//		zap.String("key", string(key)),
		//		zap.Error(err))
		//	return nil, err
		//}
		//
		////if req.Setting.LoadBar != 0 {
		////	// todo 获取历史数据，然后调用on_bar
		////	runtime.OnBar(nil)
		////}
		//
		//Log.Info("添加策略成功", zap.String("key", string(key)))

		return nil, nil
	})
	return err
}

func (e *Engine) Run() {

}

type userCtx map[object.UserId]*userContext

func (uc userCtx) GetUserContext(id object.UserId) *userContext {
	uCtx, ok := uc[id]
	if !ok {
		return nil
	}
	return uCtx
}

type userGateway map[object.Gateway]gateway.UserGateway
type userStrategy map[object.VtSymbol]strategy

type userContext struct {
	Id         object.UserId
	gateways   userGateway
	strategies userStrategy
}

func NewUserContext(id object.UserId) *userContext {
	return &userContext{
		Id:         id,
		gateways:   make(userGateway),
		strategies: make(userStrategy),
	}
}

func (uc *userContext) GetUserGateway(gateway object.Gateway) gateway.UserGateway {
	g, ok := uc.gateways[gateway]
	if !ok {
		return nil
	}
	return g
}

func (uc *userContext) AddUserGateway(userGateway gateway.UserGateway) {
	uc.gateways[userGateway.Name()] = userGateway
}

type strategy struct {
}

type PubSub interface{}
