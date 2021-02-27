package gateway

import (
	"errors"
	"github.com/chuwt/zing/db"
	"github.com/chuwt/zing/object"
)

var (
	GFactory = NewFactory()
)

// 公共接口
type Global interface {
	Gateway
	GetContract() ([]*db.Contract, error)
	NewUserGateway(object.UserId, object.ApiAuthJson) (UserGateway, error)
}

// 用户登陆的接口
type UserGateway interface {
	Gateway
}

type Gateway interface {
	Name() object.Gateway
	Init() error
}

// 鉴权接口
type ApiAuth interface {
	Sign(method, host, path string, params object.Params) string
	NewSignParams() object.Params
	NewWsSignParams() object.Params
}

type (
	NewGatewayFunc func(object.UserId, object.ApiAuthJson) (UserGateway, error)
	newGatewayFunc map[object.Gateway]NewGatewayFunc
	contract       map[object.VtSymbol]*db.Contract
	gateways       map[object.Gateway]Global
)

var (
	gatewayNotExitedError  = errors.New("gateway not existed")
	contractNotExitedError = errors.New("contract not existed")
)

type Factory struct {
	gateways       gateways       // Global
	newGatewayFunc newGatewayFunc // gateway创建组
	contract       contract       // 交易对维护列表

	//userCtx map[object.UserId]*UserCtx
}

func NewFactory() Factory {
	return Factory{
		newGatewayFunc: make(newGatewayFunc),
		contract:       make(contract),
	}
}

func (f *Factory) NewUserGateway(userId object.UserId, gateway object.Gateway, auth object.ApiAuthJson) (Gateway, error) {
	if gatewayFunc, ok := f.newGatewayFunc[gateway]; !ok {
		return nil, gatewayNotExitedError
	} else {
		newGateway, err := gatewayFunc(userId, auth)
		if err != nil {
			return nil, err
		}
		return newGateway, nil
	}
}

//
//func (f *Factory) GetGatewaysByUserId(userId object.UserId) map[object.Gateway]*Gateway {
//	return f.userCtx[userId].gateway
//}

//func (f *Factory) GetGatewayByGatewayName(userId object.UserId, gatewayName object.Gateway) *Gateway {
//	gateways := f.GetGatewaysByUserId(userId)
//	if gateways == nil {
//		return nil
//	}
//	if gateway, ok := gateways[gatewayName]; !ok {
//		return nil
//	} else {
//		return gateway
//	}
//}

// 添加全局gateway
func (f *Factory) AddGateway(gateway object.Gateway, newFunc func() Global) error {
	p := newFunc()
	if err := p.Init(); err != nil {
		return err
	}
	contractList, err := p.GetContract()
	if err != nil {
		return err
	}
	for _, contract := range contractList {
		f.addContract(contract.VtSymbol(), contract)
	}
	f.newGatewayFunc[gateway] = p.NewUserGateway
	f.gateways[gateway] = p
	return nil
}

func (f *Factory) GetContract(vtSymbol object.VtSymbol) (*db.Contract, error) {
	if contract, ok := f.contract[vtSymbol]; !ok {
		// todo 通过接口获取

		return nil, contractNotExitedError
	} else {
		return contract, nil
	}
}

func (f *Factory) addContract(vtSymbol object.VtSymbol, contract *db.Contract) {
	f.contract[vtSymbol] = contract
}

//func (f *Factory) AddBalance(userId object.UserId, balance *db.Balance) error {
//	userCtx, ok := f.userCtx[userId]
//	if !ok {
//		return errors.New("userCtx not existed")
//	}
//	userCtx.balance[balance.VtBalance()] = balance
//	return nil
//}

//var (
//	Factor = NewFactor()
//)

type UserCtx struct {
	userId  object.UserId
	balance map[object.VtBalance]*db.Balance // 余额管理
	gateway map[object.Gateway]*Gateway      // gateway
	order   map[object.VtSymbol][]string     // 订单
}

func NewUserCtx(userId object.UserId) *UserCtx {
	return &UserCtx{
		userId:  userId,
		balance: make(map[object.VtBalance]*db.Balance),
		gateway: make(map[object.Gateway]*Gateway),
		order:   make(map[object.VtSymbol][]string),
	}
}
