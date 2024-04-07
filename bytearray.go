package ergo

/*
#include "ergo.h"
*/
import "C"
import (
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

	errPtr := C.ergo_lib_byte_array_from_raw_parts((*C.uchar)(byteData), C.ulong(len(bytes)), &p)
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
	Len() uint32
	Get(index uint32) (ByteArray, error)
	Add(byteArray ByteArray)
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

func (b *byteArrays) Len() uint32 {
	res := C.ergo_lib_byte_arrays_len(b.p)
	return uint32(res)
}

func (b *byteArrays) Get(index uint32) (ByteArray, error) {
	var p C.ByteArrayPtr

	res := C.ergo_lib_byte_arrays_get(b.p, C.ulong(index), &p)
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

func (b *byteArrays) pointer() C.ByteArraysPtr {
	return b.p
}

func finalizeByteArrays(b *byteArrays) {
	C.ergo_lib_byte_arrays_delete(b.p)
}
