package ws

type HuobiPosition struct {
	Currency  string `json:"currency"`
	AccountId int64  `json:"accountId"`
	Balance   string `json:"balance"`   // 余额
	Available string `json:"available"` // 可用余额
	// 余额变动类型，有效值：order-place(订单创建)，order-match(订单成交)，order-refund(订单成交退款)，order-cancel(订单撤销)，order-fee-refund(点卡抵扣交易手续费)，margin-transfer(杠杆账户划转)，margin-loan(借币本金)，margin-interest(借币计息)，margin-repay(归还借币本金币息)，deposit (充币）, withdraw (提币)，other(其他资产变化)
	ChangeType string `json:"changeType"`
	ChangeTime int64
}
