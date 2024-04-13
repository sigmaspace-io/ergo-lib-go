package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageSigning(t *testing.T) {
	key := NewSecretKey()
	addr := key.Address()
	keys := NewSecretKeys()
	keys.Add(key)
	testWallet := NewWalletFromSecretKeys(keys)
	msg := []byte("this is a message")
	sig, sigErr := testWallet.SignMessageUsingP2PK(addr, msg)
	assert.NoError(t, sigErr)
	res, verifyErr := VerifySignature(addr, msg, sig)
	assert.NoError(t, verifyErr)
	assert.True(t, res)
}
