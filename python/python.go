package python

import (
	"errors"
	"fmt"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/python/lib"
	"os"
	"runtime"
	"strings"
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

func (pe *PyEngine) ObjectCallFunc(obj *lib.PyObject, funcName string, args ...string) (*resp, error) {

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()

	pyArgs := lib.PyTuple_New(len(args))

	for index, arg := range args {
		s := lib.PyUnicode_FromString(arg)
		lib.PyTuple_SetItem(pyArgs, index, s)
		//s.DecRef()
	}
	//defer pyArgs.DecRef()

	pyFunc := obj.GetAttrString(funcName)
	resObj := pyFunc.Call(pyArgs, nil)
	//pyFunc.DecRef()
	res := lib.PyUnicode_AsUTF8(resObj.Repr())
	lib.PyGILState_Release(gil)
	if res == "" || res == "None" {
		return nil, nil
	}
	rp := new(resp)
	if err := json.Json.Unmarshal([]byte(strings.ReplaceAll(res, "'", "")), rp); err != nil {
		return nil, err
	}

	return rp, nil
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
	_, err := s.engine.ObjectCallFunc(s.pyObject, "on_start")
	return err
}

func (s *Strategy) OnBar(bar *object.BarData) error {
	return nil
}

func (s *Strategy) OnContract(contract *object.ContractData) error {
	contractBytes, _ := json.Json.Marshal(contract)
	_, err := s.engine.ObjectCallFunc(s.pyObject, "on_contract", string(contractBytes))
	return err
}

func (s *Strategy) OnInit() error {
	res, err := s.engine.ObjectCallFunc(s.pyObject, "on_init")
	if err != nil {
		return err
	}
	if res.Status == 0 {
		return errors.New(res.Msg)
	}
	return nil
}

type resp struct {
	Msg    string `json:"msg"`
	Status int32  `json:"status"`
}

// todo python 返回的处理
