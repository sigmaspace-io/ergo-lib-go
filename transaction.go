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

// TxId represents Transaction id
type TxId struct {
	p C.TxIdPtr
}

func newTxId(t *TxId) *TxId {
	runtime.AddCleanup(t, finalizeTxId, t.p)
	return t
}

// NewTxId creates TxId from hex-encoded string
func NewTxId(s string) (*TxId, error) {
	txIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(txIdStr))

	var p C.TxIdPtr

	errPtr := C.ergo_lib_tx_id_from_str(txIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &TxId{p}

	return newTxId(t), nil
}

// String returns TxId as string
func (t *TxId) String() (string, error) {
	var outTxIdStr *C.char

	errPtr := C.ergo_lib_tx_id_to_str(t.p, &outTxIdStr)
	runtime.KeepAlive(t)
	err := newError(errPtr)
	if err.isError() {
		return "", err.error()
	}
	defer C.ergo_lib_delete_string(outTxIdStr)

	return C.GoString(outTxIdStr), nil
}

// Equals checks if provided TxId is same
func (t *TxId) Equals(txId *TxId) bool {
	res := C.ergo_lib_tx_id_eq(t.p, txId.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(txId)
	return bool(res)
}

func (t *TxId) pointer() C.TxIdPtr {
	return t.p
}

func finalizeTxId(p C.TxIdPtr) {
	C.ergo_lib_tx_id_delete(p)
}

// CommitmentHint is a family of hints which are about a correspondence between a public image of a secret image and prover's commitment
// to randomness ("a" in a sigma protocol).
type CommitmentHint struct {
	p C.CommitmentHintPtr
}

func newCommitmentHint(c *CommitmentHint) *CommitmentHint {
	runtime.AddCleanup(c, finalizeCommitmentHint, c.p)
	return c
}

func (c *CommitmentHint) pointer() C.CommitmentHintPtr {
	return c.p
}

func finalizeCommitmentHint(p C.CommitmentHintPtr) {
	C.ergo_lib_commitment_hint_delete(p)
}

// HintsBag is a collection of CommitmentHint to be used by a prover
type HintsBag struct {
	p C.HintsBagPtr
}

func newHintsBag(h *HintsBag) *HintsBag {
	runtime.AddCleanup(h, finalizeHintsBag, h.p)
	return h
}

// NewHintsBag creates an empty HintsBag
func NewHintsBag() *HintsBag {
	var p C.HintsBagPtr
	C.ergo_lib_hints_bag_empty(&p)

	h := &HintsBag{p: p}
	return newHintsBag(h)
}

// Add adds CommitmentHint to the bag
func (h *HintsBag) Add(hint *CommitmentHint) {
	C.ergo_lib_hints_bag_add_commitment(h.p, hint.pointer())
	runtime.KeepAlive(h)
	runtime.KeepAlive(hint)
}

// Len returns the length of the HintsBag
func (h *HintsBag) Len() int {
	res := C.ergo_lib_hints_bag_len(h.p)
	return int(res)
}

// Get returns the CommitmentHint at the provided index if it exists
func (h *HintsBag) Get(index int) (*CommitmentHint, error) {
	var p C.CommitmentHintPtr

	res := C.ergo_lib_hints_bag_get(h.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(h)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		c := &CommitmentHint{p: p}
		return newCommitmentHint(c), nil
	}

	return nil, nil
}

// All returns an iterator over all CommitmentHint inside the collection
func (h *HintsBag) All() iter.Seq2[int, *CommitmentHint] {
	return func(yield func(int, *CommitmentHint) bool) {
		for i := 0; i < h.Len(); i++ {
			tk, err := h.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(h)
	}
}

func (h *HintsBag) pointer() C.HintsBagPtr {
	return h.p
}

func finalizeHintsBag(p C.HintsBagPtr) {
	C.ergo_lib_hints_bag_delete(p)
}

type TransactionHintsBag struct {
	p C.TransactionHintsBagPtr
}

func newTransactionHintsBag(t *TransactionHintsBag) *TransactionHintsBag {
	runtime.AddCleanup(t, finalizeTransactionHintsBag, t.p)
	return t
}

// NewTransactionHintsBag creates empty TransactionHintsBag
func NewTransactionHintsBag() *TransactionHintsBag {
	var p C.TransactionHintsBagPtr
	C.ergo_lib_transaction_hints_bag_empty(&p)

	t := &TransactionHintsBag{p: p}

	return newTransactionHintsBag(t)
}

// AddHintsForInput adds hints for Input
func (t *TransactionHintsBag) AddHintsForInput(index uint32, hintsBag *HintsBag) {
	C.ergo_lib_transaction_hints_bag_add_hints_for_input(t.p, C.uintptr_t(index), hintsBag.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(hintsBag)
}

// AllHintsForInput gets HintsBag corresponding to Input index
func (t *TransactionHintsBag) AllHintsForInput(index uint32) *HintsBag {
	var p C.HintsBagPtr
	C.ergo_lib_transaction_hints_bag_all_hints_for_input(t.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(t)
	h := &HintsBag{p: p}
	return newHintsBag(h)
}

func (t *TransactionHintsBag) pointer() C.TransactionHintsBagPtr {
	return t.p
}

func finalizeTransactionHintsBag(p C.TransactionHintsBagPtr) {
	C.ergo_lib_transaction_hints_bag_delete(p)
}

// ExtractHintsFromSignedTransaction extracts hints from signed Transaction
func ExtractHintsFromSignedTransaction(
	transaction *Transaction,
	stateContext *StateContext,
	boxesToSpend *Boxes,
	dataBoxes *Boxes,
	realPropositions *Propositions,
	simulatedPropositions *Propositions) (*TransactionHintsBag, error) {
	var p C.TransactionHintsBagPtr

	errPtr := C.ergo_lib_transaction_extract_hints(
		transaction.pointer(),
		stateContext.pointer(),
		boxesToSpend.pointer(),
		dataBoxes.pointer(),
		realPropositions.pointer(),
		simulatedPropositions.pointer(),
		&p)
	runtime.KeepAlive(transaction)
	runtime.KeepAlive(stateContext)
	runtime.KeepAlive(boxesToSpend)
	runtime.KeepAlive(dataBoxes)
	runtime.KeepAlive(realPropositions)
	runtime.KeepAlive(simulatedPropositions)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	th := &TransactionHintsBag{p: p}

	return newTransactionHintsBag(th), nil
}

// UnsignedTransaction represents an unsigned Transaction (inputs without proofs)
type UnsignedTransaction struct {
	p C.UnsignedTransactionPtr
}

func newUnsignedTransaction(u *UnsignedTransaction) *UnsignedTransaction {
	runtime.AddCleanup(u, finalizeUnsignedTransaction, u.p)
	return u
}

// NewUnsignedTransactionFromJson parse UnsignedTransaction from JSON. Supports Ergo Node/Explorer API and Box values and Token amount encoded as strings.
func NewUnsignedTransactionFromJson(json string) (*UnsignedTransaction, error) {
	unsTxJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(unsTxJsonStr))

	var p C.UnsignedTransactionPtr

	errPtr := C.ergo_lib_unsigned_tx_from_json(unsTxJsonStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ut := &UnsignedTransaction{p: p}
	return newUnsignedTransaction(ut), nil
}

// TxId returns TxId for this UnsignedTransaction
func (u *UnsignedTransaction) TxId() *TxId {
	var p C.TxIdPtr
	C.ergo_lib_unsigned_tx_id(u.p, &p)
	runtime.KeepAlive(u)
	ti := &TxId{p: p}
	return newTxId(ti)
}

// UnsignedInputs returns UnsignedInputs for this UnsignedTransaction
func (u *UnsignedTransaction) UnsignedInputs() *UnsignedInputs {
	var p C.UnsignedInputsPtr
	C.ergo_lib_unsigned_tx_inputs(u.p, &p)
	runtime.KeepAlive(u)
	ui := &UnsignedInputs{p: p}
	return newUnsignedInputs(ui)
}

// DataInputs returns DataInputs for this UnsignedTransaction
func (u *UnsignedTransaction) DataInputs() *DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_unsigned_tx_data_inputs(u.p, &p)
	runtime.KeepAlive(u)
	di := &DataInputs{p: p}
	return newDataInputs(di)
}

// OutputCandidates returns BoxCandidates for this UnsignedTransaction
func (u *UnsignedTransaction) OutputCandidates() *BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_unsigned_tx_output_candidates(u.p, &p)
	runtime.KeepAlive(u)
	bc := &BoxCandidates{p: p}
	return newBoxCandidates(bc)
}

// Json returns json representation of UnsignedTransaction as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
func (u *UnsignedTransaction) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_unsigned_tx_to_json(u.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(u)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// JsonEIP12 returns json representation of UnsignedTransaction as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
func (u *UnsignedTransaction) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_unsigned_tx_to_json_eip12(u.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(u)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (u *UnsignedTransaction) pointer() C.UnsignedTransactionPtr {
	return u.p
}

func finalizeUnsignedTransaction(p C.UnsignedTransactionPtr) {
	C.ergo_lib_unsigned_tx_delete(p)
}

// Transaction is an atomic state transition operation. It destroys Boxes from the state
// and creates new ones. If transaction is spending boxes protected by some non-trivial scripts,
// its inputs should also contain proof of spending correctness - context extension (user-defined
// key-value map) and data inputs (links to existing boxes in the state) that may be used during
// script reduction to crypto, signatures that satisfies the remaining cryptographic protection
// of the script.
// Transactions are not encrypted, so it is possible to browse and view every transaction ever
// collected into a block.
type Transaction struct {
	p C.TransactionPtr
}

func newTransaction(t *Transaction) *Transaction {
	runtime.AddCleanup(t, finalizeTransaction, t.p)
	return t
}

// NewTransaction creates Transaction from UnsignedTransaction and an array of proofs in the same order as
// UnsignedTransaction inputs with empty proof indicated with empty ByteArray
func NewTransaction(unsignedTx *UnsignedTransaction, proofs *ByteArrays) (*Transaction, error) {
	var p C.TransactionPtr

	errPtr := C.ergo_lib_tx_from_unsigned_tx(unsignedTx.pointer(), proofs.pointer(), &p)
	runtime.KeepAlive(unsignedTx)
	runtime.KeepAlive(proofs)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &Transaction{p: p}
	return newTransaction(t), nil
}

// NewTransactionFromJson parse Transaction from JSON. Supports Ergo Node/Explorer API and box values and token amount encoded as strings.
func NewTransactionFromJson(json string) (*Transaction, error) {
	txJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(txJsonStr))

	var p C.TransactionPtr

	errPtr := C.ergo_lib_tx_from_json(txJsonStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &Transaction{p: p}
	return newTransaction(t), nil
}

// TxId returns TxId for this Transaction
func (t *Transaction) TxId() *TxId {
	var p C.TxIdPtr
	C.ergo_lib_tx_id(t.p, &p)
	runtime.KeepAlive(t)
	ti := &TxId{p: p}
	return newTxId(ti)
}

// Inputs returns Inputs for this Transaction
func (t *Transaction) Inputs() *Inputs {
	var p C.InputsPtr
	C.ergo_lib_tx_inputs(t.p, &p)
	runtime.KeepAlive(t)
	i := &Inputs{p: p}
	return newInputs(i)
}

// DataInputs returns DataInputs for this Transaction
func (t *Transaction) DataInputs() *DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_tx_data_inputs(t.p, &p)
	runtime.KeepAlive(t)
	di := &DataInputs{p: p}
	return newDataInputs(di)
}

// OutputCandidates returns BoxCandidates for this Transaction
func (t *Transaction) OutputCandidates() *BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_tx_output_candidates(t.p, &p)
	runtime.KeepAlive(t)
	bc := &BoxCandidates{p: p}
	return newBoxCandidates(bc)
}

// Outputs returns Boxes for this Transaction
func (t *Transaction) Outputs() *Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_tx_outputs(t.p, &p)
	runtime.KeepAlive(t)
	b := &Boxes{p: p}
	return newBoxes(b)
}

// Json returns json representation of Transaction as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
func (t *Transaction) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_tx_to_json(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// JsonEIP12 returns json representation of Transaction as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
func (t *Transaction) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_tx_to_json_eip12(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// Validate validates the current Transaction
func (t *Transaction) Validate(stateContext *StateContext, boxesToSpent *Boxes, dataBoxes *Boxes) error {
	errPtr := C.ergo_lib_tx_validate(t.p, stateContext.pointer(), boxesToSpent.pointer(), dataBoxes.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(dataBoxes)
	runtime.KeepAlive(stateContext)
	runtime.KeepAlive(boxesToSpent)
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (t *Transaction) pointer() C.TransactionPtr {
	return t.p
}

func finalizeTransaction(p C.TransactionPtr) {
	C.ergo_lib_tx_delete(p)
}
