package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type BatchMerkleProof interface {
	Valid(expectedRoot []byte) bool
}

type batchMerkleProof struct {
	p C.BatchMerkleProofPtr
}

func newBatchMerkleProof(b *batchMerkleProof) BatchMerkleProof {
	runtime.SetFinalizer(b, finalizeBatchMerkleProof)
	return b
}

func NewBatchMerkleProof(json string) (BatchMerkleProof, error) {
	jsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(jsonStr))

	var p C.BatchMerkleProofPtr
	errPtr := C.ergo_lib_batch_merkle_proof_from_json(jsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	b := &batchMerkleProof{p: p}
	return newBatchMerkleProof(b), nil
}

func (b *batchMerkleProof) Valid(expectedRoot []byte) bool {
	byteData := C.CBytes(expectedRoot)
	defer C.free(unsafe.Pointer(byteData))
	res := C.ergo_lib_batch_merkle_proof_valid(b.p, (*C.uchar)(byteData), C.ulong(len(expectedRoot)))
	return bool(res)
}

func finalizeBatchMerkleProof(b *batchMerkleProof) {
	C.ergo_lib_batch_merkle_proof_delete(b.p)
}
