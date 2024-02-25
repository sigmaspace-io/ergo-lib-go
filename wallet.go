package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type MnemonicGenerator interface {
	Generate() (string, error)
	GenerateFromEntropy(entropy []byte) (string, error)
}

type mnemonicGenerator struct {
	p C.MnemonicGeneratorPtr
}

func NewMnemonicGenerator(language string, strength uint32) (MnemonicGenerator, error) {
	languageStr := C.CString(language)
	defer C.free(unsafe.Pointer(languageStr))

	var p C.MnemonicGeneratorPtr

	errPtr := C.ergo_lib_mnemonic_generator(languageStr, C.uint(strength), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	m := &mnemonicGenerator{p: p}
	runtime.SetFinalizer(m, finalizeMnemonicGenerator)

	return m, nil
}

func (m *mnemonicGenerator) Generate() (string, error) {
	var returnStr C.ReturnString

	returnStr = C.ergo_lib_mnemonic_generator_generate(m.p)
	defer C.ergo_lib_mnemonic_generator_free_mnemonic(returnStr.value)
	err := newError(returnStr.error)
	if err.isError() {
		return "", err.error()
	}

	mnemonic := C.GoString(returnStr.value)

	return mnemonic, nil
}

func (m *mnemonicGenerator) GenerateFromEntropy(entropy []byte) (string, error) {
	var returnStr C.ReturnString

	byteData := C.CBytes(entropy)
	defer C.free(unsafe.Pointer(byteData))

	returnStr = C.ergo_lib_mnemonic_generator_generate_from_entropy(m.p, (*C.uchar)(byteData), C.ulong(len(entropy)))
	defer C.ergo_lib_mnemonic_generator_free_mnemonic(returnStr.value)
	err := newError(returnStr.error)
	if err.isError() {
		return "", err.error()
	}

	mnemonic := C.GoString(returnStr.value)

	return mnemonic, nil
}

func finalizeMnemonicGenerator(m *mnemonicGenerator) {
	C.free(unsafe.Pointer(m.p))
}

type Wallet interface {
}

type wallet struct {
	p C.WalletPtr
}

func NewWallet(mnemonicPhrase string, mnemonicPassword string) (Wallet, error) {
	mnemonic := C.CString(mnemonicPhrase)
	defer C.free(unsafe.Pointer(mnemonic))
	password := C.CString(mnemonicPassword)
	defer C.free(unsafe.Pointer(password))

	var p C.WalletPtr

	errPtr := C.ergo_lib_wallet_from_mnemonic(mnemonic, password, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	w := &wallet{p: p}

	runtime.SetFinalizer(w, finalizeWallet)

	return w, nil
}

func finalizeWallet(w *wallet) {
	C.ergo_lib_wallet_delete(w.p)
}
