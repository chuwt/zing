package python

import (
	"fmt"
	"sync"
	"testing"
	"vngo/object"
	"vngo/python/lib"
)

func TestPython(t *testing.T) {
	pe := NewPyEngine("/Volumes/hdd1000gb/workspace/src/vngo/python/vngo/strategies")
	if err := pe.Init(); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer pe.Close()

	state := lib.PyEval_SaveThread()

	strategy := pe.NewStrategyInstance(
		"MaDingStrategy",
		object.StrategyId(1),
		object.VtSymbol{
			GatewayName: "huobi",
			Symbol:      "btcusdt",
		},
		"")

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			pe.ObjectCallFunc(strategy, "on_start", nil)
			wg.Done()
		}()
	}
	wg.Wait()
	lib.PyEval_RestoreThread(state)

}
