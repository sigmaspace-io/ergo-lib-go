package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// SecretKey represents secret key for the prover
type SecretKey interface {
	// Address returns address of the SecretKey
	Address() Address
	// Bytes returns SecretKey encoded to bytes
	Bytes() []byte
	pointer() C.SecretKeyPtr
}

type secretKey struct {
	p C.SecretKeyPtr
}

func newSecretKey(s *secretKey) SecretKey {
	runtime.SetFinalizer(s, finalizeSecretKey)
	return s
}

// NewSecretKey generates new random SecretKey
func NewSecretKey() SecretKey {
	var p C.SecretKeyPtr
	C.ergo_lib_secret_key_generate_random(&p)
	s := &secretKey{p: p}
	return newSecretKey(s)
}

// NewSecretKeyFromBytes parses dlog secret key from bytes (SEC-1-encoded scalar)
func NewSecretKeyFromBytes(bytes []byte) (SecretKey, error) {
	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	var p C.SecretKeyPtr
	errPtr := C.ergo_lib_secret_key_from_bytes((*C.uchar)(byteData), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	s := &secretKey{p: p}
	return newSecretKey(s), nil
}

func (s *secretKey) Address() Address {
	var p C.AddressPtr
	C.ergo_lib_secret_key_get_address(s.p, &p)
	a := &address{p}
	return newAddress(a)
}

func (s *secretKey) Bytes() []byte {
	bytes := C.malloc(C.ulong(32))
	C.ergo_lib_secret_key_to_bytes(s.p, (*C.uint8_t)(bytes))
	defer C.free(unsafe.Pointer(bytes))
	result := C.GoBytes(bytes, C.int(32))
	return result
}

func (s *secretKey) pointer() C.SecretKeyPtr {
	return s.p
}

func finalizeSecretKey(s *secretKey) {
	C.ergo_lib_secret_key_delete(s.p)
}

// SecretKeys an ordered collection of SecretKey
type SecretKeys interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the SecretKey at the provided index if it exists
	Get(index uint32) (SecretKey, error)
	// Add adds provided SecretKey to the end of the collection
	Add(secretKey SecretKey)
	pointer() C.SecretKeysPtr
}

type secretKeys struct {
	p C.SecretKeysPtr
}

func newSecretKeys(s *secretKeys) SecretKeys {
	runtime.SetFinalizer(s, finalizeSecretKeys)
	return s
}

// NewSecretKeys creates an empty SecretKeys collection
func NewSecretKeys() SecretKeys {
	var p C.SecretKeysPtr
	C.ergo_lib_secret_keys_new(&p)
	s := &secretKeys{p: p}
	return newSecretKeys(s)
}

func (s *secretKeys) Len() uint32 {
	res := C.ergo_lib_secret_keys_len(s.p)
	return uint32(res)
}

func (s *secretKeys) Get(index uint32) (SecretKey, error) {
	var p C.SecretKeyPtr

	res := C.ergo_lib_secret_keys_get(s.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		sk := &secretKey{p: p}
		return newSecretKey(sk), nil
	}

	return nil, nil
}

func (s *secretKeys) Add(secretKey SecretKey) {
	C.ergo_lib_secret_keys_add(secretKey.pointer(), s.p)
}

func (s *secretKeys) pointer() C.SecretKeysPtr {
	return s.p
}

func finalizeSecretKeys(s *secretKeys) {
	C.ergo_lib_secret_keys_delete(s.p)
}
