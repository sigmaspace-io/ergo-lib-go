package ergo

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewMnemonicGenerator(t *testing.T) {
	mnemonicGenerator, _ := NewMnemonicGenerator("english", 128)
	assert.NotNil(t, mnemonicGenerator)
}

func TestMnemonicGenerator_Generate(t *testing.T) {
	strengths := map[uint32]int{
		128: 12,
		160: 15,
		192: 18,
		224: 21,
		256: 24,
	}

	for strength := range strengths {
		mnemonicGenerator, _ := NewMnemonicGenerator("english", strength)
		mnemonic, _ := mnemonicGenerator.Generate()
		words := strings.Split(mnemonic, " ")
		assert.Equal(t, strengths[strength], len(words))
	}
}

func TestMnemonicGenerator_GenerateFromEntropy(t *testing.T) {
	mnemonicGenerator, _ := NewMnemonicGenerator("english", 128)
	entropy := []byte{39, 77, 111, 111, 102, 33, 39, 0, 39, 77, 111, 111, 102, 33, 39, 0}

	mnemonic, err := mnemonicGenerator.GenerateFromEntropy(entropy)

	assert.Nil(t, err)
	assert.Equal(t, "chef hidden swift slush bar length outdoor pupil hunt country endorse accuse", mnemonic)
}

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet("chef hidden swift slush bar length outdoor pupil hunt country endorse accuse", "testPass")

	assert.Nil(t, err)
	assert.NotNil(t, wallet)
}
