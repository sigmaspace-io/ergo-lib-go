package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// TxId represents transaction id
type TxId interface {
	// String returns TxId as string
	String() (string, error)
	pointer() C.TxIdPtr
}

type txId struct {
	p C.TxIdPtr
}

func newTxId(t *txId) TxId {
	runtime.SetFinalizer(t, finalizeTxId)
	return t
}

// NewTxId creates TxId from hex-encoded string
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

func (t *txId) String() (string, error) {
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

// CommitmentHint is a family of hints which are about a correspondence between a public image of a secret image and prover's commitment
// to randomness ("a" in a sigma protocol).
type CommitmentHint interface {
	pointer() C.CommitmentHintPtr
}

type commitmentHint struct {
	p C.CommitmentHintPtr
}

func newCommitmentHint(c *commitmentHint) CommitmentHint {
	runtime.SetFinalizer(c, finalizeCommitmentHint)
	return c
}

func (c *commitmentHint) pointer() C.CommitmentHintPtr {
	return c.p
}

func finalizeCommitmentHint(c *commitmentHint) {
	C.ergo_lib_commitment_hint_delete(c.p)
}

// HintsBag is a collection of CommitmentHint to be used by a prover
type HintsBag interface {
	// Add adds CommitmentHint to the bag
	Add(hint CommitmentHint)
	// Len returns the length of the HintsBag
	Len() uint32
	// Get returns the CommitmentHint at the provided index if it exists
	Get(index uint32) (CommitmentHint, error)
	pointer() C.HintsBagPtr
}

type hintsBag struct {
	p C.HintsBagPtr
}

func newHintsBag(h *hintsBag) HintsBag {
	runtime.SetFinalizer(h, finalizeHintsBag)
	return h
}

// NewHintsBag creates an empty HintsBag
func NewHintsBag() HintsBag {
	var p C.HintsBagPtr
	C.ergo_lib_hints_bag_empty(&p)

	h := &hintsBag{p: p}
	return newHintsBag(h)
}

func (h *hintsBag) Add(hint CommitmentHint) {
	C.ergo_lib_hints_bag_add_commitment(h.p, hint.pointer())
}

func (h *hintsBag) Len() uint32 {
	res := C.ergo_lib_hints_bag_len(h.p)
	return uint32(res)
}

func (h *hintsBag) Get(index uint32) (CommitmentHint, error) {
	var p C.CommitmentHintPtr

	res := C.ergo_lib_hints_bag_get(h.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		c := &commitmentHint{p: p}
		return newCommitmentHint(c), nil
	}

	return nil, nil
}

func (h *hintsBag) pointer() C.HintsBagPtr {
	return h.p
}

func finalizeHintsBag(h *hintsBag) {
	C.ergo_lib_hints_bag_delete(h.p)
}

type TransactionHintsBag interface {
	// AddHintsForInput adds hints for input
	AddHintsForInput(index uint32, hintsBag HintsBag)
	// AllHintsForInput gets HintsBag corresponding to input index
	AllHintsForInput(index uint32) HintsBag
	pointer() C.TransactionHintsBagPtr
}

type transactionHintsBag struct {
	p C.TransactionHintsBagPtr
}

func newTransactionHintsBag(t *transactionHintsBag) TransactionHintsBag {
	runtime.SetFinalizer(t, finalizeTransactionHintsBag)
	return t
}

// NewTransactionHintsBag creates empty TransactionHintsBag
func NewTransactionHintsBag() TransactionHintsBag {
	var p C.TransactionHintsBagPtr
	C.ergo_lib_transaction_hints_bag_empty(&p)

	t := &transactionHintsBag{p: p}

	return newTransactionHintsBag(t)
}

func (t *transactionHintsBag) AddHintsForInput(index uint32, hintsBag HintsBag) {
	C.ergo_lib_transaction_hints_bag_add_hints_for_input(t.p, C.uintptr_t(index), hintsBag.pointer())
}

func (t *transactionHintsBag) AllHintsForInput(index uint32) HintsBag {
	var p C.HintsBagPtr
	C.ergo_lib_transaction_hints_bag_all_hints_for_input(t.p, C.uintptr_t(index), &p)
	h := &hintsBag{p: p}
	return newHintsBag(h)
}

func (t *transactionHintsBag) pointer() C.TransactionHintsBagPtr {
	return t.p
}

func finalizeTransactionHintsBag(t *transactionHintsBag) {
	C.ergo_lib_transaction_hints_bag_delete(t.p)
}

// ExtractHintsFromSignedTransaction extracts hints from signed transaction
func ExtractHintsFromSignedTransaction(
	transaction Transaction,
	stateContext StateContext,
	boxesToSpend Boxes,
	dataBoxes Boxes,
	realPropositions Propositions,
	simulatedPropositions Propositions) (TransactionHintsBag, error) {
	var p C.TransactionHintsBagPtr

	errPtr := C.ergo_lib_transaction_extract_hints(
		transaction.pointer(),
		stateContext.pointer(),
		boxesToSpend.pointer(),
		dataBoxes.pointer(),
		realPropositions.pointer(),
		simulatedPropositions.pointer(),
		&p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	th := &transactionHintsBag{p: p}

	return newTransactionHintsBag(th), nil
}

// UnsignedTransaction represents an unsigned transaction (inputs without proofs)
type UnsignedTransaction interface {
	// TxId returns TxId for this UnsignedTransaction
	TxId() TxId
	// UnsignedInputs returns UnsignedInputs for this UnsignedTransaction
	UnsignedInputs() UnsignedInputs
	// DataInputs returns DataInputs for this UnsignedTransaction
	DataInputs() DataInputs
	// OutputCandidates returns BoxCandidates for this UnsignedTransaction
	OutputCandidates() BoxCandidates
	// Json returns json representation of UnsignedTransaction as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
	Json() (string, error)
	// JsonEIP12 returns json representation of UnsignedTransaction as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
	JsonEIP12() (string, error)
	pointer() C.UnsignedTransactionPtr
}

type unsignedTransaction struct {
	p C.UnsignedTransactionPtr
}

func newUnsignedTransaction(u *unsignedTransaction) UnsignedTransaction {
	runtime.SetFinalizer(u, finalizeUnsignedTransaction)
	return u
}

// NewUnsignedTransactionFromJson parse UnsignedTransaction from JSON. Supports Ergo Node/Explorer API and box values and token amount encoded as strings.
func NewUnsignedTransactionFromJson(json string) (UnsignedTransaction, error) {
	unsTxJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(unsTxJsonStr))

	var p C.UnsignedTransactionPtr

	errPtr := C.ergo_lib_unsigned_tx_from_json(unsTxJsonStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ut := &unsignedTransaction{p: p}
	return newUnsignedTransaction(ut), nil
}

func (u *unsignedTransaction) TxId() TxId {
	var p C.TxIdPtr
	C.ergo_lib_unsigned_tx_id(u.p, &p)
	ti := &txId{p: p}
	return newTxId(ti)
}

func (u *unsignedTransaction) UnsignedInputs() UnsignedInputs {
	var p C.UnsignedInputsPtr
	C.ergo_lib_unsigned_tx_inputs(u.p, &p)
	ui := &unsignedInputs{p: p}
	return newUnsignedInputs(ui)
}

func (u *unsignedTransaction) DataInputs() DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_unsigned_tx_data_inputs(u.p, &p)
	di := &dataInputs{p: p}
	return newDataInputs(di)
}

func (u *unsignedTransaction) OutputCandidates() BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_unsigned_tx_output_candidates(u.p, &p)
	bc := &boxCandidates{p: p}
	return newBoxCandidates(bc)
}

