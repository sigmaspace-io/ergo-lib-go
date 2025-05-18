package ergo

/*
#include "ergo.h"
*/
import "C"
import "runtime"

// TxBuilder builds UnsignedTransaction
type TxBuilder struct {
	p C.TxBuilderPtr
}

func newTxBuilder(t *TxBuilder) *TxBuilder {
	runtime.AddCleanup(t, finalizeTxBuilder, t.p)
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
	boxSelection *BoxSelection,
	outputCandidates *BoxCandidates,
	currentHeight uint32,
	feeAmount *BoxValue,
	changeAddress *Address) *TxBuilder {
	var p C.TxBuilderPtr
	C.ergo_lib_tx_builder_new(
		boxSelection.pointer(),
		outputCandidates.pointer(),
		C.uint32_t(currentHeight),
		feeAmount.pointer(),
		changeAddress.pointer(),
		&p)
	runtime.KeepAlive(boxSelection)
	runtime.KeepAlive(outputCandidates)
	runtime.KeepAlive(feeAmount)
	runtime.KeepAlive(changeAddress)
	tb := &TxBuilder{p: p}
	return newTxBuilder(tb)
}

// SetDataInputs set data Inputs for Transaction
func (t *TxBuilder) SetDataInputs(dataInputs *DataInputs) {
	C.ergo_lib_tx_builder_set_data_inputs(t.p, dataInputs.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(dataInputs)
}

// SetContextExtension sets context extension for a given Input
func (t *TxBuilder) SetContextExtension(boxId *BoxId, contextExtension *ContextExtension) {
	C.ergo_lib_tx_builder_set_context_extension(t.p, boxId.pointer(), contextExtension.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(boxId)
	runtime.KeepAlive(contextExtension)
}

// SetTokenBurnPermit permits the burn of the given Token amount, i.e. allows this Token amount to be omitted in the outputs
func (t *TxBuilder) SetTokenBurnPermit(tokens *Tokens) {
	C.ergo_lib_tx_builder_set_token_burn_permit(t.p, tokens.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(tokens)
}

// DataInputs returns DataInputs of the TxBuilder
func (t *TxBuilder) DataInputs() *DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_tx_builder_data_inputs(t.p, &p)
	runtime.KeepAlive(t)
	di := &DataInputs{p: p}
	return newDataInputs(di)
}

// BoxSelection returns BoxSelection of the TxBuilder
func (t *TxBuilder) BoxSelection() *BoxSelection {
	var p C.BoxSelectionPtr
	C.ergo_lib_tx_builder_box_selection(t.p, &p)
	runtime.KeepAlive(t)
	bs := &BoxSelection{p: p}
	return newBoxSelection(bs)
}

// OutputCandidates returns BoxCandidates of the TxBuilder
func (t *TxBuilder) OutputCandidates() *BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_tx_builder_output_candidates(t.p, &p)
	runtime.KeepAlive(t)
	bc := &BoxCandidates{p: p}
	return newBoxCandidates(bc)
}

// CurrentHeight returns the current height
func (t *TxBuilder) CurrentHeight() uint32 {
	res := C.ergo_lib_tx_builder_current_height(t.p)
	runtime.KeepAlive(t)
	return uint32(res)
}

// FeeAmount returns the fee amount of the TxBuilder
func (t *TxBuilder) FeeAmount() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_tx_builder_fee_amount(t.p, &p)
	runtime.KeepAlive(t)
	bv := &BoxValue{p: p}
	return newBoxValue(bv)
}

// ChangeAddress returns the change Address of the TxBuilder
func (t *TxBuilder) ChangeAddress() *Address {
	var p C.AddressPtr
	C.ergo_lib_tx_builder_change_address(t.p, &p)
	runtime.KeepAlive(t)
	a := &Address{p: p}
	return newAddress(a)
}

// Build builds the UnsignedTransaction
func (t *TxBuilder) Build() (*UnsignedTransaction, error) {
	var p C.UnsignedTransactionPtr

	errPtr := C.ergo_lib_tx_builder_build(t.p, &p)
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ut := &UnsignedTransaction{p: p}
	return newUnsignedTransaction(ut), nil
}

func finalizeTxBuilder(p C.TxBuilderPtr) {
	C.ergo_lib_tx_builder_delete(p)
}

// SuggestedTxFee returns the suggested Transaction fee (semi-default value used across wallets and dApp as of Oct 2020)
func SuggestedTxFee() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_tx_builder_suggested_tx_fee(&p)
	bv := &BoxValue{p: p}
	return newBoxValue(bv)
}
