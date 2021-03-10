package object

type Event struct {
	Type int32
	Data interface{}
}

const (
	EventTypeNone = iota
	EventTypeConnection
	EventTypePosition
	EventTypeOrder
	EventTypeTrade
)
