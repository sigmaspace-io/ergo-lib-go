package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"iter"
	"runtime"
	"unsafe"
)

// BlockHeader represents data of the block header available in Sigma proposition
type BlockHeader struct {
	p C.BlockHeaderPtr
}

func newBlockHeader(b *BlockHeader) *BlockHeader {
	runtime.AddCleanup(b, finalizeBlockHeader, b.p)
	return b
}

// NewBlockHeader creates a new BlockHeader from block header array JSON (Node API)
func NewBlockHeader(json string) (*BlockHeader, error) {
	blockHeaderJson := C.CString(json)
	defer C.free(unsafe.Pointer(blockHeaderJson))

	var p C.BlockHeaderPtr

	errPtr := C.ergo_lib_block_header_from_json(blockHeaderJson, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &BlockHeader{p: p}

	return newBlockHeader(b), nil
}

// BlockId returns the BlockId of the BlockHeader
func (b *BlockHeader) BlockId() *BlockId {
	var p C.BlockIdPtr

	C.ergo_lib_block_header_id(b.p, &p)
	runtime.KeepAlive(b)

	bi := &BlockId{p: p}

	return newBlockId(bi)
}

// Equals checks if provided BlockHeader is same
func (b *BlockHeader) Equals(blockHeader *BlockHeader) bool {
	res := C.ergo_lib_block_header_eq(b.p, blockHeader.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(blockHeader)
	return bool(res)
}

func (b *BlockHeader) pointer() C.BlockHeaderPtr {
	return b.p
}

func finalizeBlockHeader(p C.BlockHeaderPtr) {
	C.ergo_lib_block_header_delete(p)
}

// BlockId represents the id of a BlockHeader
type BlockId struct {
	p C.BlockIdPtr
}

func newBlockId(b *BlockId) *BlockId {
	runtime.AddCleanup(b, finalizeBlockId, b.p)
	return b
}

// NewBlockId creates a new BlockId from hex-encoded string
func NewBlockId(s string) (*BlockId, error) {
	blockIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(blockIdStr))

	var p C.BlockIdPtr

	errPtr := C.ergo_lib_block_id_from_str(blockIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &BlockId{p: p}

	return newBlockId(b), nil
}

// Equals checks if provided BlockId is same
func (b *BlockId) Equals(blockId *BlockId) bool {
	res := C.ergo_lib_block_id_eq(b.p, blockId.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(blockId)
	return bool(res)
}

func (b *BlockId) pointer() C.BlockIdPtr {
	return b.p
}

func finalizeBlockId(p C.BlockIdPtr) {
	C.ergo_lib_block_id_delete(p)
}

// BlockHeaders an ordered collection of BlockHeader
type BlockHeaders struct {
	p C.BlockHeadersPtr
}

func newBlockHeaders(b *BlockHeaders) *BlockHeaders {
	runtime.AddCleanup(b, finalizeBlockHeaders, b.p)
	return b
}

// NewBlockHeaders creates an empty BlockHeaders collection
func NewBlockHeaders() *BlockHeaders {
	var p C.BlockHeadersPtr
	C.ergo_lib_block_headers_new(&p)
	b := &BlockHeaders{p: p}

	return newBlockHeaders(b)
}

// Len returns the length of the collection
func (b *BlockHeaders) Len() int {
	res := C.ergo_lib_block_headers_len(b.p)
	runtime.KeepAlive(b)
	return int(res)
}

// Get returns the BlockHeader at the provided index if it exists
func (b *BlockHeaders) Get(index int) (*BlockHeader, error) {
	var p C.BlockHeaderPtr

	res := C.ergo_lib_block_headers_get(b.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		bh := &BlockHeader{p: p}
		return newBlockHeader(bh), nil
	}

	return nil, nil
}

// Add adds provided BlockHeader to the end of the collection
func (b *BlockHeaders) Add(blockHeader *BlockHeader) {
	C.ergo_lib_block_headers_add(blockHeader.pointer(), b.p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(blockHeader)
}

// All returns an iterator over all BlockHeader inside the collection
func (b *BlockHeaders) All() iter.Seq2[int, *BlockHeader] {
	return func(yield func(int, *BlockHeader) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(b)
	}
}

func (b *BlockHeaders) pointer() C.BlockHeadersPtr {
	return b.p
}

func finalizeBlockHeaders(p C.BlockHeadersPtr) {
	C.ergo_lib_block_headers_delete(p)
}

// BlockIds an ordered collection of BlockId
type BlockIds struct {
	p C.BlockIdsPtr
}

func newBlockIds(b *BlockIds) *BlockIds {
	runtime.AddCleanup(b, finalizeBlockIds, b.p)
	return b
}

// NewBlockIds creates an empty BlockIds collection
func NewBlockIds() *BlockIds {
	var p C.BlockIdsPtr
	C.ergo_lib_block_ids_new(&p)

	b := &BlockIds{p: p}

	return newBlockIds(b)
}

// Len returns the length of the collection
func (b *BlockIds) Len() int {
	res := C.ergo_lib_block_ids_len(b.p)
	runtime.KeepAlive(b)
	return int(res)
}

// Get returns the BlockId at the provided index if it exists
func (b *BlockIds) Get(index int) (*BlockId, error) {
	var p C.BlockIdPtr

	res := C.ergo_lib_block_ids_get(b.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		bi := &BlockId{p: p}
		return newBlockId(bi), nil
	}

	return nil, nil
}

// Add adds provided BlockId to the end of the collection
func (b *BlockIds) Add(blockId *BlockId) {
	C.ergo_lib_block_ids_add(blockId.pointer(), b.p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(blockId)
}

// All returns an iterator over all BlockId inside the collection
func (b *BlockIds) All() iter.Seq2[int, *BlockId] {
	return func(yield func(int, *BlockId) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(b)
	}
}

func finalizeBlockIds(p C.BlockIdsPtr) {
	C.ergo_lib_block_ids_delete(p)
}
