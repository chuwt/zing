package lib

//#cgo pkg-config: python-3.7
//#include "Python.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func PyUnicode_AsUTF8(unicode *PyObject) string {
	cutf8 := C.PyUnicode_AsUTF8(toc(unicode))
	return C.GoString(cutf8)
}

func PyDict_New() *PyObject {
	return togo(C.PyDict_New())
}

func PyTuple_New(len int) *PyObject {
	return togo(C.PyTuple_New(C.Py_ssize_t(len)))
}

func PyTuple_SetItem(p *PyObject, pos int, o *PyObject) int {
	return int(C.PyTuple_SetItem(toc(p), C.Py_ssize_t(pos), toc(o)))
}

func Initialize() error {
	if C.Py_IsInitialized() == 0 {
		C.Py_Initialize()
	}
	if C.Py_IsInitialized() == 0 {
		return fmt.Errorf("python: could not initialize the python interpreter")
	}

	if C.PyEval_ThreadsInitialized() == 0 {
		C.PyEval_InitThreads()
	}
	if C.PyEval_ThreadsInitialized() == 0 {
		return fmt.Errorf("python: could not initialize the GIL")
	}

	return nil
}

func Finalize() error {
	C.Py_Finalize()
	return nil
}

//PyEval_InitThreads : https://docs.python.org/3/c-api/init.html#c.PyEval_InitThreads
func PyEval_InitThreads() {
	C.PyEval_InitThreads()
}

func PyImport_ImportModule(name string) *PyObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return togo(C.PyImport_ImportModule(cname))
}

//togo converts a *C.PyObject to a *PyObject
func togo(cobject *C.PyObject) *PyObject {
	return (*PyObject)(cobject)
}

func toc(object *PyObject) *C.PyObject {
	return (*C.PyObject)(object)
}

//MaxVariadicLength is the maximum number of arguments that can be passed to a variadic C function due to a cgo limitation
const MaxVariadicLength = 20

// Constants used for comparison in PyObject_RichCompareBool
var (
	Py_LT = int(C.Py_LT)
	Py_LE = int(C.Py_LE)
	Py_EQ = int(C.Py_EQ)
	Py_NE = int(C.Py_NE)
	Py_GT = int(C.Py_GT)
	Py_GE = int(C.Py_GE)
)

//None : https://docs.python.org/3/c-api/none.html#c.Py_None
var Py_None = togo(C.Py_None)

//PyObject : https://docs.python.org/3/c-api/structures.html?highlight=pyobject#c.PyObject
type PyObject C.PyObject

//IncRef : https://docs.python.org/3/c-api/refcounting.html#c.Py_INCREF
func (pyObject *PyObject) IncRef() {
	C.Py_IncRef(toc(pyObject))
}

//DecRef : https://docs.python.org/3/c-api/refcounting.html#c.Py_DECREF
func (pyObject *PyObject) DecRef() {
	C.Py_DecRef(toc(pyObject))
}

//ReprEnter : https://docs.python.org/3/c-api/exceptions.html#c.Py_ReprEnter
func (pyObject *PyObject) ReprEnter() int {
	return int(C.Py_ReprEnter(toc(pyObject)))
}

//ReprLeave : https://docs.python.org/3/c-api/exceptions.html#c.Py_ReprLeave
func (pyObject *PyObject) ReprLeave() {
	C.Py_ReprLeave(toc(pyObject))
}

//HasAttr : https://docs.python.org/3/c-api/object.html#c.PyObject_HasAttr
func (pyObject *PyObject) HasAttr(attr_name *PyObject) bool {
	return C.PyObject_HasAttr(toc(pyObject), toc(attr_name)) == 1
}

//HasAttrString : https://docs.python.org/3/c-api/object.html#c.PyObject_HasAttrString
func (pyObject *PyObject) HasAttrString(attr_name string) bool {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return C.PyObject_HasAttrString(toc(pyObject), cattr_name) == 1
}

//GetAttr : https://docs.python.org/3/c-api/object.html#c.PyObject_GetAttr
func (pyObject *PyObject) GetAttr(attr_name *PyObject) *PyObject {
	return togo(C.PyObject_GetAttr(toc(pyObject), toc(attr_name)))
}

//GetAttrString : https://docs.python.org/3/c-api/object.html#c.PyObject_GetAttrString
func (pyObject *PyObject) GetAttrString(attr_name string) *PyObject {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return togo(C.PyObject_GetAttrString(toc(pyObject), cattr_name))
}

//SetAttr : https://docs.python.org/3/c-api/object.html#c.PyObject_SetAttr
func (pyObject *PyObject) SetAttr(attr_name *PyObject, v *PyObject) int {
	return int(C.PyObject_SetAttr(toc(pyObject), toc(attr_name), toc(v)))
}

//SetAttrString : https://docs.python.org/3/c-api/object.html#c.PyObject_SetAttrString
func (pyObject *PyObject) SetAttrString(attr_name string, v *PyObject) int {
	cattr_name := C.CString(attr_name)
	defer C.free(unsafe.Pointer(cattr_name))

	return int(C.PyObject_SetAttrString(toc(pyObject), cattr_name, toc(v)))
}

