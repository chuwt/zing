package zing

import (
	"github.com/chuwt/zing/config"
	dataCenter "github.com/chuwt/zing/data_center"
	"github.com/chuwt/zing/object"
	"testing"
)

func TestUserCtx(t *testing.T) {
	userCtx := NewUserCtx(dataCenter.NewRedisTower(config.Config.Redis))

	userCtx.AddUser("chuwt")

	userCtx.AddUserStrategy("chuwt", 1, object.VtSymbol{
		GatewayName: "huobi",
		Symbol:      "btcusdt",
	}, object.StrategySetting{})

	userCtx.StartUserStrategy("chuwt", 1)

}
