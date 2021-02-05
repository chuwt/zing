package object

import (
	"fmt"
	"strings"
)

const (
	ProductNormal = "现货"
)

type VtSymbol struct {
	GatewayName GatewayName `json:"gateway"`
	Symbol      string      `json:"symbol"`
}

func (vt *VtSymbol) String() string {
	return fmt.Sprintf("%s.%s", vt.Symbol, vt.GatewayName)
}

func LoadVtSymbol(vt string) *VtSymbol {
	vtList := strings.Split(vt, ".")
	if len(vtList) == 2 {
		return &VtSymbol{
			GatewayName: GatewayName(vtList[1]),
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
	GatewayName string
	StrategyId  int64
	StrategyKey string
)

const (
	GatewayHuobi GatewayName = "huobi"
)
