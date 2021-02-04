package huobi

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type SubData struct {
	Id  int64  `json:"id"`
	Sub string `json:"sub"`
}

type Resp struct {
	Ping int64           `json:"ping"`
	Ch   string          `json:"ch"`
	Ts   int64           `json:"ts"`
	Tick json.RawMessage `json:"tick"`
}

type Pong struct {
	Pong int64 `json:"pong"`
}

type Trade struct {
	Data []struct {
		Amount decimal.Decimal `json:"amount"`
		Price  decimal.Decimal `json:"price"`
	} `json:"data"`
}
