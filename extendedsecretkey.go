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

type ExtendedSecretKey struct {
	p C.ExtSecretKeyPtr
}

func newExtendedSecretKey(e *ExtendedSecretKey) *ExtendedSecretKey {
	runtime.AddCleanup(e, finalizeExtendedSecretKey, e.p)
	return e
}

// NewExtendedSecretKey creates a new ExtendedSecretKey from secretKeyBytes, chainCode and DerivationPath
// secretKeyBytes needs to be the length of 32 bytes
// chainCode needs to be the length of 32 bytes
func NewExtendedSecretKey(secretKeyBytes []byte, chainCode []byte, derivationPath *DerivationPath) (*ExtendedSecretKey, error) {
	if len(secretKeyBytes) != 32 {
		return nil, errors.New("secretKeyBytes must be 32 bytes")
	}

	if len(chainCode) != 32 {
		return nil, errors.New("chainCode must be 32 bytes")
	}

	secretKeyByteData := C.CBytes(secretKeyBytes)
	defer C.free(unsafe.Pointer(secretKeyByteData))
	chainCodeByteData := C.CBytes(chainCode)
	defer C.free(unsafe.Pointer(chainCodeByteData))

	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_new((*C.uchar)(secretKeyByteData), (*C.uchar)(chainCodeByteData), derivationPath.pointer(), &p)
	runtime.KeepAlive(derivationPath)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	e := &ExtendedSecretKey{p: p}
	return newExtendedSecretKey(e), nil
}

// DeriveMaster derives root ExtendedSecretKey from seed bytes
func DeriveMaster(seed []byte) (*ExtendedSecretKey, error) {
	seedByteData := C.CBytes(seed)
	defer C.free(unsafe.Pointer(seedByteData))

	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_derive_master((*C.uchar)(seedByteData), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &ExtendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

// Child derives a new ExtendedSecretKey from the provided index
// The index is in the form of soft or hardened indices
// For example: 4 or 4' respectively
func (e *ExtendedSecretKey) Child(index string) (*ExtendedSecretKey, error) {
	indexStr := C.CString(index)
	defer C.free(unsafe.Pointer(indexStr))

	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_child(e.p, indexStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &ExtendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

// Path returns the DerivationPath of the ExtendedSecretKey
func (e *ExtendedSecretKey) Path() *DerivationPath {
	var p C.DerivationPathPtr
	C.ergo_lib_ext_secret_key_path(e.p, &p)
	runtime.KeepAlive(e)
	d := &DerivationPath{p: p}
	return newDerivationPath(d)
}

// SecretKey returns the SecretKey of the ExtendedSecretKey
func (e *ExtendedSecretKey) SecretKey() *SecretKey {
	var p C.SecretKeyPtr
	C.ergo_lib_ext_secret_key_get_secret_key(e.p, &p)
	runtime.KeepAlive(e)
	s := &SecretKey{p: p}
	return newSecretKey(s)
}

// ExtendedPublicKey returns the ExtendedPublicKey associated with the ExtendedSecretKey
func (e *ExtendedSecretKey) ExtendedPublicKey() *ExtendedPublicKey {
	var p C.ExtPubKeyPtr
	C.ergo_lib_ext_secret_key_public_key(e.p, &p)
	runtime.KeepAlive(e)
	ep := &ExtendedPublicKey{p: p}
	return newExtendedPublicKey(ep)
}

// Derive derives a new ExtendedSecretKey from the supplied DerivationPath
func (e *ExtendedSecretKey) Derive(derivationPath *DerivationPath) (*ExtendedSecretKey, error) {
	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_derive(e.p, derivationPath.pointer(), &p)
	runtime.KeepAlive(derivationPath)
	runtime.KeepAlive(e)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &ExtendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

func finalizeExtendedSecretKey(p C.ExtSecretKeyPtr) {
	C.ergo_lib_ext_secret_key_delete(p)
}
