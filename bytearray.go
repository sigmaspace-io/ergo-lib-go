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

type ByteArray interface {
	pointer() C.ByteArrayPtr
}

type byteArray struct {
	p C.ByteArrayPtr
}

func newByteArray(b *byteArray) ByteArray {
	runtime.SetFinalizer(b, finalizeByteArray)
	return b
}

func NewByteArray(bytes []byte) (ByteArray, error) {
	var p C.ByteArrayPtr
	byteData := C.CBytes(bytes)
	defer C.free(unsafe.Pointer(byteData))

	errPtr := C.ergo_lib_byte_array_from_raw_parts((*C.uchar)(byteData), C.uintptr_t(len(bytes)), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	ba := &byteArray{p: p}
	return newByteArray(ba), nil
}

func (b *byteArray) pointer() C.ByteArrayPtr {
	return b.p
}

func finalizeByteArray(b *byteArray) {
	C.ergo_lib_byte_array_delete(b.p)
}

type ByteArrays interface {
	Len() int
	Get(index int) (ByteArray, error)
	Add(byteArray ByteArray)
	All() iter.Seq2[int, ByteArray]
	pointer() C.ByteArraysPtr
}

type byteArrays struct {
	p C.ByteArraysPtr
}

func newByteArrays(b *byteArrays) ByteArrays {
	runtime.SetFinalizer(b, finalizeByteArrays)
	return b
}

func NewByteArrays() ByteArrays {
	var p C.ByteArraysPtr
	C.ergo_lib_byte_arrays_new(&p)
	ba := &byteArrays{p: p}
	return newByteArrays(ba)
}

func (b *byteArrays) Len() int {
	res := C.ergo_lib_byte_arrays_len(b.p)
	return int(res)
}

func (b *byteArrays) Get(index int) (ByteArray, error) {
	var p C.ByteArrayPtr

	res := C.ergo_lib_byte_arrays_get(b.p, C.uintptr_t(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		ba := &byteArray{p: p}
		return newByteArray(ba), nil
	}

	return nil, nil
}

func (b *byteArrays) Add(byteArray ByteArray) {
	C.ergo_lib_byte_arrays_add(byteArray.pointer(), b.p)
}

func (b *byteArrays) All() iter.Seq2[int, ByteArray] {
	return func(yield func(int, ByteArray) bool) {
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

func (b *byteArrays) pointer() C.ByteArraysPtr {
	return b.p
}

func finalizeByteArrays(b *byteArrays) {
	C.ergo_lib_byte_arrays_delete(b.p)
}
