package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ReducedTransaction represents reduced transaction, i.e. unsigned transaction where each unsigned input
// is augmented with ReducedInput which contains a script reduction result.
// After an unsigned transaction is reduced it can be signed without context.
// Thus, it can be serialized and transferred for example to Cold Wallet and signed
// in an environment where secrets are known.
// see EIP-19 for more details -
// https://github.com/ergoplatform/eips/blob/f280890a4163f2f2e988a0091c078e36912fc531/eip-0019.md
type ReducedTransaction interface {
	// UnsignedTransaction returns the UnsignedTransaction
	UnsignedTransaction() UnsignedTransaction
	pointer() C.ReducedTransactionPtr
}

type reducedTransaction struct {
	p C.ReducedTransactionPtr
}

func newReducedTransaction(r *reducedTransaction) ReducedTransaction {
	runtime.SetFinalizer(r, finalizeReducedTransaction)
	return r
}

// NewReducedTransaction creates a ReducedTransaction i.e unsigned transaction where each unsigned input
// is augmented with ReducedInput which contains a script reduction result
func NewReducedTransaction(unsignedTx UnsignedTransaction, boxesToSpent Boxes, dataBoxes Boxes, stateContext StateContext) (ReducedTransaction, error) {
	var p C.ReducedTransactionPtr

	errPtr := C.ergo_lib_reduced_tx_from_unsigned_tx(unsignedTx.pointer(), boxesToSpent.pointer(), dataBoxes.pointer(), stateContext.pointer(), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	r := &reducedTransaction{p: p}
	return newReducedTransaction(r), nil
}

func (r *reducedTransaction) UnsignedTransaction() UnsignedTransaction {
	var p C.UnsignedTransactionPtr
	C.ergo_lib_reduced_tx_unsigned_tx(r.p, &p)
	ut := &unsignedTransaction{p: p}
	return newUnsignedTransaction(ut)
}

func (r *reducedTransaction) pointer() C.ReducedTransactionPtr {
	return r.p
}

func finalizeReducedTransaction(r *reducedTransaction) {
	C.ergo_lib_reduced_tx_delete(r.p)
}

// Propositions list(public keys)
type Propositions interface {
	// Add adds new proposition
	Add(bytes []byte) error
	pointer() C.PropositionsPtr
}

type propositions struct {
	p C.PropositionsPtr
}

func newPropositions(p *propositions) Propositions {
	runtime.SetFinalizer(p, finalizePropositions)
	return p
}

// NewPropositions creates empty proposition holder
func NewPropositions() Propositions {
	var p C.PropositionsPtr
	C.ergo_lib_propositions_new(&p)
	prop := &propositions{p: p}
	return newPropositions(prop)
}

func (p *propositions) Add(bytes []byte) error {
	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_lib_propositions_add_proposition_from_bytes(p.p, (*C.uchar)(byteData), C.uintptr_t(len(bytes)))
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (p *propositions) pointer() C.PropositionsPtr {
	return p.p
}

func finalizePropositions(p *propositions) {
	C.ergo_lib_propositions_delete(p.p)
}
