package main

import (
	"net/http"
	_ "net/http/pprof"
	"vngo/config"
	_ "vngo/config"
	_ "vngo/db"
	"vngo/gateway"
	"vngo/gateway/huobi"
	"vngo/object"
	"vngo/recorder"
	huobird "vngo/recorder/huobi"
	"vngo/strategy"
)

func tickRecorder() error {
	rd := recorder.NewRecorder(config.Config.Redis)
	rd.AddPublisher(object.GatewayHuobi, huobird.NewPublisher)

	if err := rd.Init(); err != nil {
		return err
	}
	rd.Run()
	return nil
}

func main() {
	// 订阅
	//go tickRecorder()
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9090", nil)
	}()

	// gateway
	err := gateway.Factor.AddGateway(object.GatewayHuobi, huobi.NewPublic)
	if err != nil {
		panic(err)
	}

	// 策略
	st := strategy.NewStrategy(config.Config.Redis)
	err = st.AddStrategy(strategy.AddStrategyReq{
		UserId:     "chuwt",
		StrategyId: 1,
		VtSymbol: object.VtSymbol{
			GatewayName: "huobi",
			Symbol:      "btcusdt",
		},
		Setting: "",
	})
	if err != nil {
		return
	}

	//err = ConnectGateway("chuwt", object.GatewayHuobi)
	//if err != nil {
	//	panic(err)
	//}

	select {}
}

func ConnectGateway(userId object.UserId, name object.GatewayName) error {
	g, err := gateway.Factor.NewGateway(userId, name, &huobi.Api{
		Key:    "",
		Secret: "",
	})
	if err != nil {
		return err
	}
	err = g.Init()
	if err != nil {
		return err
	}
	return nil
}
