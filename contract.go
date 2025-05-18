package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Contract defines the Contract(script) that will be guarding Box contents
type Contract struct {
	p C.ContractPtr
}

func newContract(c *Contract) *Contract {
	runtime.AddCleanup(c, finalizeContract, c.p)
	return c
}

// NewContractFromTree creates a new Contract from ergo Tree
func NewContractFromTree(ergoTree *Tree) *Contract {
	var p C.ContractPtr
	C.ergo_lib_contract_new(ergoTree.pointer(), &p)
	runtime.KeepAlive(ergoTree)

	c := &Contract{p: p}

	return newContract(c)
}

// NewContractCompileFromString compiles a Contract from ErgoScript source code
func NewContractCompileFromString(compileFromString string) (*Contract, error) {
	contractStr := C.CString(compileFromString)
	defer C.free(unsafe.Pointer(contractStr))

	var p C.ContractPtr
	errPtr := C.ergo_lib_contract_compile(contractStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	c := &Contract{p: p}

	return newContract(c), nil
}

// NewContractPayToAddress creates a new Contract that allows spending of the guarded Box by a given recipient (Address)
func NewContractPayToAddress(payToAddress *Address) (*Contract, error) {
	var p C.ContractPtr
	errPtr := C.ergo_lib_contract_pay_to_address(payToAddress.pointer(), &p)
	runtime.KeepAlive(payToAddress)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	c := &Contract{p: p}

	return newContract(c), nil
}

// Tree returns the ergo Tree of the Contract
func (c *Contract) Tree() *Tree {
	var ergoTreePtr C.ErgoTreePtr
	C.ergo_lib_contract_ergo_tree(c.p, &ergoTreePtr)
	runtime.KeepAlive(c)

	newErgoTree := &Tree{p: ergoTreePtr}

	return newTree(newErgoTree)
}

// Equals checks if provided Contract is same
func (c *Contract) Equals(contract *Contract) bool {
	res := C.ergo_lib_contract_eq(c.p, contract.pointer())
	runtime.KeepAlive(c)
	runtime.KeepAlive(contract)
	return bool(res)
}

func (c *Contract) pointer() C.ContractPtr {
	return c.p
}

func finalizeContract(p C.ContractPtr) {
	C.ergo_lib_contract_delete(p)
}