func (u *unsignedTransaction) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_unsigned_tx_to_json(u.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (u *unsignedTransaction) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_unsigned_tx_to_json_eip12(u.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (u *unsignedTransaction) pointer() C.UnsignedTransactionPtr {
	return u.p
}

func finalizeUnsignedTransaction(u *unsignedTransaction) {
	C.ergo_lib_unsigned_tx_delete(u.p)
}

// Transaction is an atomic state transition operation. It destroys Boxes from the state
// and creates new ones. If transaction is spending boxes protected by some non-trivial scripts,
// its inputs should also contain proof of spending correctness - context extension (user-defined
// key-value map) and data inputs (links to existing boxes in the state) that may be used during
// script reduction to crypto, signatures that satisfies the remaining cryptographic protection
// of the script.
// Transactions are not encrypted, so it is possible to browse and view every transaction ever
// collected into a block.
type Transaction interface {
	// TxId returns TxId for this Transaction
	TxId() TxId
	// Inputs returns Inputs for this Transaction
	Inputs() Inputs
	// DataInputs returns DataInputs for this Transaction
	DataInputs() DataInputs
	// OutputCandidates returns BoxCandidates for this Transaction
	OutputCandidates() BoxCandidates
	// Outputs returns Boxes for this Transaction
	Outputs() Boxes
	// Json returns json representation of Transaction as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
	Json() (string, error)
	// JsonEIP12 returns json representation of Transaction as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
	JsonEIP12() (string, error)
	// Validate validates the current Transaction
	Validate(stateContext StateContext, boxesToSpent Boxes, dataBoxes Boxes) error
	pointer() C.TransactionPtr
}

type transaction struct {
	p C.TransactionPtr
}

func newTransaction(t *transaction) Transaction {
	runtime.SetFinalizer(t, finalizeTransaction)
	return t
}

// NewTransaction creates Transaction from UnsignedTransaction and an array of proofs in the same order as
// UnsignedTransaction inputs with empty proof indicated with empty ByteArray
func NewTransaction(unsignedTx UnsignedTransaction, proofs ByteArrays) (Transaction, error) {
	var p C.TransactionPtr

	errPtr := C.ergo_lib_tx_from_unsigned_tx(unsignedTx.pointer(), proofs.pointer(), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &transaction{p: p}
	return newTransaction(t), nil
}

// NewTransactionFromJson parse Transaction from JSON. Supports Ergo Node/Explorer API and box values and token amount encoded as strings.
func NewTransactionFromJson(json string) (Transaction, error) {
	txJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(txJsonStr))

	var p C.TransactionPtr

	errPtr := C.ergo_lib_tx_from_json(txJsonStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &transaction{p: p}
	return newTransaction(t), nil
}

func (t *transaction) TxId() TxId {
	var p C.TxIdPtr
	C.ergo_lib_tx_id(t.p, &p)
	ti := &txId{p: p}
	return newTxId(ti)
}

func (t *transaction) Inputs() Inputs {
	var p C.InputsPtr
	C.ergo_lib_tx_inputs(t.p, &p)
	i := &inputs{p: p}
	return newInputs(i)
}

func (t *transaction) DataInputs() DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_tx_data_inputs(t.p, &p)
	di := &dataInputs{p: p}
	return newDataInputs(di)
}

func (t *transaction) OutputCandidates() BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_tx_output_candidates(t.p, &p)
	bc := &boxCandidates{p: p}
	return newBoxCandidates(bc)
}

func (t *transaction) Outputs() Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_tx_outputs(t.p, &p)
	b := &boxes{p: p}
	return newBoxes(b)
}

func (t *transaction) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_tx_to_json(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (t *transaction) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_tx_to_json_eip12(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (t *transaction) Validate(stateContext StateContext, boxesToSpent Boxes, dataBoxes Boxes) error {
	errPtr := C.ergo_lib_tx_validate(t.p, stateContext.pointer(), boxesToSpent.pointer(), dataBoxes.pointer())
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (t *transaction) pointer() C.TransactionPtr {
	return t.p
}

func finalizeTransaction(t *transaction) {
	C.ergo_lib_tx_delete(t.p)
}