////DelAttr : https://docs.python.org/3/c-api/object.html#c.PyObject_DelAttr
//func (pyObject *PyObject) DelAttr(attr_name *PyObject) int {
//	return int(C._go_PyObject_DelAttr(toc(pyObject), toc(attr_name)))
//}
//
////DelAttrString : https://docs.python.org/3/c-api/object.html#c.PyObject_DelAttrString
//func (pyObject *PyObject) DelAttrString(attr_name string) int {
//	cattr_name := C.CString(attr_name)
//	defer C.free(unsafe.Pointer(cattr_name))
//
//	return int(C._go_PyObject_DelAttrString(toc(pyObject), cattr_name))
//}

//RichCompare : https://docs.python.org/3/c-api/object.html#c.PyObject_RichCompare
func (pyObject *PyObject) RichCompare(o *PyObject, opid int) *PyObject {
	return togo(C.PyObject_RichCompare(toc(pyObject), toc(o), C.int(opid)))
}

//RichCompareBool : https://docs.python.org/3/c-api/object.html#c.PyObject_RichCompareBool
func (pyObject *PyObject) RichCompareBool(o *PyObject, opid int) int {
	return int(C.PyObject_RichCompareBool(toc(pyObject), toc(o), C.int(opid)))
}

//Repr : https://docs.python.org/3/c-api/object.html#c.PyObject_Repr
func (pyObject *PyObject) Repr() *PyObject {
	return togo(C.PyObject_Repr(toc(pyObject)))
}

//ASCII : https://docs.python.org/3/c-api/object.html#c.PyObject_ASCII
func (pyObject *PyObject) ASCII() *PyObject {
	return togo(C.PyObject_ASCII(toc(pyObject)))
}

//Str : https://docs.python.org/3/c-api/object.html#c.PyObject_Str
func (pyObject *PyObject) Str() *PyObject {
	return togo(C.PyObject_Str(toc(pyObject)))
}

//Bytes : https://docs.python.org/3/c-api/object.html#c.PyObject_Bytes
func (pyObject *PyObject) Bytes() *PyObject {
	return togo(C.PyObject_Bytes(toc(pyObject)))
}

//IsSubclass : https://docs.python.org/3/c-api/object.html#c.PyObject_IsSubclass
func (pyObject *PyObject) IsSubclass(cls *PyObject) int {
	return int(C.PyObject_IsSubclass(toc(pyObject), toc(cls)))
}

//IsInstance : https://docs.python.org/3/c-api/object.html#c.PyObject_IsInstance
func (pyObject *PyObject) IsInstance(cls *PyObject) int {
	return int(C.PyObject_IsInstance(toc(pyObject), toc(cls)))
}

// PyCallable_Check : https://docs.python.org/3/c-api/object.html#c.PyCallable_Check
func PyCallable_Check(o *PyObject) bool {
	return C.PyCallable_Check(toc(o)) == 1
}

//Call : https://docs.python.org/3/c-api/object.html#c.PyObject_Call
func (pyObject *PyObject) Call(args *PyObject, kwargs *PyObject) *PyObject {
	return togo(C.PyObject_Call(toc(pyObject), toc(args), toc(kwargs)))
}

//CallObject : https://docs.python.org/3/c-api/object.html#c.PyObject_CallObject
func (pyObject *PyObject) CallObject(args *PyObject) *PyObject {
	return togo(C.PyObject_CallObject(toc(pyObject), toc(args)))
}

func PyGILState_Check() bool {
	return C.PyGILState_Check() == 1
}

////CallFunctionObjArgs : https://docs.python.org/3/c-api/object.html#c.PyObject_CallFunctionObjArgs
//func (pyObject *PyObject) CallFunctionObjArgs(args ...*PyObject) *PyObject {
//
//	if len(args) > MaxVariadicLength {
//		panic("CallFunctionObjArgs: too many arrguments")
//	}
//	if len(args) == 0 {
//		return togo(C._go_PyObject_CallFunctionObjArgs(toc(pyObject), 0, (**C.PyObject)(nil)))
//	}
//
//	cargs := make([]*C.PyObject, len(args), len(args))
//	for i, arg := range args {
//		cargs[i] = toc(arg)
//	}
//	return togo(C._go_PyObject_CallFunctionObjArgs(toc(pyObject), C.int(len(args)), (**C.PyObject)(unsafe.Pointer(&cargs[0]))))
//}

////CallMethodObjArgs : https://docs.python.org/3/c-api/object.html#c.PyObject_CallMethodObjArgs
//func (pyObject *PyObject) CallMethodObjArgs(name *PyObject, args ...*PyObject) *PyObject {
//	if len(args) > MaxVariadicLength {
//		panic("CallMethodObjArgs: too many arguments")
//	}
//	if len(args) == 0 {
//		return togo(C._go_PyObject_CallMethodObjArgs(toc(pyObject), toc(name), 0, (**C.PyObject)(nil)))
//	}
//
//	cargs := make([]*C.PyObject, len(args), len(args))
//	for i, arg := range args {
//		cargs[i] = toc(arg)
//	}
//	return togo(C._go_PyObject_CallMethodObjArgs(toc(pyObject), toc(name), C.int(len(args)), (**C.PyObject)(unsafe.Pointer(&cargs[0]))))
//}

