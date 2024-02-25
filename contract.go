package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Contract interface {
	ErgoTree() Tree
	pointer() C.ContractPtr
}

type contract struct {
	p C.ContractPtr
}

func newContract(c *contract) Contract {
	runtime.SetFinalizer(c, finalizeContract)
	return c
}

func NewContractFromTree(ergoTree Tree) Contract {
	var p C.ContractPtr
	C.ergo_lib_contract_new(ergoTree.pointer(), &p)

	c := &contract{p: p}

	return newContract(c)
}

func NewContractCompileFromString(compileFromString string) (Contract, error) {
	contractStr := C.CString(compileFromString)
	defer C.free(unsafe.Pointer(contractStr))

	var p C.ContractPtr
	errPtr := C.ergo_lib_contract_compile(contractStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	c := &contract{p: p}

	return newContract(c), nil
}

func NewContractPayToAddress(payToAddress Address) (Contract, error) {
	var p C.ContractPtr
	errPtr := C.ergo_lib_contract_pay_to_address(payToAddress.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	c := &contract{p: p}

	return newContract(c), nil
}

func (c *contract) ErgoTree() Tree {
	var ergoTreePtr C.ErgoTreePtr
	C.ergo_lib_contract_ergo_tree(c.p, &ergoTreePtr)

	newErgoTree := &tree{p: ergoTreePtr}

	return newTree(newErgoTree)
}

func (c *contract) pointer() C.ContractPtr {
	return c.p
}

func finalizeContract(c *contract) {
	C.ergo_lib_contract_delete(c.p)
}
