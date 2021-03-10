package object

type OrderType string

const (
	OrderTypeLIMIT  OrderType = "限价"
	OrderTypeMARKET OrderType = "市价"
	OrderTypeSTOP   OrderType = "STOP"
	OrderTypeFAK    OrderType = "FAK"
	OrderTypeFOK    OrderType = "FOK"
	OrderTypeRFQ    OrderType = "询价"
)

type Direction string

const (
	DirectionLONG  = "多"
	DirectionSHORT = "空"
)

type Offset string

const (
	OffsetOPEN  = "开"
	OffsetCLOSE = "平"
)

type OrderStatus string

const (
	OrderStatusSubmitting          = "提交中"
	OrderStatusNotTraded           = "未成交"
	OrderStatusPartTraded          = "部分成交"
	OrderStatusAllTraded           = "全部成交"
	OrderStatusCancelling          = "撤销中"
	OrderStatusCancelled           = "已撤销"
	OrderStatusRejected            = "拒单"
	OrderStatusPartTradedCancelled = "部分成交已撤单"
	OrderStatusSuccess             = "报单成功"
	OrderStatusFail                = "报单失败"
)
