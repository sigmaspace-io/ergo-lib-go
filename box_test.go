package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoxId_Base16(t *testing.T) {
	boxId, _ := NewBoxId("8452e43011f522a3432a04e4aa77e293fc8c3817a11a2088da49201b88158f8a")

	assert.Equal(t, "8452e43011f522a3432a04e4aa77e293fc8c3817a11a2088da49201b88158f8a", boxId.Base16())
}

func TestNewBoxId_Invalid(t *testing.T) {
	_, err := NewBoxId("8452e43011f522a3432a04e4aa77e293fc8c3817a11a2088da49201b88158f8a,")

	assert.Error(t, err)
}
