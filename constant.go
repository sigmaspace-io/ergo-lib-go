package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// Constant represents Ergo constant(evaluated) values
type Constant interface {
	// Base16 encode as Base16-encoded ErgoTree serialized value or throw an error if serialization failed
	Base16() (string, error)
	// Type returns the Constant type as string
	Type() (string, error)
	// Value returns the Constant value as string
	Value() (string, error)
}

type constant struct {
	p C.ConstantPtr
}

func newConstant(c *constant) Constant {
	runtime.SetFinalizer(c, finalizeConstant)
	return c
}

// NewConstant creates a new Constant from Base16-encoded ErgoTree serialized value
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

	errPtr := C.ergo_lib_constant_to_base16(c.p, &constantStr)
	defer C.ergo_lib_delete_string(constantStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantStr), nil
}

func (c *constant) Type() (string, error) {
	var constantTypeStr *C.char

	errPtr := C.ergo_lib_constant_type_to_dbg_str(c.p, &constantTypeStr)
	defer C.ergo_lib_delete_string(constantTypeStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantTypeStr), nil
}

func (c *constant) Value() (string, error) {
	var constantValueStr *C.char

	errPtr := C.ergo_lib_constant_value_to_dbg_str(c.p, &constantValueStr)
	defer C.ergo_lib_delete_string(constantValueStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return strings.ReplaceAll(C.GoString(constantValueStr), " ", ""), nil
}

func finalizeConstant(c *constant) {
	C.ergo_lib_constant_delete(c.p)
}
