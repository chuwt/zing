package object

type Event struct {
	Type    int32
	Gateway Gateway
	Data    interface{}
}

const (
	EventTypeNone = iota
	EventTypeConnection
	EventTypePosition
	EventTypeOrder
	EventTypeTrade
)
