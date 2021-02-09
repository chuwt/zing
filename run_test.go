package main

import (
	"net/http"
	"testing"
	"vngo/config"
	"vngo/gateway"
	"vngo/gateway/huobi"
	"vngo/object"
	"vngo/strategy"
)

func TestRun(t *testing.T) {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9090", nil)
	}()

	var err error

	// gateway
	err = gateway.Factor.AddGateway(object.GatewayHuobi, huobi.NewPublic)
	if err != nil {
		panic(err)
	}

	st := strategy.NewStrategy(config.Config.Redis, config.Config.Strategy)
	if err = st.Init(); err != nil {
		panic(err)
	}

	for i := 0; i < 1; i++ {
		i := i
		go func() {
			if err = st.AddStrategy(strategy.AddStrategyReq{
				ApiReq: strategy.ApiReq{
					UserId:     "chuwt",
					StrategyId: object.StrategyId(i),
				},
				StrategyClassName: "TestStrategy",
				VtSymbol: object.VtSymbol{
					GatewayName: "huobi",
					Symbol:      "btcusdt",
				},
				Setting: strategy.Setting{
					RuntimeSetting: "{\"size\":2,\"grid_number\":8,\"buy_width\":3,\"buy_callback\":0.5,\"sell_width\":1.3,\"sell_callback\":0.3,\"grid_loop\":false,\"grid_clear\":false}",
					LoadBar:        0,
					Contract:       true,
				},
			}); err != nil {
				return
			}
			//if err = st.StartStrategy(strategy.StartStrategyReq{
			//	ApiReq: strategy.ApiReq{
			//		UserId:     "chuwt",
			//		StrategyId: object.StrategyId(i),
			//	},
			//}); err != nil {
			//	panic(err)
			//}
		}()
	}
	select {}
}
