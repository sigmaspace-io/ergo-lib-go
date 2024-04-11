package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecretKey(t *testing.T) {
	key := NewSecretKey()
	bytes := key.Bytes()
	newKey, newKeyErr := NewSecretKeyFromBytes(bytes)
	assert.NoError(t, newKeyErr)
	assert.Equal(t, key, newKey)
}
