package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
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
type BoxId interface {
	// Base16 returns the BoxId as base16 encoded string
	Base16() string
	pointer() C.BoxIdPtr
}

type boxId struct {
	p C.BoxIdPtr
}

func newBoxId(b *boxId) BoxId {
	runtime.SetFinalizer(b, finalizeBoxId)
	return b
}

// NewBoxId creates a new ergo BoxId from the supplied base16 string.
func NewBoxId(s string) (BoxId, error) {
	boxIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(boxIdStr))

	var p C.BoxIdPtr

	errPtr := C.ergo_lib_box_id_from_str(boxIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	b := &boxId{p}

	return newBoxId(b), nil
}

func (b *boxId) Base16() string {
	var boxIdStr *C.char

	C.ergo_lib_box_id_to_str(b.p, &boxIdStr)
	defer C.ergo_lib_delete_string(boxIdStr)

	return C.GoString(boxIdStr)
}

func (b *boxId) pointer() C.BoxIdPtr {
	return b.p
}

func finalizeBoxId(b *boxId) {
	C.ergo_lib_box_id_delete(b.p)
}

// BoxValue in nanoERGs with bound checks
type BoxValue interface {
	// Int64 returns BoxValue value as int64
	Int64() int64
	pointer() C.BoxValuePtr
}

type boxValue struct {
	p C.BoxValuePtr
}

func newBoxValue(b *boxValue) BoxValue {
	runtime.SetFinalizer(b, finalizeBoxValue)
	return b
}

// NewBoxValue creates a BoxValue from int64
func NewBoxValue(value int64) (BoxValue, error) {
	var p C.BoxValuePtr

	errPtr := C.ergo_lib_box_value_from_i64(C.int64_t(value), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &boxValue{p: p}

	return newBoxValue(b), nil
}

func (b *boxValue) Int64() int64 {
	value := C.ergo_lib_box_value_as_i64(b.p)
	return int64(value)
}

func (b *boxValue) pointer() C.BoxValuePtr {
	return b.p
}

func finalizeBoxValue(b *boxValue) {
	C.ergo_lib_box_value_delete(b.p)
}

// SafeUserMinBoxValue returns recommended (safe) minimal BoxValue to use in case box size estimation is unavailable.
// Allows box size upto 2777 bytes with current min box value per byte of 360 nanoERGs
func SafeUserMinBoxValue() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_box_value_safe_user_min(&p)

	b := &boxValue{p: p}

	return newBoxValue(b)
}

// UnitsPerErgo returns number of units inside one ERGO (i.e. one ERG using nano ERG representation)
func UnitsPerErgo() int64 {
	units := C.ergo_lib_box_value_units_per_ergo()
	return int64(units)
}

