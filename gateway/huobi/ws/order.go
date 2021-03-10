package ws

import "github.com/chuwt/zing/object"

type HuobiOrder struct {
	EventType string // 事件类型，有效值：creation
	Symbol    string // 交易代码
	//accountId	long	// 账户ID
	OrderId       int64  // 订单ID
	ClientOrderId string // 用户自编订单号（如有）
	//orderSource	string	// 订单来源
	OrderPrice      string // 订单价格
	OrderSize       string // 订单数量（对市价买单无效）
	OrderValue      string // 订单金额（仅对市价买单有效）
	Type            string // 订单类型，有效值：buy-market, sell-market, buy-limit, sell-limit, buy-limit-maker, sell-limit-maker, buy-ioc, sell-ioc, buy-limit-fok, sell-limit-fok
	OrderStatus     string // 订单状态，有效值：submitted
	OrderCreateTime int64  // 订单创建时间
	RemainAmt       string // 该订单未成交数量（市价买单为未成交金额）
	ExecAmt         string // 该订单累计成交量
}

func switchType(orderType string) (object.OrderType, object.Direction, object.Offset) {
	switch orderType {
	case "buy-market":
		return object.OrderTypeMARKET, object.DirectionLONG, object.OffsetOPEN
	case "sell-market":
		return object.OrderTypeMARKET, object.DirectionSHORT, object.OffsetCLOSE
	case "buy-limit":
		return object.OrderTypeLIMIT, object.DirectionLONG, object.OffsetOPEN
	case "sell-limit":
		return object.OrderTypeLIMIT, object.DirectionSHORT, object.OffsetCLOSE
	}
	return "", "", ""
}

func switchStatus(orderStatus string) object.OrderStatus {
	switch orderStatus {
	case "submitted":
		return object.OrderStatusNotTraded
	case "partial-filled":
		return object.OrderStatusPartTraded
	case "filled":
		return object.OrderStatusAllTraded
	case "cancelling":
		return object.OrderStatusCancelling
	case "partial-canceled":
		return object.OrderStatusPartTradedCancelled
	case "canceled":
		return object.OrderStatusCancelled
	}
	return ""
}
