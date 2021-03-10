package object

import (
	"github.com/shopspring/decimal"
)

type TickData struct {
	VtSymbol
	Timestamp int64 `json:"timestamp"`

	LastPrice  decimal.Decimal `json:"last_price"`
	LastVolume decimal.Decimal `json:"last_volume"`

	Tick

	BidPrice1 decimal.Decimal
	BidPrice2 decimal.Decimal
	BidPrice3 decimal.Decimal
	BidPrice4 decimal.Decimal
	BidPrice5 decimal.Decimal

	AskPrice1 decimal.Decimal
	AskPrice2 decimal.Decimal
	AskPrice3 decimal.Decimal
	AskPrice4 decimal.Decimal
	AskPrice5 decimal.Decimal

	LastTime int64 `json:"-"`
}

type Tick struct {
	High  decimal.Decimal `json:"high"`
	Open  decimal.Decimal `json:"open"`
	Low   decimal.Decimal `json:"low"`
	Close decimal.Decimal `json:"close"`
	Vol   decimal.Decimal `json:"vol"`
}

type BarData struct {
}

type ContractData struct {
	VtSymbol
	Name    string `json:"name"`
	Product string `json:"product"`

	MinOrderAmt   decimal.Decimal `json:"min_order_amt"`
	MinOrderValue decimal.Decimal `json:"min_order_value"`
	MinVolume     decimal.Decimal `json:"min_volume"`
}

type PositionData struct {
	Gateway   Gateway `json:"gateway"`
	Currency  string  `json:"currency"`
	Balance   string  `json:"balance"`   // 余额
	Available string  `json:"available"` // 可用余额
}

func (pd *PositionData) VtCurrency() VtCurrency {
	return VtCurrency(string(pd.Gateway) + "." + pd.Currency)
}

type OrderData struct {
	Gateway       Gateway
	Symbol        string      // 交易代码
	OrderId       int64       // 订单ID
	ClientOrderId string      // 用户自编订单号（如有）
	Type          OrderType   // 订单类型
	Direction     Direction   // 方向
	Offset        Offset      // 开平
	Price         string      // 下单价格
	Size          string      // 下单数量
	Value         string      // 下单金额
	Status        OrderStatus // 订单状态，有效值：submitted
	CreateTime    int64       // 订单创建时间
	RemainAmt     string      // 订单未成交量（金额）
	ExecAmt       string      // 订单累计成交量（金额）
}

type TradeData struct {
	Gateway     Gateway
	Symbol      string
	OrderId     int64
	TradeId     int64
	Direction   Direction
	Offset      Offset
	Price       string // 价格
	Volume      string // 数量（未减去手续费）
	Fee         string // 手续费
	CreateTime  int64
	OrderStatus OrderStatus
}
