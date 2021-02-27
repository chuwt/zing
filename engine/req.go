package engine

import (
	"fmt"
	"github.com/chuwt/zing/object"
)

type ApiReq struct {
	UserId     object.UserId
	StrategyId object.StrategyId
}

// userId.strategyId
func (ar *ApiReq) Key() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%d", ar.UserId, ar.StrategyId))
}

type AddStrategyReq struct {
	ApiReq
	StrategyClassName string // 策略类名
	VtSymbol          object.VtSymbol
	Setting           Setting
}

type Setting struct {
	RuntimeSetting string // 策略初始化的参数
	LoadBar        int    // 加载的历史数据天数
	Contract       bool   // 是否需要初始化contract
}

type AddUserGatewayReq struct {
	UserId  object.UserId
	Gateway object.Gateway
	Auth    object.ApiAuthJson
}

func (r *AddUserGatewayReq) Key() object.StrategyKey {
	return object.StrategyKey(fmt.Sprintf("%s.%s", r.UserId, r.Gateway))
}
