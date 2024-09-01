package ergo

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMnemonicToSeed(t *testing.T) {
	mnemonic := "change me do not use me change me do not use me"
	seed := MnemonicToSeed(mnemonic, "")
	assert.Equal(t, "c5b2537b52b27b903b34c423783ced17c489e4385ec6d49d6a19a7f892ecd3917db36675de36bcbe3b8dbc6f803877f4155bdf83482ca5f0fc4282a61ac842a3", hex.EncodeToString(seed))
}

func TestMnemonicToSeeWithPass(t *testing.T) {
	mnemonic := "change me do not use me change me do not use me"
	seed := MnemonicToSeed(mnemonic, "password123")
	assert.Equal(t, "dfe3088b88e2eb8588482e8c56d9cde497c4e1f63fd29b480cbb0ed0227331d51301cfc2d461acce642868ecb618a37b4fd75d48dc6189674c55fbafd807d69c", hex.EncodeToString(seed))
}
