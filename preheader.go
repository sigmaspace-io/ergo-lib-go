package ergo

/*
   #include "ergo.h"
*/
import "C"
import "runtime"

// PreHeader represents a block header with the current SpendingTransaction, that can be predicted by a miner before its formation
type PreHeader interface {
	// Equals checks if provided PreHeader is same
	Equals(preHeader PreHeader) bool
	pointer() C.PreHeaderPtr
}

type preHeader struct {
	p C.PreHeaderPtr
}

func newPreHeader(h *preHeader) PreHeader {
	runtime.AddCleanup(h, finalizePreHeader, h.p)
	return h
}

// NewPreHeader creates PreHeader using data from BlockHeader
func NewPreHeader(header BlockHeader) PreHeader {
	var p C.PreHeaderPtr

	C.ergo_lib_preheader_from_block_header(header.pointer(), &p)
	runtime.KeepAlive(header)

	ph := &preHeader{p: p}

	return newPreHeader(ph)
}

func (h *preHeader) Equals(preHeader PreHeader) bool {
	res := C.ergo_lib_pre_header_eq(h.p, preHeader.pointer())
	runtime.KeepAlive(h)
	runtime.KeepAlive(preHeader)
	return bool(res)
}

func (h *preHeader) pointer() C.PreHeaderPtr {
	return h.p
}

func finalizePreHeader(p C.PreHeaderPtr) {
	C.ergo_lib_preheader_delete(p)
}
