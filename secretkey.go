package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"errors"
	"iter"
	"runtime"
	"unsafe"
)

// SecretKey represents secret key for the prover
type SecretKey struct {
	p C.SecretKeyPtr
}

func newSecretKey(s *SecretKey) *SecretKey {
	runtime.AddCleanup(s, finalizeSecretKey, s.p)
	return s
}

// NewSecretKey generates new random SecretKey
func NewSecretKey() *SecretKey {
	var p C.SecretKeyPtr
	C.ergo_lib_secret_key_generate_random(&p)
	s := &SecretKey{p: p}
	return newSecretKey(s)
}

// NewSecretKeyFromBytes parses dlog secret key from bytes (SEC-1-encoded scalar)
// provided secret key bytes must be of length 32
func NewSecretKeyFromBytes(bytes []byte) (*SecretKey, error) {
	if len(bytes) != 32 {
		return nil, errors.New("secret key size must be 32 bytes")
	}

	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	var p C.SecretKeyPtr
	errPtr := C.ergo_lib_secret_key_from_bytes((*C.uchar)(byteData), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	s := &SecretKey{p: p}
	return newSecretKey(s), nil
}

// Address returns Address of the SecretKey
func (s *SecretKey) Address() *Address {
	var p C.AddressPtr
	C.ergo_lib_secret_key_get_address(s.p, &p)
	runtime.KeepAlive(s)
	a := &Address{p}
	return newAddress(a)
}

// Bytes returns SecretKey encoded to bytes
func (s *SecretKey) Bytes() []byte {
	bytes := C.malloc(C.uintptr_t(32))
	C.ergo_lib_secret_key_to_bytes(s.p, (*C.uint8_t)(bytes))
	defer C.free(unsafe.Pointer(bytes))
	runtime.KeepAlive(s)
	result := C.GoBytes(bytes, C.int(32))
	return result
}

func (s *SecretKey) pointer() C.SecretKeyPtr {
	return s.p
}

func finalizeSecretKey(p C.SecretKeyPtr) {
	C.ergo_lib_secret_key_delete(p)
}

// SecretKeys an ordered collection of SecretKey
type SecretKeys struct {
	p C.SecretKeysPtr
}

func newSecretKeys(s *SecretKeys) *SecretKeys {
	runtime.AddCleanup(s, finalizeSecretKeys, s.p)
	return s
}

// NewSecretKeys creates an empty SecretKeys collection
func NewSecretKeys() *SecretKeys {
	var p C.SecretKeysPtr
	C.ergo_lib_secret_keys_new(&p)
	s := &SecretKeys{p: p}
	return newSecretKeys(s)
}

// Len returns the length of the collection
func (s *SecretKeys) Len() int {
	res := C.ergo_lib_secret_keys_len(s.p)
	runtime.KeepAlive(s)
	return int(res)
}

// Get returns the SecretKey at the provided index if it exists
func (s *SecretKeys) Get(index int) (*SecretKey, error) {
	var p C.SecretKeyPtr

	res := C.ergo_lib_secret_keys_get(s.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(s)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		sk := &SecretKey{p: p}
		return newSecretKey(sk), nil
	}

	return nil, nil
}

// Add adds provided SecretKey to the end of the collection
func (s *SecretKeys) Add(secretKey *SecretKey) {
	C.ergo_lib_secret_keys_add(secretKey.pointer(), s.p)
	runtime.KeepAlive(s)
	runtime.KeepAlive(secretKey)
}

// All returns an iterator over all SecretKey inside the collection
func (s *SecretKeys) All() iter.Seq2[int, *SecretKey] {
	return func(yield func(int, *SecretKey) bool) {
		for i := 0; i < s.Len(); i++ {
			tk, err := s.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(s)
	}
}

func (s *SecretKeys) pointer() C.SecretKeysPtr {
	return s.p
}

func finalizeSecretKeys(p C.SecretKeysPtr) {
	C.ergo_lib_secret_keys_delete(p)
}
