package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErgoError_IsError(t *testing.T) {
	_, err := NewAddress("9hdxkYakTHWXR992umPcvh8bAEGG9Sdoi7uW8TKXk1enXCDFBVJ,")

	assert.Error(t, err)
}

func TestErgoError_Error(t *testing.T) {
	_, err := NewAddress("9hdxkYakTHWXR992umPcvh8bAEGG9Sdoi7uW8TKXk1enXCDFBVJ,")

	assert.EqualError(t, err, "error: Base58 decoding error: provided string contained invalid character ',' at byte 51")
}

func TestErgoError_Error_NilError(t *testing.T) {
	err := newError(nil)

	assert.Nil(t, err.error())
}
