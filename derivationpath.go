package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type DerivationPath struct {
	p C.DerivationPathPtr
}

func newDerivationPath(d *DerivationPath) *DerivationPath {
	runtime.AddCleanup(d, finalizeDerivationPath, d.p)
	return d
}

// NewDerivationPath creates DerivationPath from account index and Address indices
func NewDerivationPath(account uint32, addressIndices []uint32) (*DerivationPath, error) {
	var p C.DerivationPathPtr

	errPtr := C.ergo_lib_derivation_path_new(C.uint32_t(account), (*C.uint32_t)(&addressIndices[0]), C.uintptr_t(len(addressIndices)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	d := &DerivationPath{p: p}
	return newDerivationPath(d), nil
}

// NewDerivationPathFromString creates DerivationPath from string which
// should be in the form of m/44/429/acc'/0/addr
func NewDerivationPathFromString(s string) (*DerivationPath, error) {
	derivationPathStr := C.CString(s)
	defer C.free(unsafe.Pointer(derivationPathStr))

	var p C.DerivationPathPtr
	errPtr := C.ergo_lib_derivation_path_from_str(derivationPathStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	d := &DerivationPath{p: p}
	return newDerivationPath(d), nil
}

// String returns the DerivationPath formatted as string in the m/44/429/acc'/0/addr format
func (d *DerivationPath) String() string {
	var derivationPathStr *C.char

	C.ergo_lib_derivation_path_to_str(d.p, &derivationPathStr)
	defer C.ergo_lib_delete_string(derivationPathStr)
	runtime.KeepAlive(d)

	return C.GoString(derivationPathStr)
}

// Depth returns the length of the DerivationPath
func (d *DerivationPath) Depth() uint32 {
	res := C.ergo_lib_derivation_path_depth(d.p)
	runtime.KeepAlive(d)
	return uint32(res)
}

// Next returns a new DerivationPath with the last element of the derivation path being increased, e.g. m/1/2 -> m/1/3
func (d *DerivationPath) Next() (*DerivationPath, error) {
	var p C.DerivationPathPtr

	errPtr := C.ergo_lib_derivation_path_next(d.p, &p)
	runtime.KeepAlive(d)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	dp := &DerivationPath{p: p}
	return newDerivationPath(dp), nil
}

func (d *DerivationPath) pointer() C.DerivationPathPtr {
	return d.p
}

func finalizeDerivationPath(p C.DerivationPathPtr) {
	C.ergo_lib_derivation_path_delete(p)
}
