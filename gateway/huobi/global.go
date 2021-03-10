package huobi

import (
	"errors"
	"github.com/chuwt/zing/db"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/gateway/huobi/rest"
	"github.com/chuwt/zing/gateway/huobi/ws"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	"go.uber.org/zap"
)

var (
	Log = zap.L().With(zap.Namespace("huobi-gateway"))
)

type Global struct {
	rest rest.HuoBi
	ws   ws.HuoBi
}

func NewGlobal() gateway.Global {
	return &Global{
		rest: rest.NewRest("https://api.huobipro.com"),
		ws:   ws.NewWs("wss://api.huobi.pro/ws/v2"),
	}
}

func (p *Global) Name() object.Gateway {
	return object.GatewayHuobi
}

func (p *Global) NewUserGateway(userId object.UserId, auth object.ApiAuthJson) (gateway.UserGateway, error) {
	hb := HuoBi{
		Global: p,
		userId: userId,
	}
	if auth == "" {
		return nil, errors.New("nil apiKey")
	}
	apiAuth := new(ApiAuth)
	if err := json.Json.Unmarshal([]byte(auth), apiAuth); err != nil {
		return nil, errors.New("apiKey format error")
	}
	hb.rest.AddAuth(apiAuth)
	hb.ws.AddAuth(apiAuth)
	return &hb, nil
}

func (p *Global) Init() error {
	// 一些初始化
	return nil
}

func (p *Global) Start(chan object.Event) {}

func (p *Global) GetContract() ([]*db.Contract, error) {
	Log.Info("开始获取交易对信息")
	symbols, err := p.rest.CommonSymbols(nil)
	if err != nil {
		Log.Error(
			"接口获取交易对失败",
			zap.Error(err),
		)
		// todo 是否可以从数据库获取
		return nil, err
	}
	dbContract := make([]*db.Contract, 0, len(symbols.Data))
	// 然后设置交易对
	for _, symbol := range symbols.Data {
		// 入库
		contract := &db.Contract{
			Symbol:          symbol.Symbol,
			Gateway:         string(p.Name()),
			Product:         object.ProductNormal,
			BaseCurrency:    symbol.BaseCurrency,
			QuoteCurrency:   symbol.QuoteCurrency,
			PricePrecision:  symbol.PricePrecision,
			AmountPrecision: symbol.AmountPrecision,
			ValuePrecision:  symbol.ValuePrecision,
			MinOrderValue:   symbol.MinOrderValue.String(),
			LimitMinAmt:     symbol.LimitOrderMinOrderAmt.String(),
			LimitMaxAmt:     symbol.LimitOrderMaxOrderAmt.String(),
		}
		_, err := db.CreateDupEntry(db.GetEngine(), contract)
		if err != nil {
			Log.Error("数据库插入交易对错误",
				zap.Error(err))
			return nil, err
		}
		dbContract = append(dbContract, contract)
	}
	Log.Info("交易对信息获取成功")
	return dbContract, nil
}
