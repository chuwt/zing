package ws

import (
	"github.com/chuwt/zing/client/ws"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	"go.uber.org/zap"
	"strings"
)

var Log = zap.L().With(zap.Namespace("huobi user ws"))

type HuoBi struct {
	*ws.Websocket
	Name object.Gateway
	Host string
	api  gateway.ApiAuth
}

func NewWs(host string) HuoBi {
	return HuoBi{
		Host: host,
		Name: object.GatewayHuobi,
	}
}

func (h *HuoBi) AddAuth(api gateway.ApiAuth) {
	h.api = api
}

func (h *HuoBi) Init() error {
	// 连接
	w, err := ws.NewWsClient(h.Host)
	if err != nil {
		return err
	}
	h.Websocket = &w
	//var err error
	//if err = h.Connect(); err != nil {
	//	return err
	//}
	//// 登陆
	//if err = h.Login(); err != nil {
	//	return err
	//}
	return nil
}

func (h *HuoBi) Start(receiver chan object.Event) error {
	// 数据处理
	h.dataProcess(receiver)
	// 消息接收
	h.DataReceiver(nil)
	// 登陆
	return h.Login()
}

func (h *HuoBi) Connect() error {
	w, err := ws.NewWsClient(h.Host)
	if err != nil {
		return err
	}
	h.Websocket = &w
	//// 数据处理
	//h.dataProcess()
	//// 消息接收
	//h.DataReceiver(nil)
	return nil
}

func (h *HuoBi) dataProcess(receiver chan object.Event) {
	var (
		err   error
		event object.Event
	)
	go func() {
		for {
			select {
			case msg := <-h.RecvChan:
				if event, err = h.OnReceive(msg); err != nil {
					h.CtxF()
					h.Err <- err
					return
				} else {
					if event.Type != object.EventTypeNone {
						// 发送数据
						receiver <- event
					}
				}
			case <-h.Ctx.Done():
				return
			}
		}
	}()
}

func (h *HuoBi) OnReceive(data []byte) (object.Event, error) {
	res := new(ResData)
	if err := json.Json.Unmarshal(data, res); err != nil {
		return object.Event{}, err
	}
	switch res.Action {
	case "req":
		switch res.Ch {
		case "auth":
			// 登陆成功
			_ = h.UserOrder()
			_ = h.UserTrade()
			_ = h.UserAccount()
			// 成功后需要发送成功信息，此时接收端会知道重连了
			return object.Event{
				Type: object.EventTypeConnection,
				Data: h.Name,
			}, nil
		}
	case "push":
		chs := strings.Split(res.Ch, ".")
		switch chs[0] {
		case "accounts":
			gatewayPosition := new(HuobiPosition)
			_ = json.Json.Unmarshal(res.Data, gatewayPosition)

			position := object.PositionData{
				Gateway:   h.Name,
				Currency:  gatewayPosition.Currency,
				Balance:   gatewayPosition.Balance,
				Available: gatewayPosition.Available,
			}

			return object.Event{
				Type: object.EventTypePosition,
				Data: position,
			}, nil
		case "trade":
			gatewayTrade := new(HuobiTrade)
			_ = json.Json.Unmarshal(res.Data, gatewayTrade)

			trade := object.TradeData{
				Gateway:    h.Name,
				Symbol:     gatewayTrade.Symbol,
				OrderId:    gatewayTrade.OrderId,
				TradeId:    gatewayTrade.TradeId,
				Price:      gatewayTrade.OrderPrice,
				Volume:     gatewayTrade.TradeVolume,
				CreateTime: gatewayTrade.OrderCreateTime / 1000,
			}
			_, trade.Direction, trade.Offset = switchType(gatewayTrade.OrderType)
			trade.OrderStatus = switchStatus(gatewayTrade.OrderStatus)
			// todo 关于手续费的处理, ht或者其他
			trade.Fee = gatewayTrade.TransactFee

			return object.Event{
				Type: object.EventTypeTrade,
				Data: trade,
			}, nil
		case "orders":
			gatewayOrder := new(HuobiOrder)
			_ = json.Json.Unmarshal(res.Data, gatewayOrder)
			order := object.OrderData{
				Gateway:       h.Name,
				Symbol:        gatewayOrder.Symbol,
				OrderId:       gatewayOrder.OrderId,
				ClientOrderId: gatewayOrder.ClientOrderId,
				Price:         gatewayOrder.OrderPrice,
				Size:          gatewayOrder.OrderSize,
				Value:         gatewayOrder.OrderValue,
				CreateTime:    gatewayOrder.OrderCreateTime / 1000,
				RemainAmt:     gatewayOrder.RemainAmt,
				ExecAmt:       gatewayOrder.ExecAmt,
			}
			// type direction offset status 的赋值判断
			order.Type, order.Direction, order.Offset = switchType(gatewayOrder.Type)
			order.Status = switchStatus(gatewayOrder.OrderStatus)
			//switch gatewayOrder.EventType {
			//case "creation":
			//case "trade":
			//case "cancellation":
			//}

			return object.Event{
				Type: object.EventTypeOrder,
				Data: order,
			}, nil
		default:
			Log.Warn("未知类型", zap.String("data", string(data)))
		}
	case "ping":
		ping := new(Ping)
		if err := json.Json.Unmarshal(res.Data, ping); err != nil {
			return object.Event{}, err
		}
		if err := h.SendMsg(Pong{
			Action: "pong",
			Data:   ping,
		}); err != nil {
			return object.Event{}, err
		}
	}
	return object.Event{}, nil
}

func (h *HuoBi) UserAccount() error {
	return h.SendMsg(Req{
		Action: "sub",
		Ch:     "accounts.update#1",
	})
}

func (h *HuoBi) UserTrade() error {
	return h.SendMsg(Req{
		Action: "sub",
		Ch:     "trade.clearing#*#0",
	})
}

func (h *HuoBi) UserOrder() error {
	return h.SendMsg(Req{
		Action: "sub",
		Ch:     "orders#*",
	})
}

func (h *HuoBi) Login() error {
	p := h.api.NewWsSignParams()
	auth := Auth{
		Req: Req{
			Action: "req",
			Ch:     "auth",
		},
		Params: AuthParams{
			AuthType:         "api",
			AccessKeyId:      p["accessKey"],
			SignatureMethod:  p["signatureMethod"],
			SignatureVersion: p["signatureVersion"],
			Timestamp:        p["timestamp"],
		},
	}
	auth.Params.Signature = h.api.Sign("GET", h.Host, "/ws/v2", p)
	return h.SendMsg(auth)
}

func (h *HuoBi) SendMsg(msg interface{}) error {
	msgBytes, _ := json.Json.Marshal(msg)
	return h.SendTextMsg(ws.Msg{
		Data: msgBytes,
	})
}
