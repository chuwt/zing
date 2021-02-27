package db

import (
	"github.com/shopspring/decimal"
	"github.com/chuwt/zing/object"
)

type Contract struct {
	Base            `xorm:"extends"`
	Symbol          string `xorm:"unique(g) comment(交易对)"`
	Gateway         string `xorm:"unique(g)"`
	Product         string `xorm:"comment(现货)"`
	BaseCurrency    string `xorm:"comment(基础货币)"`
	QuoteCurrency   string `xorm:"comment(报价货币)"`
	PricePrecision  int32  `xorm:"comment(价格精度，小数点后几位)"`
	AmountPrecision int32  `xorm:"comment(数量精度，小数点后几位)"`
	ValuePrecision  int32  `xorm:"comment(成交精度，小数点后几位)"`
	MinOrderValue   string `xorm:"comment(最小下单金额)"`
	LimitMinAmt     string `xorm:"comment(现价单最小下单量)"`
	LimitMaxAmt     string `xorm:"comment(现价单最大下单量)"`
}

func (c *Contract) VtSymbol() object.VtSymbol {
	return object.VtSymbol{
		GatewayName: object.Gateway(c.Gateway),
		Symbol:      c.Symbol,
	}
}

func (c *Contract) ContractData() *object.ContractData {
	return &object.ContractData{
		VtSymbol:      c.VtSymbol(),
		Product:       c.Product,
		MinOrderAmt:   decimal.RequireFromString(c.LimitMinAmt),
		MinOrderValue: decimal.RequireFromString(c.MinOrderValue),
		// 1 / 10**amountPrecision
		MinVolume: decimal.NewFromInt(1).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt32(c.AmountPrecision))),
	}
}

func GetContractByGateway(gatewayName string) ([]*Contract, error) {
	contractList := make([]*Contract, 0)
	if err := GetEngine().Where("gateway=?", gatewayName).Find(&contractList); err != nil {
		return nil, err
	}
	return contractList, nil
}
