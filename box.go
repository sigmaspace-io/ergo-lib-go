package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type BoxId interface {
	Base16() string
}

type boxId struct {
	p C.BoxIdPtr
}

func newBoxId(b *boxId) BoxId {
	runtime.SetFinalizer(b, finalizeBoxId)

	return b
}

// NewBoxId creates a new ergo box id from the supplied base16 string.
func NewBoxId(s string) (BoxId, error) {
	boxIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(boxIdStr))

	var p C.BoxIdPtr

	errPtr := C.ergo_lib_box_id_from_str(boxIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &boxId{p}

	return newBoxId(b), nil
}

func (b *boxId) Base16() string {
	var boxIdStr *C.char
	defer C.ergo_lib_delete_string(boxIdStr)

	C.ergo_lib_box_id_to_str(b.p, &boxIdStr)

	return C.GoString(boxIdStr)
}

func finalizeBoxId(b *boxId) {
	C.ergo_lib_box_id_delete(b.p)
}
