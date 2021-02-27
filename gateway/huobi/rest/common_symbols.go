package rest

import (
	"github.com/shopspring/decimal"
	"github.com/chuwt/zing/client/http"
)

// 获取所有交易对
func (r *HuoBi) CommonSymbols(callback http.Callback) (*CommonSymbolsRes, error) {
	path := r.GetUrl("/v1/common/symbols")
	if callback != nil {
		// 异步
		return nil, nil
	} else {
		// 同步
		commonSymbolsRes := new(CommonSymbolsRes)
		err := http.SyncGetRequest(path, commonSymbolsRes, nil)
		if err != nil {
			return nil, err
		}
		return commonSymbolsRes, nil
	}
}

type CommonSymbolsRes struct {
	BaseRes
	Data []struct {
		BaseCurrency           string          `json:"base-currency"`              // 交易对中的基础币种
		QuoteCurrency          string          `json:"quote-currency"`             // 交易对中的报价币种
		PricePrecision         int32           `json:"price-precision"`            // 交易对报价的精度（小数点后位数）
		AmountPrecision        int32           `json:"amount-precision"`           // 交易对基础币种计数精度（小数点后位数）
		SymbolPartition        string          `json:"symbol-partition"`           // 交易区，可能值: [main，innovation]
		Symbol                 string          `json:"symbol"`                     // 交易对
		State                  string          `json:"state"`                      // 交易对状态；可能值: [online，offline,suspend] online - 已上线；offline - 交易对已下线，不可交易；suspend -- 交易暂停；pre-online - 即将上线
		ValuePrecision         int32           `json:"value-precision"`            // 交易对交易金额的精度（小数点后位数）
		MinOrderValue          decimal.Decimal `json:"min-order-value"`            // 交易对限价单和市价买单最小下单金额 ，以计价币种为单位
		LimitOrderMinOrderAmt  decimal.Decimal `json:"limit-order-min-order-amt"`  // 交易对限价单最小下单量 ，以基础币种为单位（NEW）
		LimitOrderMaxOrderAmt  decimal.Decimal `json:"limit-order-max-order-amt"`  // 交易对限价单最大下单量 ，以基础币种为单位（NEW）
		SellMarketMinOrderAmt  decimal.Decimal `json:"sell-market-min-order-amt"`  // 交易对市价卖单最小下单量，以基础币种为单位（NEW）
		SellMarketMaxOrderAmt  decimal.Decimal `json:"sell-market-max-order-amt"`  // 交易对市价卖单最大下单量，以基础币种为单位（NEW）
		BuyMarketMaxOrderValue decimal.Decimal `json:"buy-market-max-order-value"` // 交易对市价买单最大下单金额，以计价币种为单位（NEW）
		ApiTrading             string          `json:"api-trading"`                // API交易使能标记（有效值：enabled, disabled）
	} `json:"data"`
}
