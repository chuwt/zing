package receiver

import "github.com/chuwt/zing/object"

type Tick struct {
	id       object.StrategyId
	queue    chan object.TickData
	isClosed bool
	cap      int
}

func NewTickReceiver(id object.StrategyId, cap int64) *Tick {
	return &Tick{
		id:    id,
		queue: make(chan object.TickData, cap),
		cap:   int(cap),
	}
}

func (tr *Tick) Id() object.StrategyId {
	return tr.id
}

func (tr *Tick) Get() chan object.TickData {
	return tr.queue
}

func (tr *Tick) OnReceive() chan object.TickData {
	return tr.queue
}

func (tr *Tick) IsClosed() bool {
	return tr.isClosed
}

// 需要手动关闭
func (tr *Tick) Close() {
	tr.isClosed = true
}

func (tr *Tick) IsFull() bool {
	if tr.cap <= len(tr.queue) {
		return true
	}
	return false
}
