package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxId(t *testing.T) {
	txIdStr := "93d344aa527e18e5a221db060ea1a868f46b61e4537e6e5f69ecc40334c15e38"

	testTxId, _ := NewTxId(txIdStr)
	resTxIdStr, _ := testTxId.ToString()

	assert.Equal(t, txIdStr, resTxIdStr)
}
