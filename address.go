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
	// MainnetPrefix is the network prefix used in mainnet address encoding.
	MainnetPrefix networkPrefix = 0

	// TestnetPrefix is the network prefix used in testnet address encoding.
	TestnetPrefix = 16
)

type addressTypePrefix uint8

const (
	P2PkPrefix addressTypePrefix = 1
	Pay2ShPrefix
	Pay2SPrefix
)

type Address interface {
	// Base58 converts an Address to a base58 string using the provided networkPrefix.
	Base58(prefix networkPrefix) string

	// TypePrefix returns the networkPrefix for the address.
	// 0x01 - Pay-to-PublicKey(P2PK) address.
	// 0x02 - Pay-to-Script-Hash(P2SH).
	// 0x03 - Pay-to-Script(P2S).
	TypePrefix() addressTypePrefix
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
	defer C.ergo_lib_delete_string(outAddrStr)
	cPrefix := C.uchar(prefix)

	C.ergo_lib_address_to_base58(a.p, cPrefix, &outAddrStr)

	return C.GoString(outAddrStr)
}

func (a *address) TypePrefix() addressTypePrefix {
	prefix := C.ergo_lib_address_type_prefix(a.p)

	return addressTypePrefix(prefix)
}

func finalizeAddress(a *address) {
	C.ergo_lib_address_delete(a.p)
}
