package ergo

/*
   #include "ergo.h"
*/
import "C"
import "runtime"

// StateContext represents blockchain state (last headers, etc.)
type StateContext interface {
	pointer() C.ErgoStateContextPtr
}

type stateContext struct {
	p C.ErgoStateContextPtr
}

func newStateContext(s *stateContext) StateContext {
	runtime.SetFinalizer(s, finalizeStateContext)
	return s
}

// NewStateContext creates StateContext from PreHeader and BlockHeaders
func NewStateContext(preHeader PreHeader, headers BlockHeaders, parameters Parameters) (StateContext, error) {
	var p C.ErgoStateContextPtr

	errPtr := C.ergo_lib_ergo_state_context_new(preHeader.pointer(), headers.pointer(), parameters.pointer(), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	st := &stateContext{p: p}

	return newStateContext(st), nil
}

func (s *stateContext) pointer() C.ErgoStateContextPtr {
	return s.p
}

func finalizeStateContext(s *stateContext) {
	C.ergo_lib_ergo_state_context_delete(s.p)
}
