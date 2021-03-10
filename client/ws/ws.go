package ws

import (
	"context"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

var Log = zap.L().With(zap.Namespace("websocket"))

type Websocket struct {
	Ctx      context.Context
	CtxF     context.CancelFunc
	conn     *websocket.Conn
	RecvChan chan []byte
	Err      chan error
}

func NewWsClient(addr string) (Websocket, error) {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 20 * time.Second
	c, _, err := dialer.Dial(addr, nil)
	if err != nil {
		return Websocket{}, err
	}

	ws := Websocket{
		conn:     c,
		RecvChan: make(chan []byte, 1024),
		Err:      make(chan error, 1),
	}
	ws.Ctx, ws.CtxF = context.WithCancel(context.Background())
	return ws, nil
}

func (ws *Websocket) DataReceiver(unCompress func([]byte) ([]byte, error)) {
	go ws.onReceive(unCompress)
}

func (ws *Websocket) SendStruct(msg interface{}) error {
	msgBytes, _ := json.Json.Marshal(msg)
	Log.Debug("发送消息", zap.String("data", string(msgBytes)))
	err := ws.conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		Log.Error("发送消息失败",
			zap.String("data", string(msgBytes)),
		)
		return err
	}
	return nil
}

func (ws *Websocket) SendTextMsg(msg Msg) error {
	Log.Debug("发送消息", zap.String("data", string(msg.Data)))
	err := ws.conn.WriteMessage(websocket.TextMessage, msg.Data)
	if err != nil {
		Log.Error("发送消息失败",
			zap.String("gateway", string(msg.GatewayName)),
			zap.String("userId", string(msg.UserId)),
			zap.String("data", string(msg.Data)),
		)
		return err
	}
	return nil
}

func (ws *Websocket) onReceive(unCompress func([]byte) ([]byte, error)) {
	var (
		message []byte
		err     error
	)
	for {
		select {
		case <-ws.Ctx.Done():
			return
		default:
			_, message, err = ws.conn.ReadMessage()
			if err != nil {
				ws.CtxF()
				ws.Err <- err
				Log.Error("读取消息失败", zap.Error(err))
				return
			}
			if unCompress != nil {
				if message, err = unCompress(message); err != nil {
					Log.Warn("读取消息解压缩失败", zap.Error(err))
					continue
				}
			}
			ws.RecvChan <- message
			Log.Debug("读取消息", zap.String("data", string(message)))
		}
	}
}

type Msg struct {
	GatewayName object.Gateway
	UserId      object.UserId
	Data        []byte
}
