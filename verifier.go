package ergo

/*
#include "ergo.h"
*/
import "C"
import "unsafe"

// VerifySignature verifies that the signature is presented to satisfy SigmaProp conditions
func VerifySignature(address Address, message []byte, signature SignedMessage) (bool, error) {
	byteData := C.CBytes(message)
	defer C.free(unsafe.Pointer(byteData))

	res := C.ergo_lib_verify_signature(address.pointer(), (*C.uchar)(byteData), C.uintptr_t(len(message)), signature.pointer())
	err := newError(res.error)
	if err.isError() {
		return false, err.error()
	}
	return bool(res.value), nil
}