////CallMethodArgs : same as CallMethodObjArgs but with name as go string
//func (pyObject *PyObject) CallMethodArgs(name string, args ...*PyObject) *PyObject {
//	pyName := PyUnicode_FromString(name)
//	defer pyName.DecRef()
//
//	return pyObject.CallMethodObjArgs(pyName, args...)
//}

func PyUnicode_FromString(u string) *PyObject {
	cu := C.CString(u)
	defer C.free(unsafe.Pointer(cu))

	return togo(C.PyUnicode_FromString(cu))
}

//PyGILState is an opaque “handle” to the thread state when PyGILState_Ensure() was called, and must be passed to PyGILState_Release() to ensure Python is left in the same state
type PyGILState C.PyGILState_STATE

//PyGILState_Ensure : https://docs.python.org/3/c-api/init.html#c.PyGILState_Ensure
func PyGILState_Ensure() PyGILState {
	return PyGILState(C.PyGILState_Ensure())
}

//PyGILState_Release : https://docs.python.org/3/c-api/init.html#c.PyGILState_Release
func PyGILState_Release(state PyGILState) {
	C.PyGILState_Release(C.PyGILState_STATE(state))
}

//PyThreadState : https://docs.python.org/3/c-api/init.html#c.PyThreadState
type PyThreadState C.PyThreadState

//PyEval_SaveThread : https://docs.python.org/3/c-api/init.html#c.PyEval_SaveThread
func PyEval_SaveThread() *PyThreadState {
	return (*PyThreadState)(C.PyEval_SaveThread())
}

//PyEval_RestoreThread : https://docs.python.org/3/c-api/init.html#c.PyEval_RestoreThread
func PyEval_RestoreThread(tstate *PyThreadState) {
	C.PyEval_RestoreThread((*C.PyThreadState)(tstate))
}

//Hash : https://docs.python.org/3/c-api/object.html#c.PyObject_Hash
func (pyObject *PyObject) Hash() int {
	return int(C.PyObject_Hash(toc(pyObject)))
}

//HashNotImplemented : https://docs.python.org/3/c-api/object.html#c.PyObject_HashNotImplemented
func (pyObject *PyObject) HashNotImplemented() int {
	return int(C.PyObject_HashNotImplemented(toc(pyObject)))
}

//IsTrue : https://docs.python.org/3/c-api/object.html#c.PyObject_IsTrue
func (pyObject *PyObject) IsTrue() int {
	return int(C.PyObject_IsTrue(toc(pyObject)))
}

//Not : https://docs.python.org/3/c-api/object.html#c.PyObject_Not
func (pyObject *PyObject) Not() int {
	return int(C.PyObject_Not(toc(pyObject)))
}

//Type : https://docs.python.org/3/c-api/object.html#c.PyObject_Type
func (pyObject *PyObject) Type() *PyObject {
	return togo(C.PyObject_Type(toc(pyObject)))
}

//Length : https://docs.python.org/3/c-api/object.html#c.PyObject_Length
func (pyObject *PyObject) Length() int {
	return int(C.PyObject_Length(toc(pyObject)))
}

//LengthHint : https://docs.python.org/3/c-api/object.html#c.PyObject_LengthHint
func (pyObject *PyObject) LengthHint(pyDefault int) int {
	return int(C.PyObject_LengthHint(toc(pyObject), C.Py_ssize_t(pyDefault)))
}

//GetItem : https://docs.python.org/3/c-api/object.html#c.PyObject_GetItem
func (pyObject *PyObject) GetItem(key *PyObject) *PyObject {
	return togo(C.PyObject_GetItem(toc(pyObject), toc(key)))
}

//SetItem : https://docs.python.org/3/c-api/object.html#c.PyObject_SetItem
func (pyObject *PyObject) SetItem(key, v *PyObject) int {
	return int(C.PyObject_SetItem(toc(pyObject), toc(key), toc(v)))
}

//DelItem : https://docs.python.org/3/c-api/object.html#c.PyObject_DelItem
func (pyObject *PyObject) DelItem(key *PyObject) int {
	return int(C.PyObject_DelItem(toc(pyObject), toc(key)))
}

//Dir : https://docs.python.org/3/c-api/object.html#c.PyObject_Dir
func (pyObject *PyObject) Dir() *PyObject {
	return togo(C.PyObject_Dir(toc(pyObject)))
}

//GetIter : https://docs.python.org/3/c-api/object.html#c.PyObject_GetIter
func (pyObject *PyObject) GetIter() *PyObject {
	return togo(C.PyObject_GetIter(toc(pyObject)))
}
