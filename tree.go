package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"unsafe"
)

// Tree is the root of ErgoScript IR. Serialized instances of Tree are self-sufficient and can be passed around
type Tree interface {
	// Base16 converts the Tree to a base16 encoded string.
	Base16() (string, error)
	// Address converts the Tree to an Address.
	Address() (Address, error)
	// TemplateBytesLength determines the length of the byte array
	TemplateBytesLength() (int, error)
	// TemplateHash returns the hash of the template bytes as string
	TemplateHash() (string, error)
	// ConstantsLength returns the number of constants stored in the serialized ErgoTree or throws error if the parsing of constants failed
	ConstantsLength() (int, error)
	// Constant returns Constant with given index (as stored in serialized ErgoTree) if it exists or throws error if the parsing of constants failed
	Constant(index int) (Constant, error)
	// Constants returns all Constant within the Tree or throws error if the parsing of constants failed
	Constants() ([]Constant, error)
	// Equals checks if provided Tree is same
	Equals(tree Tree) bool
	pointer() C.ErgoTreePtr
}

type tree struct {
	p C.ErgoTreePtr
}

func newTree(t *tree) Tree {
	runtime.SetFinalizer(t, finalizeTree)
	return t
}

// NewTree creates a new ergo tree from the supplied base16 string.
func NewTree(s string) (Tree, error) {
	treeStr := C.CString(s)
	defer C.free(unsafe.Pointer(treeStr))

	var p C.ErgoTreePtr

	errPtr := C.ergo_lib_ergo_tree_from_base16_bytes(treeStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &tree{p}

	return newTree(t), nil
}

func (t *tree) Base16() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_tree_to_base16_bytes(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (t *tree) Address() (Address, error) {
	var p C.AddressPtr

	errPtr := C.ergo_lib_address_from_ergo_tree(t.p, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	a := &address{p}

	return newAddress(a), nil
}

func (t *tree) TemplateBytesLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_template_bytes_len(t.p)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	size := C.ulong(returnNum.value)

	return int(size), nil
}

func (t *tree) TemplateHash() (string, error) {
	bytesLength, byteErr := t.TemplateBytesLength()
	if byteErr != nil {
		return "", byteErr
	}

	output := C.malloc(C.uintptr_t(bytesLength))
	defer C.free(unsafe.Pointer(output))

	errPtr := C.ergo_lib_ergo_tree_template_bytes(t.p, (*C.uint8_t)(output))
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoBytes(unsafe.Pointer(output), C.int(bytesLength))

	hash := sha256.Sum256(result[:])
	return hex.EncodeToString(hash[:]), nil
}

func (t *tree) ConstantsLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_constants_len(t.p)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	length := C.ulong(returnNum.value)

	return int(length), nil
}

func (t *tree) Constant(index int) (Constant, error) {
	var constantOut C.ConstantPtr
	var returnOption C.ReturnOption

	indexNumber := C.uintptr_t(index)

	returnOption = C.ergo_lib_ergo_tree_get_constant(t.p, indexNumber, &constantOut)
	err := newError(returnOption.error)

	if err.isError() {
		return &constant{}, err.error()
	}

	c := &constant{p: constantOut}

	return newConstant(c), nil
}

func (t *tree) Constants() ([]Constant, error) {
	length, err := t.ConstantsLength()
	if err != nil {
		return nil, err
	}
	var constants []Constant
	for i := 0; i < length; i++ {
		ergoTreeConstant, constErr := t.Constant(i)
		if constErr != nil {
			return nil, constErr
		}
		constants = append(constants, ergoTreeConstant)
	}
	return constants, nil
}

func (t *tree) Equals(tree Tree) bool {
	res := C.ergo_lib_ergo_tree_eq(t.p, tree.pointer())
	return bool(res)
}

func (t *tree) pointer() C.ErgoTreePtr {
	return t.p
}

func finalizeTree(t *tree) {
	C.ergo_lib_ergo_tree_delete(t.p)
}
