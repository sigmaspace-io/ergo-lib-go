package ergo

/*
   #include "ergo.h"
*/
import "C"
import "runtime"

// PreHeader represents a block header with the current SpendingTransaction, that can be predicted by a miner before its formation
type PreHeader struct {
	p C.PreHeaderPtr
}

func newPreHeader(h *PreHeader) *PreHeader {
	runtime.AddCleanup(h, finalizePreHeader, h.p)
	return h
}

// NewPreHeader creates PreHeader using data from BlockHeader
func NewPreHeader(header *BlockHeader) *PreHeader {
	var p C.PreHeaderPtr

	C.ergo_lib_preheader_from_block_header(header.pointer(), &p)
	runtime.KeepAlive(header)

	ph := &PreHeader{p: p}

	return newPreHeader(ph)
}

// Equals checks if provided PreHeader is same
func (h *PreHeader) Equals(preHeader *PreHeader) bool {
	res := C.ergo_lib_pre_header_eq(h.p, preHeader.pointer())
	runtime.KeepAlive(h)
	runtime.KeepAlive(preHeader)
	return bool(res)
}

func (h *PreHeader) pointer() C.PreHeaderPtr {
	return h.p
}

func finalizePreHeader(p C.PreHeaderPtr) {
	C.ergo_lib_preheader_delete(p)
}
