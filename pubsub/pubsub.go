package pubsub

type PubSub interface {
	Init() error
	Pub() error
	Sub() error
}
