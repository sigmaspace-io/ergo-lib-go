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

type nonMandatoryRegisterId uint8

const (
	// R4 id for the non-mandatory R4 register
	R4 nonMandatoryRegisterId = 4
	// R5 id for the non-mandatory R5 register
	R5 nonMandatoryRegisterId = 5
	// R6 id for the non-mandatory R6 register
	R6 nonMandatoryRegisterId = 6
	// R7 id for the non-mandatory R7 register
	R7 nonMandatoryRegisterId = 7
	// R8 id for the non-mandatory R8 register
	R8 nonMandatoryRegisterId = 8
	// R9 id for the non-mandatory R9 register
	R9 nonMandatoryRegisterId = 9
)

// BoxId (32-byte digest)
type BoxId struct {
	p C.BoxIdPtr
}

func newBoxId(b *BoxId) *BoxId {
	runtime.AddCleanup(b, finalizeBoxId, b.p)
	return b
}

// NewBoxId creates a new ergo BoxId from the supplied base16 string.
func NewBoxId(s string) (*BoxId, error) {
	boxIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(boxIdStr))

	var p C.BoxIdPtr

	errPtr := C.ergo_lib_box_id_from_str(boxIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &BoxId{p}

	return newBoxId(b), nil
}

// Base16 returns the BoxId as base16 encoded string
func (b *BoxId) Base16() string {
	var boxIdStr *C.char

	C.ergo_lib_box_id_to_str(b.p, &boxIdStr)
	defer C.ergo_lib_delete_string(boxIdStr)
	runtime.KeepAlive(b)

	return C.GoString(boxIdStr)
}

// Equals checks if provided BoxId is same
func (b *BoxId) Equals(boxId *BoxId) bool {
	res := C.ergo_lib_box_id_eq(b.p, boxId.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxId)
	return bool(res)
}

func (b *BoxId) pointer() C.BoxIdPtr {
	return b.p
}

func finalizeBoxId(p C.BoxIdPtr) {
	C.ergo_lib_box_id_delete(p)
}

// BoxValue in nanoERGs with bound checks
type BoxValue struct {
	p C.BoxValuePtr
}

func newBoxValue(b *BoxValue) *BoxValue {
	runtime.AddCleanup(b, finalizeBoxValue, b.p)
	return b
}

// NewBoxValue creates a BoxValue from int64
func NewBoxValue(value int64) (*BoxValue, error) {
	var p C.BoxValuePtr

	errPtr := C.ergo_lib_box_value_from_i64(C.int64_t(value), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &BoxValue{p: p}

	return newBoxValue(b), nil
}

// Int64 returns BoxValue value as int64
func (b *BoxValue) Int64() int64 {
	value := C.ergo_lib_box_value_as_i64(b.p)
	runtime.KeepAlive(b)
	return int64(value)
}

// Equals checks if provided BoxValue is same
func (b *BoxValue) Equals(boxValue *BoxValue) bool {
	res := C.ergo_lib_box_value_eq(b.p, boxValue.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxValue)
	return bool(res)
}

func (b *BoxValue) pointer() C.BoxValuePtr {
	return b.p
}

func finalizeBoxValue(p C.BoxValuePtr) {
	C.ergo_lib_box_value_delete(p)
}

// SafeUserMinBoxValue returns recommended (safe) minimal BoxValue to use in case Box size estimation is unavailable.
// Allows Box size upto 2777 bytes with current min Box value per byte of 360 nanoERGs
func SafeUserMinBoxValue() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_box_value_safe_user_min(&p)

	b := &BoxValue{p: p}

	return newBoxValue(b)
}

// UnitsPerErgo returns number of units inside one ERGO (i.e. one ERG using nano ERG representation)
func UnitsPerErgo() int64 {
	units := C.ergo_lib_box_value_units_per_ergo()
	return int64(units)
}

// SumOfBoxValues creates a new BoxValue which is the sum of the arguments, throwing error if value is out of bounds
func SumOfBoxValues(boxValue0 *BoxValue, boxValue1 *BoxValue) (*BoxValue, error) {
	var p C.BoxValuePtr
	errPtr := C.ergo_lib_box_value_sum_of(boxValue0.pointer(), boxValue1.pointer(), &p)
	runtime.KeepAlive(boxValue0)
	runtime.KeepAlive(boxValue1)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &BoxValue{p: p}

	return newBoxValue(b), nil
}

// BoxCandidate contains the same fields as Box except for Transaction id and index, that will be calculated
// after full Transaction formation. Use BoxCandidateBuilder to create an instance
type BoxCandidate struct {
	p C.ErgoBoxCandidatePtr
}

