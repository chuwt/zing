package python

import (
	"fmt"
	"sync"
	"testing"
	"time"
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

	wg := sync.WaitGroup{}
	strategies := make([]*lib.PyObject, 0)
	for i := 0; i < 500; i++ {
		i := i
		wg.Add(1)
		go func() {
			strategy := pe.NewStrategyInstance(
				"MaDingStrategy",
				object.StrategyId(i),
				object.VtSymbol{
					GatewayName: object.GatewayName(fmt.Sprintf("huobi%d", i)),
					Symbol:      "btcusdt",
				},
				"")
			strategies = append(strategies, strategy)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("init ok")

	for i := 0; i < 2; i++ {
		swg := sync.WaitGroup{}
		for _, s := range strategies {
			swg.Add(1)
			s := s
			go func() {
				pe.ObjectCallFunc(s, "test", []string{"2"})
				swg.Done()
			}()
		}
		swg.Wait()
		<-time.After(time.Second)
	}
	fmt.Println("run ok")
	gil := lib.PyGILState_Ensure()
	o := strategies[0].GetAttrString("count")
	fmt.Println(lib.PyUnicode_AsUTF8(o.Repr()))
	lib.PyGILState_Release(gil)

	lib.PyEval_RestoreThread(state)
	_ = pe.Close()
}
