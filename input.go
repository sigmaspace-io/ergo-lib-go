package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// UnsignedInput used in construction of UnsignedTransactions
type UnsignedInput interface {
	// BoxId returns the BoxId of the UnsignedInput
	BoxId() BoxId
	// ContextExtension returns the ContextExtension of the UnsignedInput
	ContextExtension() ContextExtension
	pointer() C.UnsignedInputPtr
}

type unsignedInput struct {
	p C.UnsignedInputPtr
}

func newUnsignedInput(u *unsignedInput) UnsignedInput {
	runtime.SetFinalizer(u, finalizeUnsignedInput)
	return u
}

func (u *unsignedInput) BoxId() BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_unsigned_input_box_id(u.p, &p)

	bi := &boxId{p: p}

	return newBoxId(bi)
}

func (u *unsignedInput) ContextExtension() ContextExtension {
	var p C.ContextExtensionPtr

	C.ergo_lib_unsigned_input_context_extension(u.p, &p)

	ce := &contextExtension{p: p}

	return newContextExtension(ce)
}

func (u *unsignedInput) pointer() C.UnsignedInputPtr {
	return u.p
}

func finalizeUnsignedInput(u *unsignedInput) {
	C.ergo_lib_unsigned_input_delete(u.p)
}

// Input represents signed inputs in signed transaction
type Input interface {
	// BoxId returns BoxId of Input
	BoxId() BoxId
	// SpendingProof returns spending proof of Input as ProverResult
	SpendingProof() ProverResult
	pointer() C.InputPtr
}

type input struct {
	p C.InputPtr
}

func newInput(i *input) Input {
	runtime.SetFinalizer(i, finalizeInput)
	return i
}

func (i *input) BoxId() BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_input_box_id(i.p, &p)

	bi := &boxId{p: p}

	return newBoxId(bi)
}

func (i *input) SpendingProof() ProverResult {
	var p C.ProverResultPtr

	C.ergo_lib_input_spending_proof(i.p, &p)

	pr := &proverResult{p: p}

	return newProverResult(pr)
}

func (i *input) pointer() C.InputPtr {
	return i.p
}

func finalizeInput(i *input) {
	C.ergo_lib_input_delete(i.p)
}

// ProverResult represents proof of correctness of tx spending
type ProverResult interface {
	// Bytes returns proof bytes
	Bytes() []byte
	// ContextExtension returns ContextExtension of ProverResult
	ContextExtension() ContextExtension
	// Json representation as text (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
	Json() (string, error)
}

type proverResult struct {
	p C.ProverResultPtr
}

func newProverResult(pr *proverResult) ProverResult {
	runtime.SetFinalizer(pr, finalizeProverResult)
	return pr
}

func (pr *proverResult) Bytes() []byte {
	proofLength := C.ergo_lib_prover_result_proof_len(pr.p)

	output := C.malloc(C.uintptr_t(proofLength))
	defer C.free(unsafe.Pointer(output))

	C.ergo_lib_prover_result_proof(pr.p, (*C.uint8_t)(output))

	result := C.GoBytes(unsafe.Pointer(output), C.int(proofLength))

	return result
}

func (pr *proverResult) ContextExtension() ContextExtension {
	var p C.ContextExtensionPtr

	C.ergo_lib_prover_result_context_extension(pr.p, &p)

	ce := &contextExtension{p: p}

	return newContextExtension(ce)
}

func (pr *proverResult) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_prover_result_to_json(pr.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func finalizeProverResult(pr *proverResult) {
	C.ergo_lib_prover_result_delete(pr.p)
}

// UnsignedInputs an ordered collection of UnsignedInput
type UnsignedInputs interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the UnsignedInput at the provided index if it exists
	Get(index uint32) (UnsignedInput, error)
	// Add adds provided UnsignedInput to the end of the collection
	Add(unsignedInput UnsignedInput)
}

type unsignedInputs struct {
	p C.UnsignedInputsPtr
}

func newUnsignedInputs(u *unsignedInputs) UnsignedInputs {
	runtime.SetFinalizer(u, finalizeUnsignedInputs)
	return u
}

// NewUnsignedInputs creates an empty UnsignedInputs collection
func NewUnsignedInputs() UnsignedInputs {
	var p C.UnsignedInputsPtr
	C.ergo_lib_unsigned_inputs_new(&p)

	u := &unsignedInputs{p: p}

	return newUnsignedInputs(u)
}

func (u *unsignedInputs) Len() uint32 {
	res := C.ergo_lib_unsigned_inputs_len(u.p)
	return uint32(res)
}

func (u *unsignedInputs) Get(index uint32) (UnsignedInput, error) {
	var p C.UnsignedInputPtr

	res := C.ergo_lib_unsigned_inputs_get(u.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ui := &unsignedInput{p: p}
		return newUnsignedInput(ui), nil
	}

	return nil, nil
}

func (u *unsignedInputs) Add(unsignedInput UnsignedInput) {
	C.ergo_lib_unsigned_inputs_add(unsignedInput.pointer(), u.p)
}

func finalizeUnsignedInputs(u *unsignedInputs) {
	C.ergo_lib_unsigned_inputs_delete(u.p)
}

// Inputs an ordered collection of Input
type Inputs interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the Input at the provided index if it exists
	Get(index uint32) (Input, error)
	// Add adds provided Input to the end of the collection
	Add(input Input)
}

type inputs struct {
	p C.InputsPtr
}

func newInputs(i *inputs) Inputs {
	runtime.SetFinalizer(i, finalizeInputs)
	return i
}

// NewInputs creates an empty Inputs collection
func NewInputs() Inputs {
	var p C.InputsPtr
	C.ergo_lib_inputs_new(&p)

	i := &inputs{p: p}

	return newInputs(i)
}

func (i *inputs) Len() uint32 {
	res := C.ergo_lib_inputs_len(i.p)
	return uint32(res)
}

func (i *inputs) Get(index uint32) (Input, error) {
	var p C.InputPtr

	res := C.ergo_lib_inputs_get(i.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		in := &input{p: p}
		return newInput(in), nil
	}

	return nil, nil
}

func (i *inputs) Add(input Input) {
	C.ergo_lib_inputs_add(input.pointer(), i.p)
}

func finalizeInputs(i *inputs) {
	C.ergo_lib_inputs_delete(i.p)
}
