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
	// Left represents left side the node is on in the merkle Tree
	Left nodeSide = 0
	// Right represents right side the node is on in the merkle Tree
	Right nodeSide = 1
)

type MerkleProof struct {
	p C.MerkleProofPtr
}

func newMerkleProof(m *MerkleProof) *MerkleProof {
	runtime.AddCleanup(m, finalizeMerkleProof, m.p)
	return m
}

func NewMerkleProof(leafData []byte) (*MerkleProof, error) {
	byteData := C.CBytes(leafData)
	defer C.free(unsafe.Pointer(byteData))
	var p C.MerkleProofPtr

	errPtr := C.ergo_merkle_proof_new((*C.uchar)(byteData), C.uintptr_t(len(leafData)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	m := &MerkleProof{p: p}
	return newMerkleProof(m), nil
}

func NewMerkleProofFromJson(json string) (*MerkleProof, error) {
	merkleProofJson := C.CString(json)
	defer C.free(unsafe.Pointer(merkleProofJson))

	var p C.MerkleProofPtr

	errPtr := C.ergo_merkle_proof_from_json(merkleProofJson, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	m := &MerkleProof{p: p}

	return newMerkleProof(m), nil
}

// AddNode adds a new node and it's hash to the MerkleProof. Hash must be 32 bytes in size
func (m *MerkleProof) AddNode(hash []byte, side nodeSide) error {
	byteData := C.CBytes(hash)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_merkle_proof_add_node(m.p, (*C.uchar)(byteData), C.uintptr_t(len(hash)), C.uchar(side))
	runtime.KeepAlive(m)
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

// Valid validates the MerkleProof against the provided root hash
func (m *MerkleProof) Valid(expectedRoot []byte) bool {
	byteData := C.CBytes(expectedRoot)
	defer C.free(unsafe.Pointer(byteData))
	res := C.ergo_merkle_proof_valid(m.p, (*C.uchar)(byteData), C.uintptr_t(len(expectedRoot)))
	runtime.KeepAlive(m)
	return bool(res)
}

// ValidBase16 validates the MerkleProof against the provided base16 root hash
func (m *MerkleProof) ValidBase16(expectedRoot string) bool {
	rootStr := C.CString(expectedRoot)
	defer C.free(unsafe.Pointer(rootStr))
	var res C.bool
	errPtr := C.ergo_merkle_proof_valid_base16(m.p, rootStr, &res)
	runtime.KeepAlive(m)
	err := newError(errPtr)
	if err.isError() {
		return false
	}
	return bool(res)
}

func finalizeMerkleProof(p C.MerkleProofPtr) {
	C.ergo_merkle_proof_delete(p)
}
