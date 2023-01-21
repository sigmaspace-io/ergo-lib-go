package ergo

/*
   #include "ergo_lib.h"
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

type Constant interface {
	Base16() (string, error)
	ConstantType() (string, error)
	ConstantValue() (string, error)
}

type constant struct {
	p C.ConstantPtr
}

func newConstant(c *constant) Constant {
	runtime.SetFinalizer(c, finalizeConstant)

	return c
}

func NewConstant(s string) (Constant, error) {
	base16ErgoTree := C.CString(s)
	defer C.free(unsafe.Pointer(base16ErgoTree))

	var p C.ConstantPtr

	errPtr := C.ergo_lib_constant_from_base16(base16ErgoTree, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	c := &constant{p}

	return newConstant(c), nil
}

func (c *constant) Base16() (string, error) {
	var constantStr *C.char
	defer C.free(unsafe.Pointer(constantStr))

	errPtr := C.ergo_lib_constant_to_base16(c.p, &constantStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantStr), nil
}

func (c *constant) ConstantType() (string, error) {
	var constantTypeStr *C.char
	defer C.free(unsafe.Pointer(constantTypeStr))

	errPtr := C.ergo_lib_constant_type_to_dbg_str(c.p, &constantTypeStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantTypeStr), nil
}

func (c *constant) ConstantValue() (string, error) {
	var constantValueStr *C.char
	defer C.free(unsafe.Pointer(constantValueStr))

	errPtr := C.ergo_lib_constant_value_to_dbg_str(c.p, &constantValueStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return strings.ReplaceAll(C.GoString(constantValueStr), " ", ""), nil
}

func finalizeConstant(c *constant) {
	C.ergo_lib_constant_delete(c.p)
}
