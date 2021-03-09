package zing

import (
	"github.com/chuwt/zing/config"
	dataCenter "github.com/chuwt/zing/data_center"
	"github.com/chuwt/zing/object"
	"testing"
)

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
