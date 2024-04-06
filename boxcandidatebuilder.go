package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// BoxCandidateBuilder is a builder to build a BoxCandidate
type BoxCandidateBuilder interface {
	// SetMinBoxValuePerByte sets minimal value (per byte of the serialized box size)
	SetMinBoxValuePerByte(minBoxValuePerByte uint32)
	// MinBoxValuePerByte returns minimal value (per byte of the serialized box size)
	MinBoxValuePerByte() uint32
	// SetValue sets new box value
	SetValue(boxValue BoxValue)
	// Value returns box value
	Value() BoxValue
	// CalcBoxSizeBytes calculates serialized box size(in bytes)
	CalcBoxSizeBytes() (uint32, error)
	// CalcMinBoxValue calculates minimal box value for the current box serialized size(in bytes)
	CalcMinBoxValue() (BoxValue, error)
	// SetRegisterValue sets register with a given id (R4 - R9) to the given value
	SetRegisterValue(registerId nonMandatoryRegisterId, constant Constant)
	//RegisterValue returns register value for the given register id (R4-R9), or nil if the register is empty
	RegisterValue(registerId nonMandatoryRegisterId) (Constant, error)
	// DeleteRegisterValue deletes register value(make register empty) for the given register id (R4 - R9)
	DeleteRegisterValue(registerId nonMandatoryRegisterId)
	// MintToken mints token, as defined in https://github.com/ergoplatform/eips/blob/master/eip-0004.md
	// Parameters:
	// token - token id(box id of the first input box in transaction) and token amount
	// tokenName - token name (will be encoded in R4)
	// tokenDesc - token description (will be encoded in R5)
	// numDecimals - number of decimals (will be encoded in R6)
	MintToken(token Token, tokenName string, tokenDesc string, numDecimals uint32)
	// AddToken adds given token id and token amount
	AddToken(tokenId TokenId, tokenAmount TokenAmount)
	// Build builds the box candidate
	Build() (BoxCandidate, error)
}

type boxCandidateBuilder struct {
	p C.ErgoBoxCandidateBuilderPtr
}

func newBoxCandidateBuilder(b *boxCandidateBuilder) BoxCandidateBuilder {
	runtime.SetFinalizer(b, finalizeBoxCandidateBuilder)
	return b
}

// NewBoxCandidateBuilder creates a BoxCandidateBuilder with required box parameters.
// Parameters:
// boxValue - amount of money associated with the box
// contract - guard Contract which should be evaluated to true in order to open/spend this box
// creationHeight - height when a transaction containing the box is created.
// It should not exceed the height of the block, containing the transaction with this box.
func NewBoxCandidateBuilder(boxValue BoxValue, contract Contract, creationHeight uint32) BoxCandidateBuilder {
	var p C.ErgoBoxCandidateBuilderPtr

	C.ergo_lib_ergo_box_candidate_builder_new(boxValue.pointer(), contract.pointer(), C.uint(creationHeight), &p)

	bc := &boxCandidateBuilder{p: p}

	return newBoxCandidateBuilder(bc)
}

func (b *boxCandidateBuilder) SetMinBoxValuePerByte(minBoxValuePerByte uint32) {
	C.ergo_lib_ergo_box_candidate_builder_set_min_box_value_per_byte(b.p, C.uint(minBoxValuePerByte))
}

func (b *boxCandidateBuilder) MinBoxValuePerByte() uint32 {
	res := C.ergo_lib_ergo_box_candidate_builder_min_box_value_per_byte(b.p)
	return uint32(res)
}

func (b *boxCandidateBuilder) SetValue(boxValue BoxValue) {
	C.ergo_lib_ergo_box_candidate_builder_set_value(b.p, boxValue.pointer())
}

func (b *boxCandidateBuilder) Value() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_candidate_builder_value(b.p, &p)
	bv := &boxValue{p: p}
	return newBoxValue(bv)
}

func (b *boxCandidateBuilder) CalcBoxSizeBytes() (uint32, error) {
	res := C.ergo_lib_ergo_box_candidate_builder_calc_box_size_bytes(b.p)
	err := newError(res.error)
	if err.isError() {
		return 0, err.error()
	}
	return uint32(res.value), nil
}

func (b *boxCandidateBuilder) CalcMinBoxValue() (BoxValue, error) {
	var p C.BoxValuePtr
	errPtr := C.ergo_lib_ergo_box_candidate_calc_min_box_value(b.p, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bv := &boxValue{p: p}
	return newBoxValue(bv), nil
}

func (b *boxCandidateBuilder) SetRegisterValue(registerId nonMandatoryRegisterId, constant Constant) {
	C.ergo_lib_ergo_box_candidate_builder_set_register_value(b.p, C.uchar(registerId), constant.pointer())
}

func (b *boxCandidateBuilder) RegisterValue(registerId nonMandatoryRegisterId) (Constant, error) {
	var p C.ConstantPtr
	res := C.ergo_lib_ergo_box_candidate_builder_register_value(b.p, C.uchar(registerId), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		c := &constant{p: p}
		return newConstant(c), nil
	}
	return nil, nil
}

func (b *boxCandidateBuilder) DeleteRegisterValue(registerId nonMandatoryRegisterId) {
	C.ergo_lib_ergo_box_candidate_builder_delete_register_value(b.p, C.uchar(registerId))
}

func (b *boxCandidateBuilder) MintToken(token Token, tokenName string, tokenDesc string, numDecimals uint32) {
	tknNameStr := C.CString(tokenName)
	defer C.free(unsafe.Pointer(tknNameStr))

	tknDescStr := C.CString(tokenDesc)
	defer C.free(unsafe.Pointer(tknDescStr))

	C.ergo_lib_ergo_box_candidate_builder_mint_token(b.p, token.pointer(), tknNameStr, tknDescStr, C.ulong(numDecimals))
}

func (b *boxCandidateBuilder) AddToken(tokenId TokenId, tokenAmount TokenAmount) {
	C.ergo_lib_ergo_box_candidate_builder_add_token(b.p, tokenId.pointer(), tokenAmount.pointer())
}

func (b *boxCandidateBuilder) Build() (BoxCandidate, error) {
	var p C.ErgoBoxCandidatePtr

	errPtr := C.ergo_lib_ergo_box_candidate_builder_build(b.p, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	bc := &boxCandidate{p: p}
	return newBoxCandidate(bc), nil
}

func finalizeBoxCandidateBuilder(b *boxCandidateBuilder) {
	C.ergo_lib_ergo_box_candidate_builder_delete(b.p)
}
