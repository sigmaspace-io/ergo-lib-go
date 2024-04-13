package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type nodeSide uint8

const (
	// Left represents left side the node is on in the merkle tree
	Left nodeSide = 0
	// Right represents right side the node is on in the merkle tree
	Right nodeSide = 1
)

type MerkleProof interface {
	// AddNode adds a new node and it's hash to the MerkleProof. Hash must be 32 bytes in size
	AddNode(hash []byte, side nodeSide) error
	// Valid validates the MerkleProof against the provided root hash
	Valid(expectedRoot []byte) bool
	// ValidBase16 validates the MerkleProof against the provided base16 root hash
	ValidBase16(expectedRoot string) bool
}

type merkleProof struct {
	p C.MerkleProofPtr
}

func newMerkleProof(m *merkleProof) MerkleProof {
	runtime.SetFinalizer(m, finalizeMerkleProof)
	return m
}

func NewMerkleProof(leafData []byte) (MerkleProof, error) {
	byteData := C.CBytes(leafData)
	defer C.free(unsafe.Pointer(byteData))
	var p C.MerkleProofPtr

	errPtr := C.ergo_merkle_proof_new((*C.uchar)(byteData), C.uintptr_t(len(leafData)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	m := &merkleProof{p: p}
	return newMerkleProof(m), nil
}

func NewMerkleProofFromJson(json string) (MerkleProof, error) {
	merkleProofJson := C.CString(json)
	defer C.free(unsafe.Pointer(merkleProofJson))

	var p C.MerkleProofPtr

	errPtr := C.ergo_merkle_proof_from_json(merkleProofJson, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	m := &merkleProof{p: p}

	return newMerkleProof(m), nil
}

func (m *merkleProof) AddNode(hash []byte, side nodeSide) error {
	byteData := C.CBytes(hash)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_merkle_proof_add_node(m.p, (*C.uchar)(byteData), C.uintptr_t(len(hash)), C.uchar(side))
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (m *merkleProof) Valid(expectedRoot []byte) bool {
	byteData := C.CBytes(expectedRoot)
	defer C.free(unsafe.Pointer(byteData))
	res := C.ergo_merkle_proof_valid(m.p, (*C.uchar)(byteData), C.uintptr_t(len(expectedRoot)))
	return bool(res)
}

func (m *merkleProof) ValidBase16(expectedRoot string) bool {
	rootStr := C.CString(expectedRoot)
	defer C.free(unsafe.Pointer(rootStr))
	var res C.bool
	errPtr := C.ergo_merkle_proof_valid_base16(m.p, rootStr, &res)
	err := newError(errPtr)
	if err.isError() {
		return false
	}
	return bool(res)
}

func finalizeMerkleProof(m *merkleProof) {
	C.ergo_merkle_proof_delete(m.p)
}
