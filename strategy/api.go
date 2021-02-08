package strategy

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"vngo/db"
	"vngo/object"
)

type ApiReq struct {
	UserId     object.UserId
	StrategyId object.StrategyId
}

type AddStrategyReq struct {
	ApiReq
	StrategyClassName string
	VtSymbol          object.VtSymbol
	Setting           string
}

func (ar *ApiReq) Key() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%d", ar.UserId, ar.StrategyId))
}

/*
创建策略（初始化资源）——————>终止策略（释放资源）
   |                        ^
   |						|
   |——————>启动策略—————————>|
               |		    |
               |———————> 暂停策略
*/

/* 创建策略
会初始化策略，并将策略加载到内存
*/
func (s *Strategy) AddStrategy(req AddStrategyReq) error {

	var err error
	key := req.Key()

	Log.Info("开始添加策略", zap.String("key", string(key)))

	// 同一用户的同一个策略做过滤
	_, err, _ = s.once.Do(string(key), func() (interface{}, error) {
		dbSt := &db.Strategy{
			UserId:     string(req.UserId),
			StrategyId: int64(req.StrategyId),
			Symbol:     req.VtSymbol.Symbol,
			Gateway:    string(req.VtSymbol.GatewayName),
			Setting:    req.Setting,
			ClassName:  req.StrategyClassName,
			Status:     0,
		}

		// 入库
		duplicate, err := db.CreateDupEntry(db.GetEngine(), dbSt)
		if err != nil {
			Log.Error("策略入库失败", zap.String("key", string(key)), zap.Error(err))
			return nil, err
		} else if duplicate {
			Log.Warn("策略已存在", zap.String("key", string(key)))
		}

		// 创建策略运行时
		runtime, err := s.pyEngine.NewStrategyInstance2(
			req.StrategyClassName,
			req.StrategyId,
			req.VtSymbol,
			req.Setting)
		if err != nil {
			Log.Error("添加策略失败: 创建策略运行时失败",
				zap.String("key", string(key)),
				zap.Error(err))
			return nil, nil
		}

		entity := &strategyEntity{
			Data:    dbSt,
			Runtime: runtime,
		}

		// 当前内存是否在
		if err = s.userStrategy.AddStrategy(entity); err != nil {
			Log.Error("策略已存在", zap.String("key", string(key)))
			return nil, err
		}

		Log.Info("添加策略成功", zap.String("key", string(key)))

		// todo 策略类的初始化
		// 	包括 1. 是否加载历史数据
		// 		2. 一些需要提前准备的数据

		return nil, nil
	})

	return err
}

type StartStrategyReq struct {
	ApiReq
}

/* 启动策略
会启动内存里的策略，如果内存不存在，则返回错误
*/
func (s *Strategy) StartStrategy(strategy StartStrategyReq) error {
	var err error
	key := strategy.Key()
	Log.Info("启动策略", zap.String("key", string(key)))
	// 同一用户的同一个策略做过滤
	_, err, _ = s.once.Do(string(key), func() (interface{}, error) {
		// 当前内存是否在
		entity := s.userStrategy.GetStrategyByKey(key)
		if entity == nil {
			Log.Error("策略不存在", zap.String("key", string(key)))
			return nil, errors.New("strategy not existed")
		}

		// 订阅交易对
		if err := s.Sub(object.VtSymbol{
			GatewayName: object.GatewayName(entity.Data.Gateway),
			Symbol:      entity.Data.Symbol,
		}, entity); err != nil {
			Log.Error("策略启动失败: 订阅失败",
				zap.String("key", string(key)),
				zap.Any("strategy", entity),
				zap.Error(err))
			return nil, nil
		}

		entity.Data.Status = db.StrategyStatusRunning
		// 修改状态为启动中
		if err := db.Update(db.GetEngine(), entity.Data); err != nil {
			// todo 如果出现错误，是不是要取消订阅
			Log.Error("启动策略失败: 更新数据库失败",
				zap.String("key", string(key)),
				zap.Any("strategy", entity),
				zap.Error(err))

			return nil, nil
		}

		Log.Info("启动策略成功", zap.String("key", string(key)))
		return nil, nil
	})

	return err
}

/* 修改策略
只有暂停的策略才能修改，会修改数据库和内存
*/
func (s *Strategy) EditStrategy() {

}

/* 暂停策略
会暂停策略
*/
func (s *Strategy) StopStrategy() {

}

/* 终止策略
会终止策略，删除内存里的策略
*/
func (s *Strategy) RemoveStrategy(userId object.UserId, strategyId object.StrategyId) {

}
