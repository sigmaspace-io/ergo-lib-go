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

type Tree interface {
	// Base16 converts the Tree to a base16 encoded string.
	Base16() (*string, error)

	// Address converts the Tree to an Address.
	Address() (Address, error)
	ErgoTreeTemplateBytesLength() (int, error)
	ErgoTreeTemplateHash() (string, error)
	ErgoTreeConstantsLength() (int, error)
	ErgoTreeGetConstant(index int) (Constant, error)
	ErgoTreeGetConstants() ([]Constant, error)
}

type tree struct {
	p C.ErgoTreePtr
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

	runtime.SetFinalizer(t, finalizeTree)

	return t, nil
}

func (t *tree) Base16() (*string, error) {
	var outStr *C.char
	defer C.free(unsafe.Pointer(outStr))

	errPtr := C.ergo_lib_ergo_tree_to_base16_bytes(t.p, &outStr)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	result := C.GoString(outStr)

	return &result, nil
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

func (t *tree) ErgoTreeTemplateBytesLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_template_bytes_len(t.p)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	size := C.ulong(returnNum.value)

	return int(size), nil
}

func (t *tree) ErgoTreeTemplateHash() (string, error) {
	bytesLength, byteErr := t.ErgoTreeTemplateBytesLength()
	if byteErr != nil {
		return "", byteErr
	}

	output := C.malloc(C.ulong(bytesLength))
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

func (t *tree) ErgoTreeConstantsLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_ergo_tree_constants_len(t.p)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	length := C.ulong(returnNum.value)

	return int(length), nil
}

func (t *tree) ErgoTreeGetConstant(index int) (Constant, error) {
	var constantOut C.ConstantPtr
	var returnOption C.ReturnOption

	indexNumber := C.ulong(index)

	returnOption = C.ergo_lib_ergo_tree_get_constant(t.p, indexNumber, &constantOut)
	err := newError(returnOption.error)

	if err.isError() {
		return &constant{}, err.error()
	}

	c := &constant{p: constantOut}

	return newConstant(c), nil
}

func (t *tree) ErgoTreeGetConstants() ([]Constant, error) {
	length, err := t.ErgoTreeConstantsLength()
	if err != nil {
		return nil, err
	}
	var constants []Constant
	for i := 0; i < length; i++ {
		ergoTreeConstant, constErr := t.ErgoTreeGetConstant(i)
		if constErr != nil {
			return nil, constErr
		}
		constants = append(constants, ergoTreeConstant)
	}
	return constants, nil
}

func finalizeTree(t *tree) {
	C.ergo_lib_ergo_tree_delete(t.p)
}
