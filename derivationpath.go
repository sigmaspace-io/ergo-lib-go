package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type DerivationPath interface {
	// String returns the DerivationPath formatted as string in the m/44/429/acc'/0/addr format
	String() string
	// Depth returns the length of the DerivationPath
	Depth() uint32
	// Next returns a new DerivationPath with the last element of the derivation path being increased, e.g. m/1/2 -> m/1/3
	Next() (DerivationPath, error)
	pointer() C.DerivationPathPtr
}

type derivationPath struct {
	p C.DerivationPathPtr
}

func newDerivationPath(d *derivationPath) DerivationPath {
	runtime.SetFinalizer(d, finalizeDerivationPath)
	return d
}

// NewDerivationPath creates DerivationPath from account index and address indices
func NewDerivationPath(account uint32, addressIndices []uint32) (DerivationPath, error) {
	var p C.DerivationPathPtr

	errPtr := C.ergo_lib_derivation_path_new(C.uint32_t(account), (*C.uint32_t)(&addressIndices[0]), C.uintptr_t(len(addressIndices)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	d := &derivationPath{p: p}
	return newDerivationPath(d), nil
}

// NewDerivationPathFromString creates DerivationPath from string which
// should be in the form of m/44/429/acc'/0/addr
func NewDerivationPathFromString(s string) (DerivationPath, error) {
	derivationPathStr := C.CString(s)
	defer C.free(unsafe.Pointer(derivationPathStr))

	var p C.DerivationPathPtr
	errPtr := C.ergo_lib_derivation_path_from_str(derivationPathStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	d := &derivationPath{p: p}
	return newDerivationPath(d), nil
}

func (d *derivationPath) String() string {
	var derivationPathStr *C.char

	C.ergo_lib_derivation_path_to_str(d.p, &derivationPathStr)
	defer C.ergo_lib_delete_string(derivationPathStr)

	return C.GoString(derivationPathStr)
}

func (d *derivationPath) Depth() uint32 {
	return uint32(C.ergo_lib_derivation_path_depth(d.p))
}

func (d *derivationPath) Next() (DerivationPath, error) {
	var p C.DerivationPathPtr

	errPtr := C.ergo_lib_derivation_path_next(d.p, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	dp := &derivationPath{p: p}
	return newDerivationPath(dp), nil
}

func (d *derivationPath) pointer() C.DerivationPathPtr {
	return d.p
}

func finalizeDerivationPath(d *derivationPath) {
	C.ergo_lib_derivation_path_delete(d.p)
}
