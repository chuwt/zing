package zing

import (
	"fmt"
	"github.com/chuwt/zing/config"
	dataCenter "github.com/chuwt/zing/data_center"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/gateway/huobi"
	"github.com/chuwt/zing/object"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

func TestGateway(t *testing.T) {
	var err error
	factory := gateway.NewFactory()
	err = factory.AddGateway(object.GatewayHuobi, huobi.NewGlobal)
	if err != nil {
		panic(err)
	}
	g, _ := factory.NewUserGateway("1", object.GatewayHuobi, "")
	fmt.Println(g.Name())
}

func TestUserCtx(t *testing.T) {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	var err error
	dataTower := dataCenter.NewRedisTower(config.Config.Redis)
	factory := gateway.NewFactory()
	err = factory.AddGateway(object.GatewayHuobi, huobi.NewGlobal)
	if err != nil {
		panic(err)
	}
	g, err := factory.NewUserGateway("chuwt", object.GatewayHuobi, object.ApiAuthJson(config.Config.DebugApiKey))
	if err != nil {
		t.Log(err)
		return
	}

	zing := NewZing(nil, dataTower)

	err = zing.AddUser("chuwt")
	if err != nil {
		t.Log(err)
		return
	}

	err = zing.AddUserGateway("chuwt", g)
	if err != nil {
		t.Log(err)
		return
	}

	err = zing.AdduserStrategy(
		"chuwt",
		1,
		object.VtSymbol{
			GatewayName: "huobi",
			Symbol:      "btcusdt",
		},
		object.StrategySetting{},
	)
	if err != nil {
		t.Log(err)
		return
	}

	err = zing.RunUserStrategy("chuwt", 1)
	if err != nil {
		t.Log(err)
		return
	}

	select {}

}
