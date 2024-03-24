package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type networkPrefix uint8

const (
	// MainnetPrefix is the network prefix used in mainnet address encoding
	MainnetPrefix networkPrefix = 0

	// TestnetPrefix is the network prefix used in testnet address encoding
	TestnetPrefix = 16
)

type addressTypePrefix uint8

const (
	// P2PkPrefix 0x01 - Pay-to-PublicKey(P2PK) address
	P2PkPrefix addressTypePrefix = 1
	// Pay2ShPrefix 0x02 - Pay-to-Script-Hash(P2SH)
	Pay2ShPrefix addressTypePrefix = 2
	// Pay2SPrefix 0x03 - Pay-to-Script(P2S)
	Pay2SPrefix addressTypePrefix = 3
)

type Address interface {
	// Base58 converts an Address to a base58 string using the provided networkPrefix.
	Base58(prefix networkPrefix) string

	// TypePrefix returns the addressTypePrefix for the Address.
	// 0x01 - Pay-to-PublicKey(P2PK) address.
	// 0x02 - Pay-to-Script-Hash(P2SH).
	// 0x03 - Pay-to-Script(P2S).
	TypePrefix() addressTypePrefix
	pointer() C.AddressPtr
}

type address struct {
	p C.AddressPtr
}

func newAddress(a *address) Address {
	runtime.SetFinalizer(a, finalizeAddress)

	return a
}

// NewAddress creates an Address from a base58 string.
func NewAddress(s string) (Address, error) {
	addressStr := C.CString(s)
	defer C.free(unsafe.Pointer(addressStr))

	var p C.AddressPtr

	errPtr := C.ergo_lib_address_from_base58(addressStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	a := &address{p}

	return newAddress(a), nil
}

func (a *address) Base58(prefix networkPrefix) string {
	var outAddrStr *C.char
	cPrefix := C.uchar(prefix)

	C.ergo_lib_address_to_base58(a.p, cPrefix, &outAddrStr)
	defer C.ergo_lib_delete_string(outAddrStr)

	return C.GoString(outAddrStr)
}

func (a *address) TypePrefix() addressTypePrefix {
	prefix := C.ergo_lib_address_type_prefix(a.p)

	return addressTypePrefix(prefix)
}

func (a *address) pointer() C.AddressPtr {
	return a.p
}

func finalizeAddress(a *address) {
	C.ergo_lib_address_delete(a.p)
}
