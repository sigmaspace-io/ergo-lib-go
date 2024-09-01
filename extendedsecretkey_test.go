package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtendedSecretKey(t *testing.T) {
	seedStr := "chef hidden swift slush bar length outdoor pupil hunt country endorse accuse"
	seed := MnemonicToSeed(seedStr, "")

	root, rootErr := DeriveMaster(seed)
	assert.NoError(t, rootErr)

	changePath, derivationPathErr := NewDerivationPath(0, []uint32{0})
	assert.NoError(t, derivationPathErr)

	changeSecret, deriveErr := root.Derive(changePath)
	assert.NoError(t, deriveErr)

	assert.Equal(t, "9hRTUYF37avZvhC5FG7VoSfrfWQgRMubrA4xLqwFBfes743691r", changeSecret.ExtendedPublicKey().Address().Base58(MainnetPrefix))

	nextChangePath, _ := changePath.Next()
	nextChangeSecret, nextDeriveErr := root.Derive(nextChangePath)
	assert.NoError(t, nextDeriveErr)

	assert.Equal(t, "9gYRhhA9TcFv6xWGwTBPLBJzyW1Hv3EiDzXqoivWYjq8TowWJ1h", nextChangeSecret.ExtendedPublicKey().Address().Base58(MainnetPrefix))
}
