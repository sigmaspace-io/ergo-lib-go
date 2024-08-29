package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextExtension(t *testing.T) {
	cE := NewContextExtension()
	testConst := NewConstantFromInt16(4)

	cE.Set(127, testConst)

	resConst, constErr := cE.Get(127)
	assert.NoError(t, constErr)

	resVal, resValErr := resConst.Int16()
	assert.NoError(t, resValErr)
	assert.Equal(t, int16(4), resVal)
}
