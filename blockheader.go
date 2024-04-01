package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// BlockHeader represents data of the block header available in Sigma proposition
type BlockHeader interface {
	// BlockId returns the BlockId of the BlockHeader
	BlockId() BlockId
	pointer() C.BlockHeaderPtr
}

type blockHeader struct {
	p C.BlockHeaderPtr
}

func newBlockHeader(b *blockHeader) BlockHeader {
	runtime.SetFinalizer(b, finalizeBlockHeader)
	return b
}

// NewBlockHeader creates a new BlockHeader from block header array JSON (Node API)
func NewBlockHeader(json string) (BlockHeader, error) {
	blockHeaderJson := C.CString(json)
	defer C.free(unsafe.Pointer(blockHeaderJson))

	var p C.BlockHeaderPtr

	errPtr := C.ergo_lib_block_header_from_json(blockHeaderJson, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &blockHeader{p: p}

	return newBlockHeader(b), nil
}

func (b *blockHeader) BlockId() BlockId {
	var p C.BlockIdPtr

	C.ergo_lib_block_header_id(b.p, &p)

	bi := &blockId{p: p}

	return newBlockId(bi)
}

func (b *blockHeader) pointer() C.BlockHeaderPtr {
	return b.p
}

func finalizeBlockHeader(b *blockHeader) {
	C.ergo_lib_block_header_delete(b.p)
}

// BlockId represents the id of a BlockHeader
type BlockId interface {
	pointer() C.BlockIdPtr
}

type blockId struct {
	p C.BlockIdPtr
}

func newBlockId(b *blockId) BlockId {
	runtime.SetFinalizer(b, finalizeBlockId)
	return b
}

// NewBlockId creates a new BlockId from hex-encoded string
func NewBlockId(s string) (BlockId, error) {
	blockIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(blockIdStr))

	var p C.BlockIdPtr

	errPtr := C.ergo_lib_block_id_from_str(blockIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &blockId{p: p}

	return newBlockId(b), nil
}

func (b *blockId) pointer() C.BlockIdPtr {
	return b.p
}

func finalizeBlockId(b *blockId) {
	C.ergo_lib_block_id_delete(b.p)
}

// BlockHeaders an ordered collection of BlockHeader
type BlockHeaders interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the BlockHeader at the provided index if it exists
	Get(index uint32) (BlockHeader, error)
	// Add adds provided BlockHeader to the end of the collection
	Add(blockHeader BlockHeader)
	pointer() C.BlockHeadersPtr
}

type blockHeaders struct {
	p C.BlockHeadersPtr
}

func newBlockHeaders(b *blockHeaders) BlockHeaders {
	runtime.SetFinalizer(b, finalizeBlockHeaders)
	return b
}

// NewBlockHeaders creates an empty BlockHeaders collection
func NewBlockHeaders() BlockHeaders {
	var p C.BlockHeadersPtr
	C.ergo_lib_block_headers_new(&p)
	b := &blockHeaders{p: p}

	return newBlockHeaders(b)
}

func (b *blockHeaders) Len() uint32 {
	res := C.ergo_lib_block_headers_len(b.p)
	return uint32(res)
}

func (b *blockHeaders) Get(index uint32) (BlockHeader, error) {
	var p C.BlockHeaderPtr

	res := C.ergo_lib_block_headers_get(b.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		bh := &blockHeader{p: p}
		return newBlockHeader(bh), nil
	}

	return nil, nil
}

func (b *blockHeaders) Add(blockHeader BlockHeader) {
	C.ergo_lib_block_headers_add(blockHeader.pointer(), b.p)
}

func (b *blockHeaders) pointer() C.BlockHeadersPtr {
	return b.p
}

func finalizeBlockHeaders(b *blockHeaders) {
	C.ergo_lib_block_headers_delete(b.p)
}

// BlockIds an ordered collection of BlockId
type BlockIds interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the BlockId at the provided index if it exists
	Get(index uint32) (BlockId, error)
	// Add adds provided BlockId to the end of the collection
	Add(blockId BlockId)
}

type blockIds struct {
	p C.BlockIdsPtr
}

func newBlockIds(b *blockIds) BlockIds {
	runtime.SetFinalizer(b, finalizeBlockIds)
	return b
}

// NewBlockIds creates an empty BlockIds collection
func NewBlockIds() BlockIds {
	var p C.BlockIdsPtr
	C.ergo_lib_block_ids_new(&p)

	b := &blockIds{p: p}

	return newBlockIds(b)
}

func (b *blockIds) Len() uint32 {
	res := C.ergo_lib_block_ids_len(b.p)
	return uint32(res)
}

func (b *blockIds) Get(index uint32) (BlockId, error) {
	var p C.BlockIdPtr

	res := C.ergo_lib_block_ids_get(b.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		bi := &blockId{p: p}
		return newBlockId(bi), nil
	}

	return nil, nil
}

func (b *blockIds) Add(blockId BlockId) {
	C.ergo_lib_block_ids_add(blockId.pointer(), b.p)
}

func finalizeBlockIds(b *blockIds) {
	C.ergo_lib_block_ids_delete(b.p)
}
