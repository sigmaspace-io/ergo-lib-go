package ergo

/*
   #include "ergo.h"
*/
import "C"
import "runtime"

// PreHeader represents a block header with the current SpendingTransaction, that can be predicted by a miner before its formation
type PreHeader interface {
	pointer() C.PreHeaderPtr
}

type preHeader struct {
	p C.PreHeaderPtr
}

func newPreHeader(h *preHeader) PreHeader {
	runtime.SetFinalizer(h, finalizePreHeader)
	return h
}

// NewPreHeader creates PreHeader using data from BlockHeader
func NewPreHeader(header BlockHeader) PreHeader {
	var p C.PreHeaderPtr

	C.ergo_lib_preheader_from_block_header(header.pointer(), &p)

	ph := &preHeader{p: p}

	return newPreHeader(ph)
}

func (h *preHeader) pointer() C.PreHeaderPtr {
	return h.p
}

func finalizePreHeader(h *preHeader) {
	C.ergo_lib_preheader_delete(h.p)
}
