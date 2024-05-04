package ergo

/*
#include "ergo.h"
*/
import "C"
import "unsafe"

// MnemonicToSeed converts a mnemonic phrase into a mnemonic seed
// mnemonicPassword is optional and is used to salt the seed
func MnemonicToSeed(mnemonicPhrase string, mnemonicPassword string) []byte {
	mnemonic := C.CString(mnemonicPhrase)
	defer C.free(unsafe.Pointer(mnemonic))

	password := C.CString(mnemonicPassword)
	defer C.free(unsafe.Pointer(password))

	bytes := C.malloc(C.uintptr_t(512 / 8))
	C.ergo_lib_mnemonic_to_seed(mnemonic, password, (*C.uint8_t)(bytes))
	defer C.free(unsafe.Pointer(bytes))

	result := C.GoBytes(bytes, C.int(512/8))
	return result
}
