package db

type Contract struct {
	Base           `xorm:"extends"`
	Symbol         string `xorm:"unique(g) comment(交易对)"`
	Gateway        string `xorm:"unique(g)"`
	Product        string `xorm:"comment(现货)"`
	BaseCurrency   string `xorm:"comment(基础货币)"`
	QuoteCurrency  string `xorm:"comment(报价货币)"`
	PricePrecision int32  `xorm:"comment(价格精度，小数点后几位)"`
	MinOrderValue  string `xorm:"comment(最小下单金额)"`
	LimitMinAmt    string `xorm:"comment(现价单最小下单量)"`
	LimitMaxAmt    string `xorm:"comment(现价单最大下单量)"`
}

func GetContractByGateway(gatewayName string) ([]*Contract, error) {
	contractList := make([]*Contract, 0)
	if err := GetEngine().Where("gateway=?", gatewayName).Find(&contractList); err != nil {
		return nil, err
	}
	return contractList, nil
}
