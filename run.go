package main

import (
	_ "github.com/chuwt/zing/config"
	_ "github.com/chuwt/zing/db"
	_ "net/http/pprof"
)

//func tickRecorder() error {
//	rd := recorder.NewRecorder(config.Config.Redis)
//	rd.AddPublisher(object.GatewayHuobi, huobird.NewPublisher)
//
//	if err := rd.Init(); err != nil {
//		return err
//	}
//	rd.Run()
//	return nil
//}
//
//func main() {
//	// 订阅
//	//go tickRecorder()
//
//	// debug
//	go func() {
//		_ = http.ListenAndServe("0.0.0.0:9090", nil)
//	}()
//
//	var err error
//
//	// gateway
//	factory := gateway.NewFactory()
//	err = factory.AddGateway(object.GatewayHuobi, huobi.NewGlobal)
//	if err != nil {
//		panic(err)
//	}
//
//	// engine
//	e := engine.NewEngine(factory, config.Config.Strategy)
//	e.Run()
//
//	//// 策略
//	//st := strategy.NewStrategy(config.Config.Redis, config.Config.Strategy)
//	//if err = st.Init(); err != nil {
//	//	panic(err)
//	//}
//	//
//	//if err = st.AddStrategy(strategy.AddStrategyReq{
//	//	ApiReq: strategy.ApiReq{
//	//		UserId:     "chuwt",
//	//		StrategyId: 1,
//	//	},
//	//	StrategyClassName: "MaDingStrategy",
//	//	VtSymbol: object.VtSymbol{
//	//		Gateway: "huobi",
//	//		Symbol:      "btcusdt",
//	//	},
//	//	Setting: "",
//	//}); err != nil {
//	//	return
//	//}
//	//if err = st.StartStrategy(strategy.StartStrategyReq{
//	//	ApiReq: strategy.ApiReq{
//	//		UserId:     "chuwt",
//	//		StrategyId: 1,
//	//	},
//	//}); err != nil {
//	//	panic(err)
//	//}
//
//	//if err = st.AddStrategy(strategy.AddStrategyReq{
//	//	ApiReq: strategy.ApiReq{
//	//		UserId:     "chuwt",
//	//		StrategyId: 2,
//	//	},
//	//	StrategyClassName: "MaDingStrategy",
//	//	VtSymbol: object.VtSymbol{
//	//		Gateway: "huobi",
//	//		Symbol:      "ethusdt",
//	//	},
//	//	Setting: "",
//	//}); err != nil {
//	//	return
//	//}
//
//	//err = ConnectGateway("chuwt", object.GatewayHuobi)
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	select {}
//}
//
////
////func ConnectGateway(userId object.UserId, name object.Gateway) error {
////	g, err := gateway.Factor.NewGateway(userId, name, &huobi.ApiAuth{
////		Key:    "",
////		Secret: "",
////	})
////	if err != nil {
////		return err
////	}
////	err = g.Init()
////	if err != nil {
////		return err
////	}
////	return nil
////}
