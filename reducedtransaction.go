package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ReducedTransaction represents reduced Transaction, i.e. unsigned Transaction where each unsigned Input
// is augmented with ReducedInput which contains a script reduction result.
// After an unsigned Transaction is reduced it can be signed without context.
// Thus, it can be serialized and transferred for example to Cold Wallet and signed
// in an environment where secrets are known.
// see EIP-19 for more details -
// https://github.com/ergoplatform/eips/blob/f280890a4163f2f2e988a0091c078e36912fc531/eip-0019.md
type ReducedTransaction struct {
	p C.ReducedTransactionPtr
}

func newReducedTransaction(r *ReducedTransaction) *ReducedTransaction {
	runtime.AddCleanup(r, finalizeReducedTransaction, r.p)
	return r
}

// NewReducedTransaction creates a ReducedTransaction i.e unsigned Transaction where each unsigned Input
// is augmented with ReducedInput which contains a script reduction result
func NewReducedTransaction(unsignedTx *UnsignedTransaction, boxesToSpent *Boxes, dataBoxes *Boxes, stateContext *StateContext) (*ReducedTransaction, error) {
	var p C.ReducedTransactionPtr

	errPtr := C.ergo_lib_reduced_tx_from_unsigned_tx(unsignedTx.pointer(), boxesToSpent.pointer(), dataBoxes.pointer(), stateContext.pointer(), &p)
	runtime.KeepAlive(unsignedTx)
	runtime.KeepAlive(boxesToSpent)
	runtime.KeepAlive(dataBoxes)
	runtime.KeepAlive(stateContext)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	r := &ReducedTransaction{p: p}
	return newReducedTransaction(r), nil
}

// UnsignedTransaction returns the UnsignedTransaction
func (r *ReducedTransaction) UnsignedTransaction() *UnsignedTransaction {
	var p C.UnsignedTransactionPtr
	C.ergo_lib_reduced_tx_unsigned_tx(r.p, &p)
	runtime.KeepAlive(r)
	ut := &UnsignedTransaction{p: p}
	return newUnsignedTransaction(ut)
}

func (r *ReducedTransaction) pointer() C.ReducedTransactionPtr {
	return r.p
}

func finalizeReducedTransaction(p C.ReducedTransactionPtr) {
	C.ergo_lib_reduced_tx_delete(p)
}

// Propositions list(public keys)
type Propositions struct {
	p C.PropositionsPtr
}

func newPropositions(p *Propositions) *Propositions {
	runtime.AddCleanup(p, finalizePropositions, p.p)
	return p
}

// NewPropositions creates empty proposition holder
func NewPropositions() *Propositions {
	var p C.PropositionsPtr
	C.ergo_lib_propositions_new(&p)
	prop := &Propositions{p: p}
	return newPropositions(prop)
}

// Add adds new proposition
func (p *Propositions) Add(bytes []byte) error {
	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_lib_propositions_add_proposition_from_bytes(p.p, (*C.uchar)(byteData), C.uintptr_t(len(bytes)))
	runtime.KeepAlive(p)
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (p *Propositions) pointer() C.PropositionsPtr {
	return p.p
}

func finalizePropositions(p C.PropositionsPtr) {
	C.ergo_lib_propositions_delete(p)
}
