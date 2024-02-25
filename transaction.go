package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type TxId interface {
	ToString() (string, error)
	pointer() C.TxIdPtr
}

type txId struct {
	p C.TxIdPtr
}

func newTxId(t *txId) TxId {
	runtime.SetFinalizer(t, finalizeTxId)
	return t
}

func NewTxId(s string) (TxId, error) {
	txIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(txIdStr))

	var p C.TxIdPtr

	errPtr := C.ergo_lib_tx_id_from_str(txIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &txId{p}

	return newTxId(t), nil
}

func (t *txId) ToString() (string, error) {
	var outTxIdStr *C.char

	errPtr := C.ergo_lib_tx_id_to_str(t.p, &outTxIdStr)
	err := newError(errPtr)
	if err.isError() {
		return "", err.error()
	}
	defer C.ergo_lib_delete_string(outTxIdStr)

	return C.GoString(outTxIdStr), nil
}

func (t *txId) pointer() C.TxIdPtr {
	return t.p
}

func finalizeTxId(t *txId) {
	C.ergo_lib_tx_id_delete(t.p)
}
