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

// UnsignedInput used in construction of UnsignedTransactions
type UnsignedInput struct {
	p C.UnsignedInputPtr
}

func newUnsignedInput(u *UnsignedInput) *UnsignedInput {
	runtime.AddCleanup(u, finalizeUnsignedInput, u.p)
	return u
}

// BoxId returns the BoxId of the UnsignedInput
func (u *UnsignedInput) BoxId() *BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_unsigned_input_box_id(u.p, &p)

	bi := &BoxId{p: p}

	return newBoxId(bi)
}

// ContextExtension returns the ContextExtension of the UnsignedInput
func (u *UnsignedInput) ContextExtension() *ContextExtension {
	var p C.ContextExtensionPtr

	C.ergo_lib_unsigned_input_context_extension(u.p, &p)

	ce := &ContextExtension{p: p}

	return newContextExtension(ce)
}

func (u *UnsignedInput) pointer() C.UnsignedInputPtr {
	return u.p
}

func finalizeUnsignedInput(p C.UnsignedInputPtr) {
	C.ergo_lib_unsigned_input_delete(p)
}

// Input represents signed Inputs in signed Transaction
type Input struct {
	p C.InputPtr
}

func newInput(i *Input) *Input {
	runtime.AddCleanup(i, finalizeInput, i.p)
	return i
}

// BoxId returns BoxId of Input
func (i *Input) BoxId() *BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_input_box_id(i.p, &p)

	bi := &BoxId{p: p}

	return newBoxId(bi)
}

// SpendingProof returns spending proof of Input as ProverResult
func (i *Input) SpendingProof() *ProverResult {
	var p C.ProverResultPtr

	C.ergo_lib_input_spending_proof(i.p, &p)

	pr := &ProverResult{p: p}

	return newProverResult(pr)
}

func (i *Input) pointer() C.InputPtr {
	return i.p
}

func finalizeInput(p C.InputPtr) {
	C.ergo_lib_input_delete(p)
}

// ProverResult represents proof of correctness of tx spending
type ProverResult struct {
	p C.ProverResultPtr
}

func newProverResult(pr *ProverResult) *ProverResult {
	runtime.AddCleanup(pr, finalizeProverResult, pr.p)
	return pr
}

// Bytes returns proof bytes
func (pr *ProverResult) Bytes() []byte {
	proofLength := C.ergo_lib_prover_result_proof_len(pr.p)

	output := C.malloc(C.uintptr_t(proofLength))
	defer C.free(unsafe.Pointer(output))

	C.ergo_lib_prover_result_proof(pr.p, (*C.uint8_t)(output))

	result := C.GoBytes(unsafe.Pointer(output), C.int(proofLength))

	return result
}

// ContextExtension returns ContextExtension of ProverResult
func (pr *ProverResult) ContextExtension() *ContextExtension {
	var p C.ContextExtensionPtr

	C.ergo_lib_prover_result_context_extension(pr.p, &p)

	ce := &ContextExtension{p: p}

	return newContextExtension(ce)
}

// Json representation as text (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
func (pr *ProverResult) Json() (string, error) {
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

func finalizeProverResult(p C.ProverResultPtr) {
	C.ergo_lib_prover_result_delete(p)
}

// UnsignedInputs an ordered collection of UnsignedInput
type UnsignedInputs struct {
	p C.UnsignedInputsPtr
}

func newUnsignedInputs(u *UnsignedInputs) *UnsignedInputs {
	runtime.AddCleanup(u, finalizeUnsignedInputs, u.p)
	return u
}

// NewUnsignedInputs creates an empty UnsignedInputs collection
func NewUnsignedInputs() *UnsignedInputs {
	var p C.UnsignedInputsPtr
	C.ergo_lib_unsigned_inputs_new(&p)

	u := &UnsignedInputs{p: p}

	return newUnsignedInputs(u)
}

// Len returns the length of the collection
func (u *UnsignedInputs) Len() int {
	res := C.ergo_lib_unsigned_inputs_len(u.p)
	return int(res)
}

// Get returns the UnsignedInput at the provided index if it exists
func (u *UnsignedInputs) Get(index int) (*UnsignedInput, error) {
	var p C.UnsignedInputPtr

	res := C.ergo_lib_unsigned_inputs_get(u.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ui := &UnsignedInput{p: p}
		return newUnsignedInput(ui), nil
	}

	return nil, nil
}

// Add adds provided UnsignedInput to the end of the collection
func (u *UnsignedInputs) Add(unsignedInput *UnsignedInput) {
	C.ergo_lib_unsigned_inputs_add(unsignedInput.pointer(), u.p)
}

// All returns an iterator over all UnsignedInput inside the collection
func (u *UnsignedInputs) All() iter.Seq2[int, *UnsignedInput] {
	return func(yield func(int, *UnsignedInput) bool) {
		for i := 0; i < u.Len(); i++ {
			tk, err := u.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
	}
}

func finalizeUnsignedInputs(p C.UnsignedInputsPtr) {
	C.ergo_lib_unsigned_inputs_delete(p)
}

// Inputs an ordered collection of Input
type Inputs struct {
	p C.InputsPtr
}

func newInputs(i *Inputs) *Inputs {
	runtime.AddCleanup(i, finalizeInputs, i.p)
	return i
}

// NewInputs creates an empty Inputs collection
func NewInputs() *Inputs {
	var p C.InputsPtr
	C.ergo_lib_inputs_new(&p)

	i := &Inputs{p: p}

	return newInputs(i)
}

// Len returns the length of the collection
func (i *Inputs) Len() int {
	res := C.ergo_lib_inputs_len(i.p)
	return int(res)
}

// Get returns the Input at the provided index if it exists
func (i *Inputs) Get(index int) (*Input, error) {
	var p C.InputPtr

	res := C.ergo_lib_inputs_get(i.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		in := &Input{p: p}
		return newInput(in), nil
	}

	return nil, nil
}

// Add adds provided Input to the end of the collection
func (i *Inputs) Add(input *Input) {
	C.ergo_lib_inputs_add(input.pointer(), i.p)
}

// All returns an iterator over all Input inside the collection
func (i *Inputs) All() iter.Seq2[int, *Input] {
	return func(yield func(int, *Input) bool) {
		for j := 0; j < i.Len(); j++ {
			tk, err := i.Get(j)
			if err != nil {
				return
			}
			if !yield(j, tk) {
				return
			}
		}
	}
}

func finalizeInputs(p C.InputsPtr) {
	C.ergo_lib_inputs_delete(p)
}
