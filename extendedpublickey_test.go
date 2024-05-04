package ergo

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExtendedPublicKey(t *testing.T) {
	extPubKeyBytes, _ := hex.DecodeString("02e8445082a72f29b75ca48748a914df60622a609cacfce8ed0e35804560741d29")
	chainCodeBytes, _ := hex.DecodeString("04466b9cc8e161e966409ca52986c584f07e9dc81f735db683c3ff6ec7b1503f")

	testDerivationPath, _ := NewDerivationPath(0, make([]uint32, 1))

	extPubKey, extPubKeyErr := NewExtendedPublicKey(extPubKeyBytes, chainCodeBytes, testDerivationPath)
	assert.NoError(t, extPubKeyErr)
	assert.Equal(t, "9gHMTduN2xseqb5NMKQtNSeS7Pe6wm7AwGoLoMERidWDERQunvn", extPubKey.Address().Base58(MainnetPrefix))
}
