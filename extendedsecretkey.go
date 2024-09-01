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

type ExtendedSecretKey interface {
	// Child derives a new ExtendedSecretKey from the provided index
	// The index is in the form of soft or hardened indices
	// For example: 4 or 4' respectively
	Child(index string) (ExtendedSecretKey, error)
	// Path returns the DerivationPath of the ExtendedSecretKey
	Path() DerivationPath
	// SecretKey returns the SecretKey of the ExtendedSecretKey
	SecretKey() SecretKey
	// ExtendedPublicKey returns the ExtendedPublicKey associated with the ExtendedSecretKey
	ExtendedPublicKey() ExtendedPublicKey
	// Derive derives a new ExtendedSecretKey from the supplied DerivationPath
	Derive(derivationPath DerivationPath) (ExtendedSecretKey, error)
}

type extendedSecretKey struct {
	p C.ExtSecretKeyPtr
}

func newExtendedSecretKey(e *extendedSecretKey) ExtendedSecretKey {
	runtime.SetFinalizer(e, finalizeExtendedSecretKey)
	return e
}

// NewExtendedSecretKey creates a new ExtendedSecretKey from secretKeyBytes, chainCode and derivationPath
// secretKeyBytes needs to be the length of 32 bytes
// chainCode needs to be the length of 32 bytes
func NewExtendedSecretKey(secretKeyBytes []byte, chainCode []byte, derivationPath DerivationPath) (ExtendedSecretKey, error) {
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
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	e := &extendedSecretKey{p: p}
	return newExtendedSecretKey(e), nil
}

// DeriveMaster derives root ExtendedSecretKey from seed bytes
func DeriveMaster(seed []byte) (ExtendedSecretKey, error) {
	seedByteData := C.CBytes(seed)
	defer C.free(unsafe.Pointer(seedByteData))

	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_derive_master((*C.uchar)(seedByteData), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &extendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

func (e *extendedSecretKey) Child(index string) (ExtendedSecretKey, error) {
	indexStr := C.CString(index)
	defer C.free(unsafe.Pointer(indexStr))

	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_child(e.p, indexStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &extendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

func (e *extendedSecretKey) Path() DerivationPath {
	var p C.DerivationPathPtr
	C.ergo_lib_ext_secret_key_path(e.p, &p)
	d := &derivationPath{p: p}
	return newDerivationPath(d)
}

func (e *extendedSecretKey) SecretKey() SecretKey {
	var p C.SecretKeyPtr
	C.ergo_lib_ext_secret_key_get_secret_key(e.p, &p)
	s := &secretKey{p: p}
	return newSecretKey(s)
}

func (e *extendedSecretKey) ExtendedPublicKey() ExtendedPublicKey {
	var p C.ExtPubKeyPtr
	C.ergo_lib_ext_secret_key_public_key(e.p, &p)
	ep := &extendedPublicKey{p: p}
	return newExtendedPublicKey(ep)
}

func (e *extendedSecretKey) Derive(derivationPath DerivationPath) (ExtendedSecretKey, error) {
	var p C.ExtSecretKeyPtr
	errPtr := C.ergo_lib_ext_secret_key_derive(e.p, derivationPath.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	es := &extendedSecretKey{p: p}
	return newExtendedSecretKey(es), nil
}

func finalizeExtendedSecretKey(e *extendedSecretKey) {
	C.ergo_lib_ext_secret_key_delete(e.p)
}
