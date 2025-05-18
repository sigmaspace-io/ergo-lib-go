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
type BoxCandidateBuilder struct {
	p C.ErgoBoxCandidateBuilderPtr
}

func newBoxCandidateBuilder(b *BoxCandidateBuilder) *BoxCandidateBuilder {
	runtime.AddCleanup(b, finalizeBoxCandidateBuilder, b.p)
	return b
}

// NewBoxCandidateBuilder creates a BoxCandidateBuilder with required Box Parameters.
// Parameters:
// boxValue - amount of money associated with the Box
// Contract - guard Contract which should be evaluated to true in order to open/spend this Box
// creationHeight - height when a Transaction containing the Box is created.
// It should not exceed the height of the block, containing the Transaction with this Box.
func NewBoxCandidateBuilder(boxValue *BoxValue, contract *Contract, creationHeight uint32) *BoxCandidateBuilder {
	var p C.ErgoBoxCandidateBuilderPtr

	C.ergo_lib_ergo_box_candidate_builder_new(boxValue.pointer(), contract.pointer(), C.uint32_t(creationHeight), &p)
	runtime.KeepAlive(boxValue)
	runtime.KeepAlive(contract)

	bc := &BoxCandidateBuilder{p: p}

	return newBoxCandidateBuilder(bc)
}

// SetMinBoxValuePerByte sets minimal value (per byte of the serialized Box size)
func (b *BoxCandidateBuilder) SetMinBoxValuePerByte(minBoxValuePerByte uint32) {
	C.ergo_lib_ergo_box_candidate_builder_set_min_box_value_per_byte(b.p, C.uint32_t(minBoxValuePerByte))
	runtime.KeepAlive(b)
}

// MinBoxValuePerByte returns minimal value (per byte of the serialized Box size)
func (b *BoxCandidateBuilder) MinBoxValuePerByte() uint32 {
	res := C.ergo_lib_ergo_box_candidate_builder_min_box_value_per_byte(b.p)
	runtime.KeepAlive(b)
	return uint32(res)
}

// SetValue sets new Box value
func (b *BoxCandidateBuilder) SetValue(boxValue *BoxValue) {
	C.ergo_lib_ergo_box_candidate_builder_set_value(b.p, boxValue.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxValue)
}

// Value returns Box value
func (b *BoxCandidateBuilder) Value() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_candidate_builder_value(b.p, &p)
	runtime.KeepAlive(b)
	bv := &BoxValue{p: p}
	return newBoxValue(bv)
}

// CalcBoxSizeBytes calculates serialized Box size(in bytes)
func (b *BoxCandidateBuilder) CalcBoxSizeBytes() (uint32, error) {
	res := C.ergo_lib_ergo_box_candidate_builder_calc_box_size_bytes(b.p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return 0, err.error()
	}
	return uint32(res.value), nil
}

// CalcMinBoxValue calculates minimal Box value for the current Box serialized size(in bytes)
func (b *BoxCandidateBuilder) CalcMinBoxValue() (*BoxValue, error) {
	var p C.BoxValuePtr
	errPtr := C.ergo_lib_ergo_box_candidate_calc_min_box_value(b.p, &p)
	runtime.KeepAlive(b)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bv := &BoxValue{p: p}
	return newBoxValue(bv), nil
}

// SetRegisterValue sets register with a given id (R4 - R9) to the given value
func (b *BoxCandidateBuilder) SetRegisterValue(registerId nonMandatoryRegisterId, constant *Constant) {
	C.ergo_lib_ergo_box_candidate_builder_set_register_value(b.p, C.uchar(registerId), constant.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(constant)
}

// RegisterValue returns register value for the given register id (R4 - R9), or nil if the register is empty
func (b *BoxCandidateBuilder) RegisterValue(registerId nonMandatoryRegisterId) (*Constant, error) {
	var p C.ConstantPtr
	res := C.ergo_lib_ergo_box_candidate_builder_register_value(b.p, C.uchar(registerId), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		c := &Constant{p: p}
		return newConstant(c), nil
	}
	return nil, nil
}

// DeleteRegisterValue deletes register value(make register empty) for the given register id (R4 - R9)
func (b *BoxCandidateBuilder) DeleteRegisterValue(registerId nonMandatoryRegisterId) {
	C.ergo_lib_ergo_box_candidate_builder_delete_register_value(b.p, C.uchar(registerId))
	runtime.KeepAlive(b)
}

// MintToken mints Token, as defined in https://github.com/ergoplatform/eips/blob/master/eip-0004.md
// Parameters:
// Token - Token id(Box id of the first Input Box in Transaction) and Token amount
// tokenName - Token name (will be encoded in R4)
// tokenDesc - Token description (will be encoded in R5)
// numDecimals - number of decimals (will be encoded in R6)
func (b *BoxCandidateBuilder) MintToken(token *Token, tokenName string, tokenDesc string, numDecimals uint32) {
	tknNameStr := C.CString(tokenName)
	defer C.free(unsafe.Pointer(tknNameStr))

	tknDescStr := C.CString(tokenDesc)
	defer C.free(unsafe.Pointer(tknDescStr))

	C.ergo_lib_ergo_box_candidate_builder_mint_token(b.p, token.pointer(), tknNameStr, tknDescStr, C.uintptr_t(numDecimals))
	runtime.KeepAlive(b)
	runtime.KeepAlive(token)
}

// AddToken adds given Token id and Token amount
func (b *BoxCandidateBuilder) AddToken(tokenId *TokenId, tokenAmount *TokenAmount) {
	C.ergo_lib_ergo_box_candidate_builder_add_token(b.p, tokenId.pointer(), tokenAmount.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(tokenId)
	runtime.KeepAlive(tokenAmount)
}

// Build builds the Box candidate
func (b *BoxCandidateBuilder) Build() (*BoxCandidate, error) {
	var p C.ErgoBoxCandidatePtr

	errPtr := C.ergo_lib_ergo_box_candidate_builder_build(b.p, &p)
	runtime.KeepAlive(b)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	bc := &BoxCandidate{p: p}
	return newBoxCandidate(bc), nil
}

func finalizeBoxCandidateBuilder(p C.ErgoBoxCandidateBuilderPtr) {
	C.ergo_lib_ergo_box_candidate_builder_delete(p)
}
