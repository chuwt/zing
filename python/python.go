package python

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"vngo/object"
	"vngo/python/lib"
)

// todo 暂时没搞懂垃圾回收

type PyEngine struct {
	path            string
	pythonPath      string
	init            bool
	state           *lib.PyThreadState
	strategyFactory *lib.PyObject
}

func NewPyEngine(strategyPath, pythonPath string) PyEngine {
	return PyEngine{
		path:            strategyPath,
		pythonPath:      pythonPath,
		init:            false,
		strategyFactory: nil,
	}
}

func (pe *PyEngine) Init() error {
	var err error
	if err = os.Setenv("PYTHONPATH", pe.pythonPath); err != nil {
		return err
	}
	if err = lib.Initialize(); err != nil {
		return err
	}
	pe.init = true
	return pe.Prepare()
}

func (pe *PyEngine) Close() error {
	pe.init = false
	lib.PyEval_RestoreThread(pe.state)
	return lib.Finalize()
}

func (pe *PyEngine) Prepare() error {
	module := lib.PyImport_ImportModule("main")
	pe.strategyFactory = module.GetAttrString("get_strategy_instance")
	if pe.strategyFactory == nil {
		return errors.New("create strategy factory error")
	}
	pe.state = lib.PyEval_SaveThread()
	return nil
}

func (pe *PyEngine) NewStrategyInstance2(strategyClassName string, strategyId object.StrategyId, symbol object.VtSymbol, setting string) (*Strategy, error) {

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()

	pyArgs := lib.PyTuple_New(5)
	s0 := lib.PyUnicode_FromString(pe.path)
	s1 := lib.PyUnicode_FromString(strategyClassName)
	s2 := lib.PyUnicode_FromString(fmt.Sprintf("%d", strategyId))
	s3 := lib.PyUnicode_FromString(symbol.String())
	s4 := lib.PyUnicode_FromString(setting)

	lib.PyTuple_SetItem(pyArgs, 0, s0)
	lib.PyTuple_SetItem(pyArgs, 1, s1)
	lib.PyTuple_SetItem(pyArgs, 2, s2)
	lib.PyTuple_SetItem(pyArgs, 3, s3)
	lib.PyTuple_SetItem(pyArgs, 4, s4)

	res := pe.strategyFactory.Call(pyArgs, nil)
	//s4.DecRef()
	//s3.DecRef()
	//s2.DecRef()
	//s1.DecRef()
	//s0.DecRef()
	//pyArgs.DecRef()
	lib.PyGILState_Release(gil)
	// todo res 是空的情况

	return &Strategy{
		pyObject: res,
		engine:   pe,
	}, nil
}

func (pe *PyEngine) NewStrategyInstance(strategyClassName string, strategyId object.StrategyId, symbol object.VtSymbol, setting string) *lib.PyObject {

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()

	pyArgs := lib.PyTuple_New(5)
	s0 := lib.PyUnicode_FromString(pe.path)
	s1 := lib.PyUnicode_FromString(strategyClassName)
	s2 := lib.PyUnicode_FromString(fmt.Sprintf("%d", strategyId))
	s3 := lib.PyUnicode_FromString(symbol.String())
	s4 := lib.PyUnicode_FromString(setting)

	lib.PyTuple_SetItem(pyArgs, 0, s0)
	lib.PyTuple_SetItem(pyArgs, 1, s1)
	lib.PyTuple_SetItem(pyArgs, 2, s2)
	lib.PyTuple_SetItem(pyArgs, 3, s3)
	lib.PyTuple_SetItem(pyArgs, 4, s4)

	res := pe.strategyFactory.Call(pyArgs, nil)
	//s4.DecRef()
	//s3.DecRef()
	//s2.DecRef()
	//s1.DecRef()
	//s0.DecRef()
	//pyArgs.DecRef()
	lib.PyGILState_Release(gil)

	return res
}

func (pe *PyEngine) ObjectCallFunc(obj *lib.PyObject, funcName string, args []string) *lib.PyObject {

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()

	pyArgs := lib.PyTuple_New(len(args))

	for index, arg := range args {
		s := lib.PyUnicode_FromString(arg)
		lib.PyTuple_SetItem(pyArgs, index, s)
		s.DecRef()
	}
	//defer pyArgs.DecRef()

	pyFunc := obj.GetAttrString(funcName)
	res := pyFunc.Call(pyArgs, nil)
	//pyFunc.DecRef()
	lib.PyGILState_Release(gil)

	return res
}

type Strategy struct {
	pyObject *lib.PyObject
	engine   *PyEngine
}

func (s *Strategy) Init() {

}

func (s *Strategy) Start() {

}

func (s *Strategy) OnTick(tick *object.TickData) error {
	count := s.engine.ObjectCallFunc(s.pyObject, "on_start", nil)
	gil := lib.PyGILState_Ensure()
	fmt.Println(lib.PyUnicode_AsUTF8(count.Repr()))
	lib.PyGILState_Release(gil)
	return nil
}
