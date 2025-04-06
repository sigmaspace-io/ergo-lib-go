package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"errors"
	"runtime"
)

// nilErrorStr is the value C.ergo_lib_error_to_string() returns
// if there is no error contained in the error pointer.
const nilErrorStr = "success"

type ergoError struct {
	p C.ErrorPtr
}

// newError creates an Error from the supplied ergo ErrorPtr.
func newError(err C.ErrorPtr) ergoError {
	e := ergoError{p: err}

	runtime.AddCleanup(&e, finalizeError, e.p)

	return e
}

func (e ergoError) isError() bool {
	return e.p != nil
}

func (e ergoError) error() error {
	cStr := C.ergo_lib_error_to_string(e.p)
	defer C.ergo_lib_delete_string(cStr)
	runtime.KeepAlive(e)
	s := C.GoString(cStr)

	if s == nilErrorStr {
		return nil
	}

	return errors.New(s)
}

func finalizeError(p C.ErrorPtr) {
	C.ergo_lib_delete_error(p)
}