// SumOfBoxValues creates a new BoxValue which is the sum of the arguments, throwing error if value is out of bounds
func SumOfBoxValues(boxValue0 BoxValue, boxValue1 BoxValue) (BoxValue, error) {
	var p C.BoxValuePtr
	errPtr := C.ergo_lib_box_value_sum_of(boxValue0.pointer(), boxValue1.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &boxValue{p: p}

	return newBoxValue(b), nil
}

// BoxCandidate contains the same fields as Box except for transaction id and index, that will be calculated
// after full transaction formation. Use BoxCandidateBuilder to create an instance
type BoxCandidate interface {
	// RegisterValue returns value (Constant) stored in the register or nil if the register is empty
	RegisterValue(registerId nonMandatoryRegisterId) (Constant, error)
	// CreationHeight returns the creation height of the BoxCandidate
	CreationHeight() uint32
	// Tokens returns the ergo Tokens for the BoxCandidate
	Tokens() Tokens
	// Tree returns the ergo Tree for the BoxCandidate
	Tree() Tree
	// BoxValue returns the BoxValue of the BoxCandidate
	BoxValue() BoxValue
	pointer() C.ErgoBoxCandidatePtr
}

type boxCandidate struct {
	p C.ErgoBoxCandidatePtr
}

func newBoxCandidate(b *boxCandidate) BoxCandidate {
	runtime.SetFinalizer(b, finalizeBoxCandidate)
	return b
}

func (b *boxCandidate) RegisterValue(registerId nonMandatoryRegisterId) (Constant, error) {
	var p C.ConstantPtr
	rId := C.uchar(registerId)

	res := C.ergo_lib_ergo_box_candidate_register_value(b.p, rId, &p)
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

func (b *boxCandidate) CreationHeight() uint32 {
	height := C.ergo_lib_ergo_box_candidate_creation_height(b.p)
	return uint32(height)
}

func (b *boxCandidate) Tokens() Tokens {
	var p C.TokensPtr

	C.ergo_lib_ergo_box_candidate_tokens(b.p, &p)

	t := &tokens{p: p}

	return newTokens(t)
}

func (b *boxCandidate) Tree() Tree {
	var p C.ErgoTreePtr

	C.ergo_lib_ergo_box_candidate_ergo_tree(b.p, &p)

	t := &tree{p: p}

	return newTree(t)
}

func (b *boxCandidate) BoxValue() BoxValue {
	var p C.BoxValuePtr

	C.ergo_lib_ergo_box_candidate_box_value(b.p, &p)

	bv := &boxValue{p: p}

	return newBoxValue(bv)
}

func (b *boxCandidate) pointer() C.ErgoBoxCandidatePtr {
	return b.p
}

func finalizeBoxCandidate(b *boxCandidate) {
	C.ergo_lib_ergo_box_candidate_delete(b.p)
}

// Box that is taking part in some transaction on the chain Differs with BoxCandidate
// by added transaction id and an index in the input of that transaction
type Box interface {
	// BoxId returns the BoxId of the Box
	BoxId() BoxId
	// RegisterValue returns value (Constant) stored in the register or nil if the register is empty
	RegisterValue(registerId nonMandatoryRegisterId) (Constant, error)
	// CreationHeight returns the creation height of the Box
	CreationHeight() uint32
	// Tokens returns the ergo Tokens for the Box
	Tokens() Tokens
	// Tree returns the ergo Tree for the Box
	Tree() Tree
	// BoxValue returns the BoxValue of the Box
	BoxValue() BoxValue
	// Json returns json representation of Box as string (compatible with Ergo Node/Explorer API, numbers are encoded as numbers)
	Json() (string, error)
	// JsonEIP12 returns json representation of Box as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
	JsonEIP12() (string, error)
	pointer() C.ErgoBoxPtr
}

type box struct {
	p C.ErgoBoxPtr
}

func newBox(b *box) Box {
	runtime.SetFinalizer(b, finalizeBox)
	return b
}

// NewBox creates a new Box from provided parameters:
// boxValue - amount of money associated with the box
// creationHeight - height when a transaction containing the box is created.
// contract - guarding contract(Contract), which should be evaluated to true in order to open(spend) this box
// txId - transaction id in which this box was "created" (participated in outputs)
// index - index (in outputs) in the transaction
func NewBox(boxValue BoxValue, creationHeight uint32, contract Contract, txId TxId, index uint16, tokens Tokens) (Box, error) {
	var p C.ErgoBoxPtr

	errPtr := C.ergo_lib_ergo_box_new(boxValue.pointer(), C.uint(creationHeight), contract.pointer(), txId.pointer(), C.ushort(index), tokens.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &box{p: p}

	return newBox(b), nil
}

// NewBoxFromJson parse Box from JSON. Supports Ergo Node/Explorer API and box values and token amount encoded as strings.
func NewBoxFromJson(json string) (Box, error) {
	boxJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(boxJsonStr))

	var p C.ErgoBoxPtr

	errPtr := C.ergo_lib_ergo_box_from_json(boxJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	b := &box{p: p}

	return newBox(b), nil
}

func (b *box) BoxId() BoxId {
	var p C.BoxIdPtr

	C.ergo_lib_ergo_box_id(b.p, &p)

	bi := &boxId{p: p}

	return newBoxId(bi)
}

func (b *box) RegisterValue(registerId nonMandatoryRegisterId) (Constant, error) {
	var p C.ConstantPtr
	rId := C.uchar(registerId)

	res := C.ergo_lib_ergo_box_register_value(b.p, rId, &p)
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

func (b *box) CreationHeight() uint32 {
	height := C.ergo_lib_ergo_box_creation_height(b.p)
	return uint32(height)
}

func (b *box) Tokens() Tokens {
	var p C.TokensPtr
	C.ergo_lib_ergo_box_tokens(b.p, &p)

	t := &tokens{p: p}

	return newTokens(t)
}

func (b *box) Tree() Tree {
	var p C.ErgoTreePtr
	C.ergo_lib_ergo_box_ergo_tree(b.p, &p)

	t := &tree{p: p}

	return newTree(t)
}

func (b *box) BoxValue() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_value(b.p, &p)

	bv := &boxValue{p: p}

	return newBoxValue(bv)
}

func (b *box) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_box_to_json(b.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (b *box) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_box_to_json_eip12(b.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (b *box) pointer() C.ErgoBoxPtr {
	return b.p
}

func finalizeBox(b *box) {
	C.ergo_lib_ergo_box_delete(b.p)
}

// BoxAssetsData is a pair of value and tokens for a box
type BoxAssetsData interface {
	// BoxValue returns the BoxValue of the BoxAssetsData
	BoxValue() BoxValue
	// Tokens returns the Tokens of the BoxAssetsData
	Tokens() Tokens
	pointer() C.ErgoBoxAssetsDataPtr
}

type boxAssetsData struct {
	p C.ErgoBoxAssetsDataPtr
}

func newBoxAssetsData(b *boxAssetsData) BoxAssetsData {
	runtime.SetFinalizer(b, finalizeBoxAssetsData)
	return b
}

// NewBoxAssetsData creates a new BoxAssetsData from the supplied BoxValue and Tokens
func NewBoxAssetsData(boxValue BoxValue, tokens Tokens) BoxAssetsData {
	var p C.ErgoBoxAssetsDataPtr
	C.ergo_lib_ergo_box_assets_data_new(boxValue.pointer(), tokens.pointer(), &p)

	b := &boxAssetsData{p: p}

	return newBoxAssetsData(b)
}

func (b *boxAssetsData) BoxValue() BoxValue {
	var p C.BoxValuePtr
	C.ergo_lib_ergo_box_assets_data_value(b.p, &p)

	bv := &boxValue{p: p}

	return newBoxValue(bv)
}

func (b *boxAssetsData) Tokens() Tokens {
	var p C.TokensPtr
	C.ergo_lib_ergo_box_assets_data_tokens(b.p, &p)

	t := &tokens{p: p}

	return newTokens(t)
}

func (b *boxAssetsData) pointer() C.ErgoBoxAssetsDataPtr {
	return b.p
}

func finalizeBoxAssetsData(b *boxAssetsData) {
	C.ergo_lib_ergo_box_assets_data_delete(b.p)
}

// BoxAssetsDataList is an ordered collection of BoxAssetsData
type BoxAssetsDataList interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the BoxAssetsData at the provided index if it exists
	Get(index uint32) (BoxAssetsData, error)
	// Add adds provided BoxAssetsData to the end of the collection
	Add(boxAssetsData BoxAssetsData)
}

type boxAssetsDataList struct {
	p C.ErgoBoxAssetsDataListPtr
}

func newBoxAssetsDataList(b *boxAssetsDataList) BoxAssetsDataList {
	runtime.SetFinalizer(b, finalizeBoxAssetsDataList)
	return b
}

// NewBoxAssetsDataList creates an empty BoxAssetsDataList
func NewBoxAssetsDataList() BoxAssetsDataList {
	var p C.ErgoBoxAssetsDataListPtr
	C.ergo_lib_ergo_box_assets_data_list_new(&p)

	b := &boxAssetsDataList{p: p}

	return newBoxAssetsDataList(b)
}

func (b *boxAssetsDataList) Len() uint32 {
	res := C.ergo_lib_ergo_box_assets_data_list_len(b.p)
	return uint32(res)
}

func (b *boxAssetsDataList) Get(index uint32) (BoxAssetsData, error) {
	var p C.ErgoBoxAssetsDataPtr

	res := C.ergo_lib_ergo_box_assets_data_list_get(b.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &boxAssetsData{p: p}
		return newBoxAssetsData(ba), nil
	}

	return nil, nil
}

func (b *boxAssetsDataList) Add(boxAssetsData BoxAssetsData) {
	C.ergo_lib_ergo_box_assets_data_list_add(boxAssetsData.pointer(), b.p)
}

func finalizeBoxAssetsDataList(b *boxAssetsDataList) {
	C.ergo_lib_ergo_box_assets_data_list_delete(b.p)
}

// BoxCandidates is an ordered collection of BoxCandidate
type BoxCandidates interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the BoxCandidate at the provided index if it exists
	Get(index uint32) (BoxCandidate, error)
	// Add adds provided BoxCandidate to the end of the collection
	Add(boxCandidate BoxCandidate)
}

type boxCandidates struct {
	p C.ErgoBoxCandidatesPtr
}

func newBoxCandidates(b *boxCandidates) BoxCandidates {
	runtime.SetFinalizer(b, finalizeBoxCandidates)
	return b
}

// NewBoxCandidates creates an empty BoxCandidates collection
func NewBoxCandidates() BoxCandidates {
	var p C.ErgoBoxCandidatesPtr
	C.ergo_lib_ergo_box_candidates_new(&p)

	b := &boxCandidates{p: p}

	return newBoxCandidates(b)
}

func (b *boxCandidates) Len() uint32 {
	res := C.ergo_lib_ergo_box_candidates_len(b.p)
	return uint32(res)
}

func (b *boxCandidates) Get(index uint32) (BoxCandidate, error) {
	var p C.ErgoBoxCandidatePtr

	res := C.ergo_lib_ergo_box_candidates_get(b.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &boxCandidate{p: p}
		return newBoxCandidate(ba), nil
	}

	return nil, nil
}

func (b *boxCandidates) Add(boxCandidate BoxCandidate) {
	C.ergo_lib_ergo_box_candidates_add(boxCandidate.pointer(), b.p)
}

func finalizeBoxCandidates(b *boxCandidates) {
	C.ergo_lib_ergo_box_candidates_delete(b.p)
}

// Boxes an ordered collection of Box
type Boxes interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the Box at the provided index if it exists
	Get(index uint32) (Box, error)
	// Add adds provided Box to the end of the collection
	Add(box Box)
	pointer() C.ErgoBoxesPtr
}

type boxes struct {
	p C.ErgoBoxesPtr
}

func newBoxes(b *boxes) Boxes {
	runtime.SetFinalizer(b, finalizeBoxes)
	return b
}

// NewBoxes creates an empty Boxes collection
func NewBoxes() Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_ergo_boxes_new(&p)

	b := &boxes{p: p}

	return newBoxes(b)
}

func (b *boxes) Len() uint32 {
	res := C.ergo_lib_ergo_boxes_len(b.p)
	return uint32(res)
}

func (b *boxes) Get(index uint32) (Box, error) {
	var p C.ErgoBoxPtr

	res := C.ergo_lib_ergo_boxes_get(b.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &box{p: p}
		return newBox(ba), nil
	}

	return nil, nil
}

func (b *boxes) Add(box Box) {
	C.ergo_lib_ergo_boxes_add(box.pointer(), b.p)
}

func (b *boxes) pointer() C.ErgoBoxesPtr {
	return b.p
}

func finalizeBoxes(b *boxes) {
	C.ergo_lib_ergo_boxes_delete(b.p)
}
