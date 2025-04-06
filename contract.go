package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Contract defines the contract(script) that will be guarding box contents
type Contract interface {
	// Tree returns the ergo Tree of the Contract
	Tree() Tree
	// Equals checks if provided Contract is same
	Equals(contract Contract) bool
	pointer() C.ContractPtr
}

type contract struct {
	p C.ContractPtr
}

func newContract(c *contract) Contract {
	runtime.AddCleanup(c, finalizeContract, c.p)
	return c
}

// NewContractFromTree creates a new Contract from ergo Tree
func NewContractFromTree(ergoTree Tree) Contract {
	var p C.ContractPtr
	C.ergo_lib_contract_new(ergoTree.pointer(), &p)
	runtime.KeepAlive(ergoTree)

	c := &contract{p: p}

	return newContract(c)
}

// NewContractCompileFromString compiles a contract from ErgoScript source code
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

// NewContractPayToAddress creates a new Contract that allows spending of the guarded box by a given recipient (Address)
func NewContractPayToAddress(payToAddress Address) (Contract, error) {
	var p C.ContractPtr
	errPtr := C.ergo_lib_contract_pay_to_address(payToAddress.pointer(), &p)
	runtime.KeepAlive(payToAddress)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	c := &contract{p: p}

	return newContract(c), nil
}

func (c *contract) Tree() Tree {
	var ergoTreePtr C.ErgoTreePtr
	C.ergo_lib_contract_ergo_tree(c.p, &ergoTreePtr)
	runtime.KeepAlive(c)

	newErgoTree := &tree{p: ergoTreePtr}

	return newTree(newErgoTree)
}

func (c *contract) Equals(contract Contract) bool {
	res := C.ergo_lib_contract_eq(c.p, contract.pointer())
	runtime.KeepAlive(c)
	runtime.KeepAlive(contract)
	return bool(res)
}

func (c *contract) pointer() C.ContractPtr {
	return c.p
}

func finalizeContract(p C.ContractPtr) {
	C.ergo_lib_contract_delete(p)
}
