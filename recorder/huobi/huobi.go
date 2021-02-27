package huobi

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"github.com/chuwt/zing/client/ws"
	"github.com/chuwt/zing/object"
)

var Log = zap.L().With(zap.Namespace("publisher-huobi"))

type Publisher struct {
	*ws.WS
	host       string
	gateway    object.Gateway
	subscribed map[string]struct{}         // 已订阅的交易对
	seq        atomic.Int64                // 订阅id
	publisher  object.PubFunc              // 消息发布方法
	tickMap    map[string]*object.TickData // 订阅的交易对tick数据
}

func NewPublisher(publisher object.PubFunc) object.DataPublisher {
	if publisher == nil {
		publisher = func(string, []byte) error { return nil }
	}
	return &Publisher{
		host:       "wss://api.huobi.pro/ws",
		gateway:    object.GatewayHuobi,
		subscribed: make(map[string]struct{}),
		seq:        *atomic.NewInt64(0),
		publisher:  publisher,
		tickMap:    make(map[string]*object.TickData),
	}
}

func (p *Publisher) Run() {
	var err error
	// todo retry 的间隔
retry:
	if err = p.init(); err != nil {
		Log.Warn("ws连接失败, 准备重连", zap.String("gateway", string(p.gateway)), zap.Error(err))
		goto retry
	}
	go p.start()
	select {
	case err = <-p.Err:
		Log.Warn("ws连接断开, 准备重连", zap.String("gateway", string(p.gateway)), zap.Error(err))
		goto retry
	}
}

func (p *Publisher) init() error {
	var err error
	p.WS, err = ws.NewWsClient(p.host)
	if err != nil {
		return err
	}
	Log.Info("ws连接成功", zap.String("gateway", string(p.gateway)))
	return nil
}

func (p *Publisher) start() {
	var (
		err error
		r   *gzip.Reader
	)
	// 启动
	p.WS.Start(func(rawMsg []byte) ([]byte, error) {
		b := new(bytes.Buffer)
		_ = binary.Write(b, binary.LittleEndian, rawMsg)
		r, err = gzip.NewReader(b)
		if err != nil {
			return nil, err
		}
		msgByte, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return msgByte, nil
	})

	// 重新订阅
	if err = p.reSubscribe(); err != nil {
		Log.Info("ws重新订阅失败", zap.String("gateway", string(p.gateway)), zap.Error(err))
		p.Err <- err
		return
	}

	// 监听数据
	for {
		select {
		case msg := <-p.RecvChan:
			if err = p.onReceiveData(msg); err != nil {
				p.Err <- err
				return
			}
		}
	}
}

func (p *Publisher) onReceiveData(msg []byte) error {
	var (
		err error
		res = new(Resp)
	)

	if err = json.Unmarshal(msg, res); err != nil {
		return err
	}

	if res.Ping != 0 {
		err = p.SendStruct(Pong{
			Pong: res.Ping,
		})
		if err != nil {
			return err
		}
	} else {
		chList := strings.Split(res.Ch, ".")
		if len(chList) >= 3 {
			symbol := chList[1]
			tick, ok := p.tickMap[symbol]
			if !ok {
				return nil
			}
			tick.Timestamp = res.Ts / 1000

			switch chList[2] {
			case "depth":
				// 深度

			case "detail":
				t := new(object.Tick)
				_ = json.Unmarshal(res.Tick, t)
				// 市场高开低收
				tick.Tick.High = t.High
				tick.Tick.Open = t.Open
				tick.Tick.Low = t.Low
				tick.Tick.Close = t.Close
				tick.Tick.Vol = t.Vol
			case "trade":
				// 最新成交信息
				t := new(Trade)
				_ = json.Unmarshal(res.Tick, t)
				tick.LastPrice = t.Data[0].Price
				tick.LastVolume = t.Data[0].Amount
			}
			// todo 这里做了一个限制，1s只推一次
			if tick.LastPrice != decimal.Zero && tick.LastTime != tick.Timestamp {
				tickBytes, _ := json.Marshal(tick)
				if err = p.publisher(tick.VtSymbol.String(), tickBytes); err != nil {
					//if err = p.publisher("huobi", tickBytes); err != nil {
					return err
				}
				tick.LastTime = tick.Timestamp
			}
		}
	}

	return nil
}

func (p *Publisher) reSubscribe() error {
	var err error
	for symbol := range p.subscribed {
		if err = p.subscribe(symbol); err != nil {
			return err
		}
	}
	return nil
}

func (p *Publisher) Subscribe(symbol string) error {
	p.subscribed[symbol] = struct{}{}
	return p.subscribe(symbol)
}

func (p *Publisher) subscribe(symbol string) error {
	var err error
	var subData SubData

	subData = SubData{
		Id:  p.incrSeq(),
		Sub: fmt.Sprintf("market.%s.depth.step0", symbol),
	}
	err = p.SendStruct(subData)
	if err != nil {
		return err
	}

	subData = SubData{
		Id:  p.incrSeq(),
		Sub: fmt.Sprintf("market.%s.detail", symbol),
	}
	err = p.SendStruct(subData)
	if err != nil {
		return err
	}

	subData = SubData{
		Id:  p.incrSeq(),
		Sub: fmt.Sprintf("market.%s.trade.detail", symbol),
	}
	err = p.SendStruct(subData)
	if err != nil {
		return err
	}
	p.tickMap[symbol] = &object.TickData{
		VtSymbol: object.VtSymbol{
			Symbol:      symbol,
			GatewayName: object.GatewayHuobi,
		},
	}
	return nil
}

func (p *Publisher) incrSeq() int64 {
	return p.seq.Inc()
}
