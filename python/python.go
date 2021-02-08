package python

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"vngo/object"
	"vngo/python/lib"
)

type PyEngine struct {
	path            string
	init            bool
	strategyFactory *lib.PyObject
}

func NewPyEngine(strategyPath string) PyEngine {
	return PyEngine{
		path:            strategyPath,
		init:            false,
		strategyFactory: nil,
	}
}

func (pe *PyEngine) Init() error {
	var err error
	if err = os.Setenv("PYTHONPATH", "./vngo"); err != nil {
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
	return lib.Finalize()
}

func (pe *PyEngine) Prepare() error {
	module := lib.PyImport_ImportModule("main")
	pe.strategyFactory = module.GetAttrString("get_strategy_instance")
	if pe.strategyFactory == nil {
		return errors.New("create strategy factory error")
	}
	return nil
}

func (pe *PyEngine) NewStrategyInstance(strategyClassName string, strategyId object.StrategyId, symbol object.VtSymbol, setting string) *lib.PyObject {

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

	defer func() {
		s4.DecRef()
		s3.DecRef()
		s2.DecRef()
		s1.DecRef()
		s0.DecRef()
	}()

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()
	res := pe.strategyFactory.Call(pyArgs, nil)
	pyArgs.DecRef()
	lib.PyGILState_Release(gil)

	return res
}

func (pe *PyEngine) ObjectCallFunc(obj *lib.PyObject, funcName string, args []string) *lib.PyObject {

	pyArgs := lib.PyTuple_New(len(args))

	for index, arg := range args {
		s := lib.PyUnicode_FromString(arg)
		lib.PyTuple_SetItem(pyArgs, index, s)
		s.DecRef()
	}
	defer pyArgs.DecRef()

	runtime.LockOSThread()
	gil := lib.PyGILState_Ensure()
	pyFunc := obj.GetAttrString(funcName)
	res := pyFunc.Call(pyArgs, nil)
	pyFunc.DecRef()
	lib.PyGILState_Release(gil)

	return res
}

type Strategy struct {
	*lib.PyObject
}

func (s *Strategy) Init() {

}

func (s *Strategy) Start() {

}
