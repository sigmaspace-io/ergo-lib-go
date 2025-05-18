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
type Tree struct {
	p C.ErgoTreePtr
}

func newTree(t *Tree) *Tree {
	runtime.AddCleanup(t, finalizeTree, t.p)
	return t
}

// NewTree creates a new ergo Tree from the supplied base16 string.
func NewTree(s string) (*Tree, error) {
	treeStr := C.CString(s)
	defer C.free(unsafe.Pointer(treeStr))

	var p C.ErgoTreePtr

	errPtr := C.ergo_lib_ergo_tree_from_base16_bytes(treeStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &Tree{p}

	return newTree(t), nil
}

// Base16 converts the Tree to a base16 encoded string.
func (t *Tree) Base16() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_ergo_tree_to_base16_bytes(t.p, &outStr)
	runtime.KeepAlive(t)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// Address converts the Tree to an Address.
func (t *Tree) Address() (*Address, error) {
	var p C.AddressPtr

	errPtr := C.ergo_lib_address_from_ergo_tree(t.p, &p)
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	a := &Address{p}

	return newAddress(a), nil
}

// TemplateBytesLength determines the length of the byte array
func (t *Tree) TemplateBytesLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_template_bytes_len(t.p)
	runtime.KeepAlive(t)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	size := C.ulong(returnNum.value)

	return int(size), nil
}

// TemplateHash returns the hash of the template bytes as string
func (t *Tree) TemplateHash() (string, error) {
	bytesLength, byteErr := t.TemplateBytesLength()
	if byteErr != nil {
		return "", byteErr
	}

	output := C.malloc(C.uintptr_t(bytesLength))
	defer C.free(unsafe.Pointer(output))

	errPtr := C.ergo_lib_ergo_tree_template_bytes(t.p, (*C.uint8_t)(output))
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoBytes(unsafe.Pointer(output), C.int(bytesLength))

	hash := sha256.Sum256(result[:])
	return hex.EncodeToString(hash[:]), nil
}

// ConstantsLength returns the number of constants stored in the serialized ErgoTree or throws error if the parsing of constants failed
func (t *Tree) ConstantsLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_constants_len(t.p)
	runtime.KeepAlive(t)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	length := C.ulong(returnNum.value)

	return int(length), nil
}

// Constant returns Constant with given index (as stored in serialized ErgoTree) if it exists or throws error if the parsing of constants failed
func (t *Tree) Constant(index int) (*Constant, error) {
	var constantOut C.ConstantPtr
	var returnOption C.ReturnOption

	indexNumber := C.uintptr_t(index)

	returnOption = C.ergo_lib_ergo_tree_get_constant(t.p, indexNumber, &constantOut)
	runtime.KeepAlive(t)
	err := newError(returnOption.error)

	if err.isError() {
		return &Constant{}, err.error()
	}

	c := &Constant{p: constantOut}

	return newConstant(c), nil
}

// Constants returns all Constant within the Tree or throws error if the parsing of constants failed
func (t *Tree) Constants() ([]*Constant, error) {
	length, err := t.ConstantsLength()
	if err != nil {
		return nil, err
	}
	var constants []*Constant
	for i := 0; i < length; i++ {
		ergoTreeConstant, constErr := t.Constant(i)
		if constErr != nil {
			return nil, constErr
		}
		constants = append(constants, ergoTreeConstant)
	}
	runtime.KeepAlive(t)
	return constants, nil
}

// Equals checks if provided Tree is same
func (t *Tree) Equals(tree *Tree) bool {
	res := C.ergo_lib_ergo_tree_eq(t.p, tree.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(tree)
	return bool(res)
}

func (t *Tree) pointer() C.ErgoTreePtr {
	return t.p
}

func finalizeTree(p C.ErgoTreePtr) {
	C.ergo_lib_ergo_tree_delete(p)
}
