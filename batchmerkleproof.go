package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type BatchMerkleProof struct {
	p C.BatchMerkleProofPtr
}

func newBatchMerkleProof(b *BatchMerkleProof) *BatchMerkleProof {
	runtime.AddCleanup(b, finalizeBatchMerkleProof, b.p)
	return b
}

func NewBatchMerkleProof(json string) (*BatchMerkleProof, error) {
	jsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(jsonStr))

	var p C.BatchMerkleProofPtr
	errPtr := C.ergo_lib_batch_merkle_proof_from_json(jsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	b := &BatchMerkleProof{p: p}
	return newBatchMerkleProof(b), nil
}

func (b *BatchMerkleProof) Valid(expectedRoot []byte) bool {
	byteData := C.CBytes(expectedRoot)
	defer C.free(unsafe.Pointer(byteData))
	res := C.ergo_lib_batch_merkle_proof_valid(b.p, (*C.uchar)(byteData), C.uintptr_t(len(expectedRoot)))
	runtime.KeepAlive(b)
	return bool(res)
}

func finalizeBatchMerkleProof(p C.BatchMerkleProofPtr) {
	C.ergo_lib_batch_merkle_proof_delete(p)
}
