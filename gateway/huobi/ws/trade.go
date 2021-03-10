package ws

type HuobiTrade struct {
	Symbol          string // 交易代码
	OrderId         int64  // 订单id
	TradePrice      string // 成交价
	TradeVolume     string // 成交量
	OrderSide       string // 订单方向，有效值： buy, sell
	OrderType       string // 订单类型，有效值： buy-market, sell-market,buy-limit,sell-limit,buy-ioc,sell-ioc,buy-limit-maker,sell-limit-maker,buy-stop-limit,sell-stop-limit,buy-limit-fok, sell-limit-fok, buy-stop-limit-fok, sell-stop-limit-fok
	aggressor       bool   // 是否交易主动方，有效值： true, false
	TradeId         int64  // 交易ID
	TradeTime       int64  // 成交时间，unix time in millisecond
	TransactFee     string // 交易手续费（正值）或交易手续费返佣（负值）
	FeeCurrency     string // 交易手续费或交易手续费返佣币种（买单的交易手续费币种为基础币种，卖单的交易手续费币种为计价币种；买单的交易手续费返佣币种为计价币种，卖单的交易手续费返佣币种为基础币种）
	FeeDeduct       string // 交易手续费抵扣
	FeeDeductType   string // 交易手续费抵扣类型，有效值： ht, point
	AccountId       int64  // 账户编号
	Source          string // 订单来源
	OrderPrice      string // 订单价格 （市价单无此字段）
	OrderSize       string // 订单数量（市价买单无此字段）
	OrderValue      string // 订单金额（仅市价买单有此字段）
	ClientOrderId   string // 用户自编订单号
	StopPrice       string // 订单触发价（仅止盈止损订单有此字段）
	Operator        string // 订单触发方向（仅止盈止损订单有此字段）
	OrderCreateTime int64  // 订单创建时间
	OrderStatus     string // 订单状态，有效值：filled, partial-filled, canceled, partial-canceled
}
