package zing

import (
	"fmt"
	"github.com/chuwt/zing/config"
	dataCenter "github.com/chuwt/zing/data_center"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/gateway/huobi"
	"github.com/chuwt/zing/object"
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
	g.Init()
	fmt.Println(g.Name())
}

func TestUserCtx(t *testing.T) {
	var err error
	dataTower := dataCenter.NewRedisTower(config.Config.Redis)

	zing := NewZing(nil, dataTower)
	err = zing.AddUser("chuwt")
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
