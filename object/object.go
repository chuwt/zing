package object

import "fmt"

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

type Contract struct {
}

type (
	UserId      string
	GatewayName string
	StrategyId  int64
)

const (
	GatewayHuobi GatewayName = "huobi"
)
