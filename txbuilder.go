package ergo

/*
#include "ergo.h"
*/
import "C"
import "runtime"

// TxBuilder builds UnsignedTransaction
type TxBuilder interface {
	// SetDataInputs set data inputs for transaction
	SetDataInputs(dataInputs DataInputs)
	// SetContextExtension sets context extension for a given input
	SetContextExtension(boxId BoxId, contextExtension ContextExtension)
	// SetTokenBurnPermit permits the burn of the given token amount, i.e. allows this token amount to be omitted in the outputs
	SetTokenBurnPermit(tokens Tokens)
	// DataInputs returns DataInputs of the TxBuilder
	DataInputs() DataInputs
	// BoxSelection returns BoxSelection of the TxBuilder
	BoxSelection() BoxSelection
	// OutputCandidates returns BoxCandidates of the TxBuilder
	OutputCandidates() BoxCandidates
	// CurrentHeight returns the current height
	CurrentHeight() uint32
	// FeeAmount returns the fee amount of the TxBuilder
	FeeAmount() BoxValue
	// ChangeAddress returns the change address of the TxBuilder
	ChangeAddress() Address
	// Build builds the UnsignedTransaction
	Build() (UnsignedTransaction, error)
}

type txBuilder struct {
	p C.TxBuilderPtr
}

func newTxBuilder(t *txBuilder) TxBuilder {
	runtime.SetFinalizer(t, finalizeTxBuilder)
	return t
}

// NewTxBuilder creates a new TxBuilder
// Parameters
// boxSelection - selected input boxes
// outputCandidates - output boxes to be "created" in this transaction
// currentHeight - chain height that will be used in additionally created boxes (change, miner's fee, etc.)
// feeAmount - miner's fee
// changeAddress - change (inputs - outputs) will be sent to this address
func NewTxBuilder(
	boxSelection BoxSelection,
	outputCandidates BoxCandidates,
	currentHeight uint32,
	feeAmount BoxValue,
	changeAddress Address) TxBuilder {
	var p C.TxBuilderPtr
	C.ergo_lib_tx_builder_new(
		boxSelection.pointer(),
		outputCandidates.pointer(),
		C.uint32_t(currentHeight),
		feeAmount.pointer(),
		changeAddress.pointer(),
		&p)
	tb := &txBuilder{p: p}
	return newTxBuilder(tb)
}

func (t *txBuilder) SetDataInputs(dataInputs DataInputs) {
	C.ergo_lib_tx_builder_set_data_inputs(t.p, dataInputs.pointer())
}

func (t *txBuilder) SetContextExtension(boxId BoxId, contextExtension ContextExtension) {
	C.ergo_lib_tx_builder_set_context_extension(t.p, boxId.pointer(), contextExtension.pointer())
}

func (t *txBuilder) SetTokenBurnPermit(tokens Tokens) {
	C.ergo_lib_tx_builder_set_token_burn_permit(t.p, tokens.pointer())
}

func (t *txBuilder) DataInputs() DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_tx_builder_data_inputs(t.p, &p)
	di := &dataInputs{p: p}
	return newDataInputs(di)
}

func (t *txBuilder) BoxSelection() BoxSelection {
	var p C.BoxSelectionPtr
	C.ergo_lib_tx_builder_box_selection(t.p, &p)
	bs := &boxSelection{p: p}
	return newBoxSelection(bs)
}

func (t *txBuilder) OutputCandidates() BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_tx_builder_output_candidates(t.p, &p)
	bc := &boxCandidates{p: p}
	return newBoxCandidates(bc)
}

func (t *txBuilder) CurrentHeight() uint32 {
	res := C.ergo_lib_tx_builder_current_height(t.p)
	return uint32(res)
}

func (t *txBuilder) FeeAmount() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_tx_builder_fee_amount(t.p, &p)
	bv := &boxValue{p: p}
	return newBoxValue(bv)
}

func (t *txBuilder) ChangeAddress() Address {
	var p C.AddressPtr
	C.ergo_lib_tx_builder_change_address(t.p, &p)
	a := &address{p: p}
	return newAddress(a)
}

func (t *txBuilder) Build() (UnsignedTransaction, error) {
	var p C.UnsignedTransactionPtr

	errPtr := C.ergo_lib_tx_builder_build(t.p, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ut := &unsignedTransaction{p: p}
	return newUnsignedTransaction(ut), nil
}

func finalizeTxBuilder(t *txBuilder) {
	C.ergo_lib_tx_builder_delete(t.p)
}

// SuggestedTxFee returns the suggested transaction fee (semi-default value used across wallets and dApp as of Oct 2020)
func SuggestedTxFee() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_tx_builder_suggested_tx_fee(&p)
	bv := &boxValue{p: p}
	return newBoxValue(bv)
}
