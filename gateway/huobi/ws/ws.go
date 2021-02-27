package ws

import (
	"encoding/json"
	"github.com/chuwt/zing/client/ws"
	"github.com/chuwt/zing/gateway"
	"github.com/chuwt/zing/object"
)

type HuoBi struct {
	Name   object.Gateway
	UserId object.UserId
	Host   string
	Conn   *ws.WS
	Err    chan error
	api    gateway.ApiAuth
}

func NewWs(host string) HuoBi {
	return HuoBi{
		Host: host,
		Name: object.GatewayHuobi,
		Err:  make(chan error),
	}
}

func (h *HuoBi) AddAuth(api gateway.ApiAuth) {
	h.api = api
}

func (h *HuoBi) Init(userId object.UserId) error {
	// 连接
	h.UserId = userId
	var err error
	if err = h.Connect(); err != nil {
		return err
	}
	// 登陆
	if err = h.Login(); err != nil {
		return err
	}
	return nil
}

func (h *HuoBi) Connect() error {
	var err error
	h.Conn, err = ws.NewWsClient(h.Host)
	if err != nil {
		return err
	}
	// 启动消息接收
	h.Conn.Start(nil)
	go func() {
		for {
			select {
			case msg := <-h.Conn.RecvChan:
				if err = h.OnReceive(msg); err != nil {
					h.Err <- err
				}
			case err := <-h.Conn.Err:
				h.Err <- err
			}
		}
	}()
	return nil
}

func (h *HuoBi) OnReceive(data []byte) error {
	res := new(ResData)
	if err := json.Unmarshal(data, res); err != nil {
		return err
	}
	switch res.Action {
	case "req":
		switch res.Ch {
		case "auth":
			// 登陆成功
			_ = h.UserOrder()
			_ = h.UserTrade()
			_ = h.UserAccount()
		}
	case "push":
		switch res.Ch {
		case "accounts.update#1":

		}
	case "ping":
		ping := new(Ping)
		if err := json.Unmarshal(res.Data, ping); err != nil {
			return err
		}
		pong := Pong{
			Action: "pong",
			Data:   ping,
		}
		pongBytes, _ := json.Marshal(pong)
		if err := h.Conn.SendMsg(ws.Msg{
			Data: pongBytes,
		}); err != nil {
			return err
		}
	}
	return nil
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
	msgBytes, _ := json.Marshal(msg)
	return h.Conn.SendMsg(ws.Msg{
		Data: msgBytes,
	})
}
