package gateway

import (
	"errors"
	"github.com/chuwt/fasthttp-client"
	"vngo/db"
	"vngo/object"
)

type Public interface {
	Gateway
	NewUserGateway(object.UserId, ApiAuth) (Gateway, error)
}

type Gateway interface {
	Name() object.GatewayName
	Init() error
}

type ApiAuth interface {
	Sign(method, host, path string, params fasthttp.Mapper) string
	NewSignParams() fasthttp.Mapper
	NewWsSignParams() fasthttp.Mapper
}

type NewGatewayFunc func(object.UserId, ApiAuth) (Gateway, error)

var (
	gatewayNotExited = errors.New("gateway not existed")
	gatewayExisted   = errors.New("gateway existed")
)

type Factory struct {
	gatewayFactor map[object.GatewayName]NewGatewayFunc // gateway创建组

	//userGateway map[object.UserId]map[object.GatewayName]*Gateway // 用户gateway的维护
	contract map[object.VtSymbol]*db.Contract // 交易对维护列表

	userCtx map[object.UserId]*UserCtx
}

type UserCtx struct {
	userId      object.UserId
	userBalance map[object.VtBalance]*db.Balance // 余额管理
	userGateway map[object.GatewayName]*Gateway  // gateway
	userOrder   map[object.VtSymbol][]string     // 订单
}

func NewUserCtx(userId object.UserId) *UserCtx {
	return &UserCtx{
		userId:      userId,
		userBalance: make(map[object.VtBalance]*db.Balance),
		userGateway: make(map[object.GatewayName]*Gateway),
		userOrder:   make(map[object.VtSymbol][]string),
	}
}

func NewFactor() Factory {
	return Factory{
		gatewayFactor: make(map[object.GatewayName]NewGatewayFunc),
		//userGateway:   make(map[object.UserId]map[object.GatewayName]*Gateway),
		userCtx:  make(map[object.UserId]*UserCtx),
		contract: make(map[object.VtSymbol]*db.Contract),
	}
}

func (f *Factory) NewGateway(userId object.UserId, gatewayName object.GatewayName, auth ApiAuth) (Gateway, error) {
	if gatewayFunc, ok := f.gatewayFactor[gatewayName]; !ok {
		return nil, gatewayNotExited
	} else {
		gateway, err := gatewayFunc(userId, auth)
		if err != nil {
			return nil, err
		}
		if _, ok := f.userCtx[userId]; !ok {
			//f.userGateway[userId] = make(map[object.GatewayName]*Gateway)
			f.userCtx[userId] = NewUserCtx(userId)
		} else {
			//if _, ok := f.userGateway[userId][gatewayName]; ok {
			if _, ok := f.userCtx[userId].userGateway[gatewayName]; ok {
				return nil, gatewayExisted
			}
		}
		f.userCtx[userId].userGateway[gatewayName] = &gateway
		return gateway, nil
	}
}

func (f *Factory) GetGatewaysByUserId(userId object.UserId) map[object.GatewayName]*Gateway {
	return f.userCtx[userId].userGateway
}

func (f *Factory) GetGatewayByGatewayName(userId object.UserId, gatewayName object.GatewayName) *Gateway {
	gateways := f.GetGatewaysByUserId(userId)
	if gateways == nil {
		return nil
	}
	if gateway, ok := gateways[gatewayName]; !ok {
		return nil
	} else {
		return gateway
	}
}

func (f *Factory) AddGateway(gatewayName object.GatewayName, newFunc func() Public) error {
	p := newFunc()
	if err := p.Init(); err != nil {
		return err
	}
	f.gatewayFactor[gatewayName] = p.NewUserGateway
	return nil
}

func (f *Factory) GetContract(vtSymbol object.VtSymbol) *db.Contract {
	if contract, ok := f.contract[vtSymbol]; !ok {
		// todo 通过接口获取

		return nil
	} else {
		return contract
	}
}

func (f *Factory) AddContract(vtSymbol object.VtSymbol, contract *db.Contract) {
	f.contract[vtSymbol] = contract
}

func (f *Factory) AddBalance(userId object.UserId, balance *db.Balance) error {
	userCtx, ok := f.userCtx[userId]
	if !ok {
		return errors.New("userCtx not existed")
	}
	userCtx.userBalance[balance.VtBalance()] = balance
	return nil
}

var (
	Factor = NewFactor()
)
