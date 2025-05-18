package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"iter"
	"runtime"
	"unsafe"
)

type ByteArray struct {
	p C.ByteArrayPtr
}

func newByteArray(b *ByteArray) *ByteArray {
	runtime.AddCleanup(b, finalizeByteArray, b.p)
	return b
}

func NewByteArray(bytes []byte) (*ByteArray, error) {
	var p C.ByteArrayPtr
	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_lib_byte_array_from_raw_parts((*C.uchar)(byteData), C.uintptr_t(len(bytes)), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ba := &ByteArray{p: p}
	return newByteArray(ba), nil
}

func (b *ByteArray) pointer() C.ByteArrayPtr {
	return b.p
}

func finalizeByteArray(p C.ByteArrayPtr) {
	C.ergo_lib_byte_array_delete(p)
}

type ByteArrays struct {
	p C.ByteArraysPtr
}

func newByteArrays(b *ByteArrays) *ByteArrays {
	runtime.AddCleanup(b, finalizeByteArrays, b.p)
	return b
}

func NewByteArrays() *ByteArrays {
	var p C.ByteArraysPtr
	C.ergo_lib_byte_arrays_new(&p)
	ba := &ByteArrays{p: p}
	return newByteArrays(ba)
}

func (b *ByteArrays) Len() int {
	res := C.ergo_lib_byte_arrays_len(b.p)
	return int(res)
}

func (b *ByteArrays) Get(index int) (*ByteArray, error) {
	var p C.ByteArrayPtr

	res := C.ergo_lib_byte_arrays_get(b.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &ByteArray{p: p}
		return newByteArray(ba), nil
	}

	return nil, nil
}

func (b *ByteArrays) Add(byteArray *ByteArray) {
	C.ergo_lib_byte_arrays_add(byteArray.pointer(), b.p)
}

func (b *ByteArrays) All() iter.Seq2[int, *ByteArray] {
	return func(yield func(int, *ByteArray) bool) {
		for i := 0; i < b.Len(); i++ {
			tk, err := b.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
	}
}

func (b *ByteArrays) pointer() C.ByteArraysPtr {
	return b.p
}

func finalizeByteArrays(p C.ByteArraysPtr) {
	C.ergo_lib_byte_arrays_delete(p)
}
