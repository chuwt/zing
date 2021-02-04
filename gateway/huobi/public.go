package huobi

import (
	"errors"
	"go.uber.org/zap"
	"vngo/db"
	"vngo/gateway"
	"vngo/gateway/huobi/rest"
	"vngo/gateway/huobi/ws"
	"vngo/object"
)

var (
	Log = zap.L().With(zap.Namespace("huobi-gateway"))
)

type Public struct {
	rest rest.HuoBi
	ws   ws.HuoBi
}

func NewPublic() gateway.Public {
	return &Public{
		rest: rest.NewRest("https://api.huobipro.com"),
		ws:   ws.NewWs("wss://api.huobi.pro/ws/v2"),
	}
}

func (p *Public) Name() object.GatewayName {
	return object.GatewayHuobi
}

func (p *Public) NewUserGateway(userId object.UserId, api gateway.ApiAuth) (gateway.Gateway, error) {
	hb := HuoBi{
		Public: *p,
		userId: userId,
	}
	if api == nil {
		return nil, errors.New("nil apiKey")
	}
	hb.rest.AddAuth(api)
	hb.ws.AddAuth(api)
	return &hb, nil
}

func (p *Public) Init() error {
	// 连接rest
	// 获取交易对信息
	var err error
	if err = p.GetContract(); err != nil {
		return err
	}

	return nil
}

func (p *Public) GetContract() error {
	Log.Info("开始获取交易对信息")
	symbols, err := p.rest.CommonSymbols(nil)
	if err != nil {
		Log.Error(
			"接口获取交易对失败",
			zap.Error(err),
		)
		// todo 是否可以从数据库获取
		return err
	}
	// 然后设置交易对
	for _, symbol := range symbols.Data {
		// 入库
		contract := &db.Contract{
			Symbol:         symbol.Symbol,
			Gateway:        string(p.Name()),
			Product:        object.ProductNormal,
			BaseCurrency:   symbol.BaseCurrency,
			QuoteCurrency:  symbol.QuoteCurrency,
			PricePrecision: symbol.PricePrecision,
			MinOrderValue:  symbol.MinOrderValue.String(),
			LimitMinAmt:    symbol.LimitOrderMinOrderAmt.String(),
			LimitMaxAmt:    symbol.LimitOrderMaxOrderAmt.String(),
		}
		_, err := db.CreateDupEntry(db.GetEngine(), contract)
		if err != nil {
			Log.Error("数据库插入交易对错误",
				zap.Error(err))
			return err
		}
		gateway.Factor.AddContract(
			object.VtSymbol{
				GatewayName: string(p.Name()),
				Symbol:      symbol.Symbol,
			},
			contract)
	}
	Log.Info("交易对信息获取成功")
	return nil
}
