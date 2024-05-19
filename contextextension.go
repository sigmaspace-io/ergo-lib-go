package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ContextExtension represent user-defined variables to be put into context
type ContextExtension interface {
	// Keys returns all keys in the map
	Keys() []byte
	// Get returns Constant at provided key or nil if it doesn't exist
	Get(key uint8) (Constant, error)
	// Set adds Constant at provided key
	Set(key uint8, constant Constant)
	pointer() C.ContextExtensionPtr
}

type contextExtension struct {
	p C.ContextExtensionPtr
}

func newContextExtension(c *contextExtension) ContextExtension {
	runtime.SetFinalizer(c, finalizeContextExtension)
	return c
}

// NewContextExtension creates new empty ContextExtension instance
func NewContextExtension() ContextExtension {
	var p C.ContextExtensionPtr
	C.ergo_lib_context_extension_empty(&p)
	c := &contextExtension{p: p}
	return newContextExtension(c)
}

func (c *contextExtension) Keys() []byte {
	bytesLength := C.ergo_lib_context_extension_len(c.p)

	output := C.malloc(C.uintptr_t(bytesLength))
	defer C.free(unsafe.Pointer(output))

	C.ergo_lib_context_extension_keys(c.p, (*C.uint8_t)(output))

	result := C.GoBytes(unsafe.Pointer(output), C.int(bytesLength))

	return result
}

func (c *contextExtension) Get(key uint8) (Constant, error) {
	var p C.ConstantPtr

	res := C.ergo_lib_context_extension_get(c.p, C.uint8_t(key), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		co := &constant{p: p}
		return newConstant(co), nil
	}

	return nil, nil
}

func (c *contextExtension) Set(key uint8, constant Constant) {
	C.ergo_lib_context_extension_set_pair(constant.pointer(), C.uint8_t(key), c.p)
}

func (c *contextExtension) pointer() C.ContextExtensionPtr {
	return c.p
}

func finalizeContextExtension(c *contextExtension) {
	C.ergo_lib_context_extension_delete(c.p)
}
