package object

import (
	"fmt"
	"strings"
)

const (
	ProductNormal = "现货"
)

// 交易对.交易所
type VtSymbol struct {
	GatewayName Gateway `json:"gateway"`
	Symbol      string  `json:"symbol"`
}

func (vt *VtSymbol) String() string {
	return fmt.Sprintf("%s.%s", vt.Symbol, vt.GatewayName)
}

func LoadVtSymbol(vt string) *VtSymbol {
	vtList := strings.Split(vt, ".")
	if len(vtList) == 2 {
		return &VtSymbol{
			GatewayName: Gateway(vtList[1]),
			Symbol:      vtList[0],
		}
	}
	return nil
}

type DataPublisher interface {
	Subscribe(symbol string) error
	Run()
}
type PubFunc func(string, []byte) error
type NewPublisher func(PubFunc) DataPublisher

type VtBalance struct {
	GatewayName string
	Currency    string
}

type (
	UserId      string
	Gateway     string
	Currency    string
	StrategyId  int64
	StrategyKey string
	ApiAuthJson string

	Params map[string]string

	ClientOrderId string
	OrderId       int64
	TradeId       int64

	VtCurrency string
)

func (sk StrategyKey) String() string {
	return string(sk)
}

type StrategySetting struct {
	RuntimeSetting string // 策略初始化的参数
	LoadBar        int    // 加载的历史数据天数
	Contract       bool   // 是否需要初始化contract
}

const (
	GatewayHuobi Gateway = "huobi"
)
