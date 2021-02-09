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
