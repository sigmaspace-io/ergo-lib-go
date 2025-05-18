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

type ExtendedPublicKey struct {
	p C.ExtPubKeyPtr
}

func newExtendedPublicKey(e *ExtendedPublicKey) *ExtendedPublicKey {
	runtime.AddCleanup(e, finalizeExtendedPublicKey, e.p)
	return e
}

// NewExtendedPublicKey creates a new ExtendedPublicKey from publicKeyBytes, chainCode and DerivationPath
// publicKeyBytes needs to be the length of 33 bytes
// chainCode needs to be the length of 32 bytes
func NewExtendedPublicKey(publicKeyBytes []byte, chainCode []byte, derivationPath *DerivationPath) (*ExtendedPublicKey, error) {
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
	runtime.KeepAlive(derivationPath)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	e := &ExtendedPublicKey{p: p}
	return newExtendedPublicKey(e), nil
}

// Child derives a new ExtendedPublicKey from the provided index
func (e *ExtendedPublicKey) Child(childIndex uint32) (*ExtendedPublicKey, error) {
	var p C.ExtPubKeyPtr
	errPtr := C.ergo_lib_ext_pub_key_child(e.p, C.uint32_t(childIndex), &p)
	runtime.KeepAlive(e)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	return newExtendedPublicKey(&ExtendedPublicKey{p}), nil
}

// Derive derives a new ExtendedPublicKey from the supplied DerivationPath
func (e *ExtendedPublicKey) Derive(derivationPath *DerivationPath) (*ExtendedPublicKey, error) {
	var p C.ExtPubKeyPtr
	errPtr := C.ergo_lib_ext_pub_key_derive(e.p, derivationPath.pointer(), &p)
	runtime.KeepAlive(e)
	runtime.KeepAlive(derivationPath)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	return newExtendedPublicKey(&ExtendedPublicKey{p}), nil
}

// Address returns the Address associated with the ExtendedPublicKey
func (e *ExtendedPublicKey) Address() *Address {
	var p C.AddressPtr
	C.ergo_lib_ext_pub_key_address(e.p, &p)
	runtime.KeepAlive(e)
	a := &Address{p: p}
	return newAddress(a)
}

func (e *ExtendedPublicKey) pointer() C.ExtPubKeyPtr {
	return e.p
}

func finalizeExtendedPublicKey(p C.ExtPubKeyPtr) {
	C.ergo_lib_ext_pub_key_delete(p)
}
