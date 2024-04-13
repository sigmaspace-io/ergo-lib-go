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
	pointer() C.ContextExtensionPtr
}

type contextExtension struct {
	p C.ContextExtensionPtr
}

func newContextExtension(c *contextExtension) ContextExtension {
	runtime.SetFinalizer(c, finalizeContextExtension)
	return c
}

func (c *contextExtension) Keys() []byte {
	bytesLength := C.ergo_lib_context_extension_len(c.p)

	output := C.malloc(C.uintptr_t(bytesLength))
	defer C.free(unsafe.Pointer(output))

	C.ergo_lib_context_extension_keys(c.p, (*C.uint8_t)(output))

	result := C.GoBytes(unsafe.Pointer(output), C.int(bytesLength))

	return result
}

func (c *contextExtension) pointer() C.ContextExtensionPtr {
	return c.p
}

func finalizeContextExtension(c *contextExtension) {
	C.ergo_lib_context_extension_delete(c.p)
}
