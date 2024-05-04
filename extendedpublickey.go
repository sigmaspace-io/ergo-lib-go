package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type ExtendedPublicKey interface {
	// Child derives a new ExtendedPublicKey from the provided index
	Child(childIndex uint32) (ExtendedPublicKey, error)
	// Derive derives a new ExtendedPublicKey from the supplied DerivationPath
	Derive(derivationPath DerivationPath) (ExtendedPublicKey, error)
	// Address returns the Address associated with the ExtendedPublicKey
	Address() Address
	pointer() C.ExtPubKeyPtr
}

type extendedPublicKey struct {
	p C.ExtPubKeyPtr
}

func newExtendedPublicKey(e *extendedPublicKey) ExtendedPublicKey {
	runtime.SetFinalizer(e, finalizeExtendedPublicKey)
	return e
}

// NewExtendedPublicKey creates a new ExtendedPublicKey from publicKeyBytes, chainCode and derivationPath
// publicKeyBytes needs to be the length of 33 bytes
// chainCode needs to be the length of 32 bytes
func NewExtendedPublicKey(publicKeyBytes []byte, chainCode []byte, derivationPath DerivationPath) (ExtendedPublicKey, error) {
	if len(publicKeyBytes) != 33 {
		return nil, errors.New("secretKeyBytes must be 32 bytes")
	}

	if len(chainCode) != 32 {
		return nil, errors.New("chainCode must be 32 bytes")
	}

	publicKeyByteData := C.CBytes(publicKeyBytes)
	defer C.free(unsafe.Pointer(publicKeyByteData))
	chainCodeByteData := C.CBytes(chainCode)
	defer C.free(unsafe.Pointer(chainCodeByteData))

	var p C.ExtPubKeyPtr
	errPtr := C.ergo_lib_ext_pub_key_new((*C.uchar)(publicKeyByteData), (*C.uchar)(chainCodeByteData), derivationPath.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	e := &extendedPublicKey{p: p}
	return newExtendedPublicKey(e), nil
}

func (e *extendedPublicKey) Child(childIndex uint32) (ExtendedPublicKey, error) {
	var p C.ExtPubKeyPtr
	errPtr := C.ergo_lib_ext_pub_key_child(e.p, C.uint32_t(childIndex), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	return newExtendedPublicKey(e), nil
}

func (e *extendedPublicKey) Derive(derivationPath DerivationPath) (ExtendedPublicKey, error) {
	var p C.ExtPubKeyPtr
	errPtr := C.ergo_lib_ext_pub_key_derive(e.p, derivationPath.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	return newExtendedPublicKey(e), nil
}

func (e *extendedPublicKey) Address() Address {
	var p C.AddressPtr
	C.ergo_lib_ext_pub_key_address(e.p, &p)
	a := &address{p: p}
	return newAddress(a)
}

func (e *extendedPublicKey) pointer() C.ExtPubKeyPtr {
	return e.p
}

func finalizeExtendedPublicKey(e *extendedPublicKey) {
	C.ergo_lib_ext_pub_key_delete(e.p)
}