func newBoxCandidate(b *BoxCandidate) *BoxCandidate {
	runtime.AddCleanup(b, finalizeBoxCandidate, b.p)
	return b
}

// RegisterValue returns value (Constant) stored in the register or nil if the register is empty
func (b *BoxCandidate) RegisterValue(registerId nonMandatoryRegisterId) (*Constant, error) {
	var p C.ConstantPtr
	rId := C.uchar(registerId)

	res := C.ergo_lib_ergo_box_candidate_register_value(b.p, rId, &p)
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

// CreationHeight returns the creation height of the BoxCandidate
func (b *BoxCandidate) CreationHeight() uint32 {
	height := C.ergo_lib_ergo_box_candidate_creation_height(b.p)
	runtime.KeepAlive(b)
	return uint32(height)
}

// Tokens returns the ergo Tokens for the BoxCandidate
func (b *BoxCandidate) Tokens() *Tokens {
	var p C.TokensPtr

	C.ergo_lib_ergo_box_candidate_tokens(b.p, &p)
	runtime.KeepAlive(b)

	t := &Tokens{p: p}

	return newTokens(t)
}

// Tree returns the ergo Tree for the BoxCandidate
func (b *BoxCandidate) Tree() *Tree {
	var p C.ErgoTreePtr

	C.ergo_lib_ergo_box_candidate_ergo_tree(b.p, &p)
	runtime.KeepAlive(b)

	t := &Tree{p: p}

	return newTree(t)
}

// BoxValue returns the BoxValue of the BoxCandidate
func (b *BoxCandidate) BoxValue() *BoxValue {
	var p C.BoxValuePtr

	C.ergo_lib_ergo_box_candidate_box_value(b.p, &p)
	runtime.KeepAlive(b)

	bv := &BoxValue{p: p}

	return newBoxValue(bv)
}

// Equals checks if provided BoxCandidate is same
func (b *BoxCandidate) Equals(candidate *BoxCandidate) bool {
	res := C.ergo_lib_ergo_box_candidate_eq(b.p, candidate.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(candidate)
	return bool(res)
}

func (b *BoxCandidate) pointer() C.ErgoBoxCandidatePtr {
	return b.p
}

func finalizeBoxCandidate(p C.ErgoBoxCandidatePtr) {
	C.ergo_lib_ergo_box_candidate_delete(p)
}

// Box that is taking part in some Transaction on the chain Differs with BoxCandidate
// by added Transaction id and an index in the Input of that Transaction
type Box struct {
	p C.ErgoBoxPtr
}

func newBox(b *Box) *Box {
	runtime.AddCleanup(b, finalizeBox, b.p)
	return b
}

// NewBox creates a new Box from provided Parameters:
// boxValue - amount of money associated with the Box
// creationHeight - height when a Transaction containing the Box is created.
// Contract - guarding Contract(Contract), which should be evaluated to true in order to open(spend) this Box
// TxId - Transaction id in which this Box was "created" (participated in outputs)
// index - index (in outputs) in the Transaction
func NewBox(boxValue *BoxValue, creationHeight uint32, contract *Contract, txId *TxId, index uint16, tokens *Tokens) (*Box, error) {
	var p C.ErgoBoxPtr

	errPtr := C.ergo_lib_ergo_box_new(boxValue.pointer(), C.uint32_t(creationHeight), contract.pointer(), txId.pointer(), C.uint16_t(index), tokens.pointer(), &p)
	runtime.KeepAlive(boxValue)
	runtime.KeepAlive(contract)
	runtime.KeepAlive(txId)
	runtime.KeepAlive(tokens)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &Box{p: p}

	return newBox(b), nil
}

// NewBoxFromJson parse Box from JSON. Supports Ergo Node/Explorer API and Box values and Token amount encoded as strings.
func NewBoxFromJson(json string) (*Box, error) {
	boxJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(boxJsonStr))

	var p C.ErgoBoxPtr

	errPtr := C.ergo_lib_ergo_box_from_json(boxJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &Box{p: p}

	return newBox(b), nil
}

// BoxId returns the BoxId of the Box
func (b *Box) BoxId() *BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_ergo_box_id(b.p, &p)
	runtime.KeepAlive(b)

	bi := &BoxId{p: p}

	return newBoxId(bi)
}

// RegisterValue returns value (Constant) stored in the register or nil if the register is empty
func (b *Box) RegisterValue(registerId nonMandatoryRegisterId) (*Constant, error) {
	var p C.ConstantPtr
	rId := C.uchar(registerId)

	res := C.ergo_lib_ergo_box_register_value(b.p, rId, &p)
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

// CreationHeight returns the creation height of the Box
func (b *Box) CreationHeight() uint32 {
	height := C.ergo_lib_ergo_box_creation_height(b.p)
	runtime.KeepAlive(b)
	return uint32(height)
}

// Tokens returns the ergo Tokens for the Box
func (b *Box) Tokens() *Tokens {
	var p C.TokensPtr
	C.ergo_lib_ergo_box_tokens(b.p, &p)
	runtime.KeepAlive(b)

	t := &Tokens{p: p}

	return newTokens(t)
}

// Tree returns the ergo Tree for the Box
func (b *Box) Tree() *Tree {
	var p C.ErgoTreePtr
	C.ergo_lib_ergo_box_ergo_tree(b.p, &p)
	runtime.KeepAlive(b)

	t := &Tree{p: p}

	return newTree(t)
}

// BoxValue returns the BoxValue of the Box
func (b *Box) BoxValue() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_value(b.p, &p)
	runtime.KeepAlive(b)

	bv := &BoxValue{p: p}

	return newBoxValue(bv)
}

// Json returns json representation of Box as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
func (b *Box) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_box_to_json(b.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(b)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// JsonEIP12 returns json representation of Box as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
func (b *Box) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_box_to_json_eip12(b.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(b)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// Size calculates serialized Box size(in bytes)
func (b *Box) Size() uint64 {
	res := C.ergo_lib_ergo_box_bytes_size(b.p)
	runtime.KeepAlive(b)
	return uint64(res)
}

// Equals checks if provided Box is same
func (b *Box) Equals(box *Box) bool {
	res := C.ergo_lib_ergo_box_eq(b.p, box.pointer())
	runtime.KeepAlive(b)
	return bool(res)
}

func (b *Box) pointer() C.ErgoBoxPtr {
	return b.p
}

func finalizeBox(p C.ErgoBoxPtr) {
	C.ergo_lib_ergo_box_delete(p)
}

// BoxAssetsData is a pair of value and Tokens for a Box
type BoxAssetsData struct {
	p C.ErgoBoxAssetsDataPtr
}

func newBoxAssetsData(b *BoxAssetsData) *BoxAssetsData {
	runtime.AddCleanup(b, finalizeBoxAssetsData, b.p)
	return b
}

// NewBoxAssetsData creates a new BoxAssetsData from the supplied BoxValue and Tokens
func NewBoxAssetsData(boxValue BoxValue, tokens Tokens) *BoxAssetsData {
	var p C.ErgoBoxAssetsDataPtr
	C.ergo_lib_ergo_box_assets_data_new(boxValue.pointer(), tokens.pointer(), &p)
	runtime.KeepAlive(boxValue)
	runtime.KeepAlive(tokens)

	b := &BoxAssetsData{p: p}

	return newBoxAssetsData(b)
}

// BoxValue returns the BoxValue of the BoxAssetsData
func (b *BoxAssetsData) BoxValue() *BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_assets_data_value(b.p, &p)
	runtime.KeepAlive(b)

	bv := &BoxValue{p: p}

	return newBoxValue(bv)
}

// Tokens returns the Tokens of the BoxAssetsData
func (b *BoxAssetsData) Tokens() *Tokens {
	var p C.TokensPtr
	C.ergo_lib_ergo_box_assets_data_tokens(b.p, &p)
	runtime.KeepAlive(b)

	t := &Tokens{p: p}

	return newTokens(t)
}

// Equals checks if provided BoxAssetsData is same
func (b *BoxAssetsData) Equals(boxAssetsData *BoxAssetsData) bool {
	res := C.ergo_lib_ergo_box_assets_data_eq(b.p, boxAssetsData.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxAssetsData)
	return bool(res)
}

func (b *BoxAssetsData) pointer() C.ErgoBoxAssetsDataPtr {
	return b.p
}

func finalizeBoxAssetsData(p C.ErgoBoxAssetsDataPtr) {
	C.ergo_lib_ergo_box_assets_data_delete(p)
}

// BoxAssetsDataList is an ordered collection of BoxAssetsData
type BoxAssetsDataList struct {
	p C.ErgoBoxAssetsDataListPtr
}

func newBoxAssetsDataList(b *BoxAssetsDataList) *BoxAssetsDataList {
	runtime.AddCleanup(b, finalizeBoxAssetsDataList, b.p)
	return b
}

// NewBoxAssetsDataList creates an empty BoxAssetsDataList
func NewBoxAssetsDataList() *BoxAssetsDataList {
	var p C.ErgoBoxAssetsDataListPtr
	C.ergo_lib_ergo_box_assets_data_list_new(&p)

	b := &BoxAssetsDataList{p: p}

	return newBoxAssetsDataList(b)
}

// Len returns the length of the collection
func (b *BoxAssetsDataList) Len() int {
	res := C.ergo_lib_ergo_box_assets_data_list_len(b.p)
	runtime.KeepAlive(b)
	return int(res)
}

// Get returns the BoxAssetsData at the provided index if it exists
func (b *BoxAssetsDataList) Get(index int) (*BoxAssetsData, error) {
	var p C.ErgoBoxAssetsDataPtr

	res := C.ergo_lib_ergo_box_assets_data_list_get(b.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &BoxAssetsData{p: p}
		return newBoxAssetsData(ba), nil
	}

	return nil, nil
}

// Add adds provided BoxAssetsData to the end of the collection
func (b *BoxAssetsDataList) Add(boxAssetsData *BoxAssetsData) {
	C.ergo_lib_ergo_box_assets_data_list_add(boxAssetsData.pointer(), b.p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxAssetsData)
}

// All returns an iterator over all BoxAssetsData inside the collection
func (b *BoxAssetsDataList) All() iter.Seq2[int, *BoxAssetsData] {
	return func(yield func(int, *BoxAssetsData) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(b)
	}
}

func (b *BoxAssetsDataList) pointer() C.ErgoBoxAssetsDataListPtr {
	return b.p
}

func finalizeBoxAssetsDataList(p C.ErgoBoxAssetsDataListPtr) {
	C.ergo_lib_ergo_box_assets_data_list_delete(p)
}

// BoxCandidates is an ordered collection of BoxCandidate
type BoxCandidates struct {
	p C.ErgoBoxCandidatesPtr
}

func newBoxCandidates(b *BoxCandidates) *BoxCandidates {
	runtime.AddCleanup(b, finalizeBoxCandidates, b.p)
	return b
}

// NewBoxCandidates creates an empty BoxCandidates collection
func NewBoxCandidates() *BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_ergo_box_candidates_new(&p)

	b := &BoxCandidates{p: p}

	return newBoxCandidates(b)
}

// Len returns the length of the collection
func (b *BoxCandidates) Len() int {
	res := C.ergo_lib_ergo_box_candidates_len(b.p)
	runtime.KeepAlive(b)
	return int(res)
}

// Get returns the BoxCandidate at the provided index if it exists
func (b *BoxCandidates) Get(index int) (*BoxCandidate, error) {
	var p C.ErgoBoxCandidatePtr

	res := C.ergo_lib_ergo_box_candidates_get(b.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &BoxCandidate{p: p}
		return newBoxCandidate(ba), nil
	}

	return nil, nil
}

// Add adds provided BoxCandidate to the end of the collection
func (b *BoxCandidates) Add(boxCandidate *BoxCandidate) {
	C.ergo_lib_ergo_box_candidates_add(boxCandidate.pointer(), b.p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxCandidate)
}

// All returns an iterator over all BoxCandidate inside the collection
func (b *BoxCandidates) All() iter.Seq2[int, *BoxCandidate] {
	return func(yield func(int, *BoxCandidate) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(b)
	}
}

func (b *BoxCandidates) pointer() C.ErgoBoxCandidatesPtr {
	return b.p
}

func finalizeBoxCandidates(p C.ErgoBoxCandidatesPtr) {
	C.ergo_lib_ergo_box_candidates_delete(p)
}

// Boxes an ordered collection of Box
type Boxes struct {
	p C.ErgoBoxesPtr
}

func newBoxes(b *Boxes) *Boxes {
	runtime.AddCleanup(b, finalizeBoxes, b.p)
	return b
}

// NewBoxes creates an empty Boxes collection
func NewBoxes() *Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_ergo_boxes_new(&p)

	b := &Boxes{p: p}

	return newBoxes(b)
}

// Len returns the length of the collection
func (b *Boxes) Len() int {
	res := C.ergo_lib_ergo_boxes_len(b.p)
	runtime.KeepAlive(b)
	return int(res)
}

// Get returns the Box at the provided index if it exists
func (b *Boxes) Get(index int) (*Box, error) {
	var p C.ErgoBoxPtr

	res := C.ergo_lib_ergo_boxes_get(b.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(b)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &Box{p: p}
		return newBox(ba), nil
	}

	return nil, nil
}

// Add adds provided Box to the end of the collection
func (b *Boxes) Add(box *Box) {
	C.ergo_lib_ergo_boxes_add(box.pointer(), b.p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(box)
}

// All returns an iterator over all Box inside the collection
func (b *Boxes) All() iter.Seq2[int, *Box] {
	return func(yield func(int, *Box) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(b)
	}
}

func (b *Boxes) pointer() C.ErgoBoxesPtr {
	return b.p
}

func finalizeBoxes(p C.ErgoBoxesPtr) {
	C.ergo_lib_ergo_boxes_delete(p)
}
