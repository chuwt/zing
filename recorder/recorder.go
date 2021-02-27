package recorder

import (
	"context"
	"errors"
	pubsub "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"github.com/chuwt/zing/client/redis"
	"github.com/chuwt/zing/object"
)

var Log = zap.L().With(zap.Namespace("recorder"))

type Recorder struct {
	ctx   context.Context
	redis *redis.Redis

	publisherFactory map[object.Gateway]object.NewPublisher  // 新对象工厂
	publisherMap     map[object.Gateway]object.DataPublisher // 创建后的保存
}

func NewRecorder(cfg redis.Config) Recorder {
	return Recorder{
		ctx:              context.Background(),
		redis:            redis.NewRedis(cfg),
		publisherFactory: make(map[object.Gateway]object.NewPublisher),
		publisherMap:     make(map[object.Gateway]object.DataPublisher),
	}
}

// 添加工厂函数
func (r *Recorder) AddPublisher(gateway object.Gateway, factory object.NewPublisher) {
	r.publisherFactory[gateway] = factory
	Log.Info(
		"添加gateway初始化方法成功",
		zap.String("gateway", string(gateway)),
	)
}

// 初始化所有订阅
func (r *Recorder) Init() error {
	for gateway := range r.publisherFactory {
		if err := r.newPublisher(gateway); err != nil {
			return err
		}
		Log.Info(
			"实例化gateway成功",
			zap.String("gateway", string(gateway)),
		)
	}
	return nil
}

// 创建订阅实例
func (r *Recorder) newPublisher(gateway object.Gateway) error {
	newFunc, ok := r.publisherFactory[gateway]
	if !ok {
		Log.Error("gateway不存在", zap.String("gateway", string(gateway)))
		return errors.New("gateway not existed")
	}
	_, ok = r.publisherMap[gateway]
	if ok {
		Log.Error("gateway已存在", zap.String("gateway", string(gateway)))
		return errors.New("gateway existed")
	}

	publisher := newFunc(r.pub)
	go publisher.Run()

	r.publisherMap[gateway] = publisher
	return nil
}

// 发布
func (r *Recorder) pub(channel string, msg []byte) error {
	return r.redis.Publish(r.ctx, channel, msg).Err()
}

// 订阅
func (r *Recorder) Run() {
	var (
		err      error
		vtSymbol *object.VtSymbol
		msg      *pubsub.Message
	)
retry:
	pubSub := r.redis.Subscribe(r.ctx, "subscribe_symbol")
	for {
		msg, err = pubSub.ReceiveMessage(r.ctx)
		if err != nil {
			Log.Error("接收订阅消息失败", zap.Error(err))
			_ = pubSub.Close()
			goto retry
		}
		Log.Debug("接收订阅消息", zap.String("data", msg.Payload))

		vtSymbol = object.LoadVtSymbol(msg.Payload)

		// 根据不同的vtSymbol进行订阅
		publisher, ok := r.publisherMap[vtSymbol.GatewayName]
		if !ok {
			Log.Error("gateway不存在", zap.String("gateway", string(vtSymbol.GatewayName)))
			continue
		}
		if err = publisher.Subscribe(vtSymbol.Symbol); err != nil {
			Log.Error(
				"订阅失败",
				zap.String("gateway", string(vtSymbol.GatewayName)),
				zap.String("symbol", vtSymbol.Symbol),
			)
			continue
		}
	}
}
