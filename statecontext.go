package ergo

/*
   #include "ergo.h"
*/
import "C"
import "runtime"

// StateContext represents blockchain state (last headers, etc.)
type StateContext interface {
}

type stateContext struct {
	p C.ErgoStateContextPtr
}

func newStateContext(s *stateContext) StateContext {
	runtime.SetFinalizer(s, finalizeStateContext)
	return s
}

// NewStateContext creates StateContext from PreHeader and BlockHeaders
func NewStateContext(preHeader PreHeader, headers BlockHeaders) (StateContext, error) {
	var p C.ErgoStateContextPtr

	errPtr := C.ergo_lib_ergo_state_context_new(preHeader.pointer(), headers.pointer(), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	st := &stateContext{p: p}

	return newStateContext(st), nil
}

func finalizeStateContext(s *stateContext) {
	C.ergo_lib_ergo_state_context_delete(s.p)
}
