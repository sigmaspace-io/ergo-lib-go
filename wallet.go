package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// MnemonicGenerator generates new mnemonic seed phrases
type MnemonicGenerator interface {
	// Generate a new mnemonic sentence using random entropy
	Generate() (string, error)
	// GenerateFromEntropy generates a new mnemonic sentence using provided entropy
	GenerateFromEntropy(entropy []byte) (string, error)
}

type mnemonicGenerator struct {
	p C.MnemonicGeneratorPtr
}

// NewMnemonicGenerator creates a new MnemonicGenerator based on supplied language and strength
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
	// AddSecret adds a secret to the wallets prover
	AddSecret(secret SecretKey) error
	// SignTransaction signs a transaction
	SignTransaction(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes) (Transaction, error)
	// SignTransactionMulti signs a multi signature transaction
	SignTransactionMulti(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes, txHints TransactionHintsBag) (Transaction, error)
	// SignReducedTransaction signs a reduced transaction (generating proofs for inputs)
	SignReducedTransaction(reducedTx ReducedTransaction) (Transaction, error)
	// SignReducedTransactionMulti signs a multi signature reduced transaction
	SignReducedTransactionMulti(reducedTx ReducedTransaction, txHints TransactionHintsBag) (Transaction, error)
	// GenerateCommitments generates Commitments for unsigned tx
	GenerateCommitments(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes) (TransactionHintsBag, error)
	// GenerateCommitmentsForReducedTransaction generates Commitments for reduced transaction
	GenerateCommitmentsForReducedTransaction(reducedTx ReducedTransaction) (TransactionHintsBag, error)
	// SignMessageUsingP2PK signs an arbitrary message using a P2PK address
	SignMessageUsingP2PK(address Address, message []byte) (SignedMessage, error)
}

type wallet struct {
	p C.WalletPtr
}

func newWallet(w *wallet) Wallet {
	runtime.SetFinalizer(w, finalizeWallet)
	return w
}

// NewWallet creates a Wallet instance loading secret key from mnemonic or throws error if a DlogSecretKey cannot be parsed from the provided phrase
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

	return newWallet(w), nil
}

// NewWalletFromSecretKeys creates a Wallet from secrets
func NewWalletFromSecretKeys(secrets SecretKeys) Wallet {
	var p C.WalletPtr
	C.ergo_lib_wallet_from_secrets(secrets.pointer(), &p)
	w := &wallet{p: p}
	return newWallet(w)
}

func (w *wallet) AddSecret(secret SecretKey) error {
	errPtr := C.ergo_lib_wallet_add_secret(w.p, secret.pointer())
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func (w *wallet) SignTransaction(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes) (Transaction, error) {
	var p C.TransactionPtr
	errPtr := C.ergo_lib_wallet_sign_transaction(w.p, stateContext.pointer(), unsignedTx.pointer(), boxesToSpend.pointer(), dataBoxes.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	t := &transaction{p: p}
	return newTransaction(t), nil
}

func (w *wallet) SignTransactionMulti(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes, txHints TransactionHintsBag) (Transaction, error) {
	var p C.TransactionPtr
	errPtr := C.ergo_lib_wallet_sign_transaction_multi(w.p, stateContext.pointer(), unsignedTx.pointer(), boxesToSpend.pointer(), dataBoxes.pointer(), txHints.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	t := &transaction{p: p}
	return newTransaction(t), nil
}

func (w *wallet) SignReducedTransaction(reducedTx ReducedTransaction) (Transaction, error) {
	var p C.TransactionPtr
	errPtr := C.ergo_lib_wallet_sign_reduced_transaction(w.p, reducedTx.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	t := &transaction{p: p}
	return newTransaction(t), nil
}

func (w *wallet) SignReducedTransactionMulti(reducedTx ReducedTransaction, txHints TransactionHintsBag) (Transaction, error) {
	var p C.TransactionPtr
	errPtr := C.ergo_lib_wallet_sign_reduced_transaction_multi(w.p, reducedTx.pointer(), txHints.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	t := &transaction{p: p}
	return newTransaction(t), nil
}

func (w *wallet) GenerateCommitments(stateContext StateContext, unsignedTx UnsignedTransaction, boxesToSpend Boxes, dataBoxes Boxes) (TransactionHintsBag, error) {
	var p C.TransactionHintsBagPtr
	errPtr := C.ergo_lib_wallet_generate_commitments(w.p, stateContext.pointer(), unsignedTx.pointer(), boxesToSpend.pointer(), dataBoxes.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	th := &transactionHintsBag{p: p}
	return newTransactionHintsBag(th), nil
}

func (w *wallet) GenerateCommitmentsForReducedTransaction(reducedTx ReducedTransaction) (TransactionHintsBag, error) {
	var p C.TransactionHintsBagPtr
	errPtr := C.ergo_lib_wallet_generate_commitments_for_reduced_transaction(w.p, reducedTx.pointer(), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	th := &transactionHintsBag{p: p}
	return newTransactionHintsBag(th), nil
}

func (w *wallet) SignMessageUsingP2PK(address Address, message []byte) (SignedMessage, error) {
	byteData := C.CBytes(message)
	defer C.free(unsafe.Pointer(byteData))

	var p C.SignedMessagePtr
	errPtr := C.ergo_lib_wallet_sign_message_using_p2pk(w.p, address.pointer(), (*C.uchar)(byteData), C.ulong(len(message)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	sm := &signedMessage{p: p}
	return newSignedMessage(sm), nil
}

func finalizeWallet(w *wallet) {
	C.ergo_lib_wallet_delete(w.p)
}

type SignedMessage interface {
	pointer() C.SignedMessagePtr
}

type signedMessage struct {
	p C.SignedMessagePtr
}

func newSignedMessage(s *signedMessage) SignedMessage {
	runtime.SetFinalizer(s, finalizeSignedMessage)
	return s
}

func (s *signedMessage) pointer() C.SignedMessagePtr {
	return s.p
}

func finalizeSignedMessage(s *signedMessage) {
	C.ergo_lib_signed_message_delete(s.p)
}
