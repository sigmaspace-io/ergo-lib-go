package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const derivationPathStr = "m/44'/429'/0'/0"

func TestNewDerivationPath(t *testing.T) {
	testDerivationPath, testDerivationPathErr := NewDerivationPath(1, []uint32{3})
	assert.NoError(t, testDerivationPathErr)
	assert.Equal(t, "m/44'/429'/1'/0/3", testDerivationPath.String())
}

func TestDerivationPath_String(t *testing.T) {
	testDerivationPath, testDerivationPathErr := NewDerivationPathFromString(derivationPathStr)
	assert.NoError(t, testDerivationPathErr)
	assert.Equal(t, derivationPathStr, testDerivationPath.String())
}

func TestDerivationPath_Depth(t *testing.T) {
	testDerivationPath, testDerivationPathErr := NewDerivationPathFromString(derivationPathStr)
	assert.NoError(t, testDerivationPathErr)
	assert.Equal(t, uint32(4), testDerivationPath.Depth())
}

func TestDerivationPath_Next(t *testing.T) {
	testDerivationPath, testDerivationPathErr := NewDerivationPathFromString(derivationPathStr)
	assert.NoError(t, testDerivationPathErr)
	nextDerivationPath, nextDerivationPathErr := testDerivationPath.Next()
	assert.NoError(t, nextDerivationPathErr)
	assert.Equal(t, "m/44'/429'/0'/1", nextDerivationPath.String())
}
