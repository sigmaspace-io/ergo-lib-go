package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstant_Base16(t *testing.T) {
	con, _ := NewConstant("0402")

	value, _ := con.Base16()

	assert.Equal(t, "0402", value)
}

func TestConstantCollByte_Base16(t *testing.T) {
	con, _ := NewConstant("0e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a70173007301")

	value, _ := con.Base16()

	assert.Equal(t, "0e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a70173007301", value)
}

func TestConstant_ConstantTypeSint(t *testing.T) {
	con, _ := NewConstant("0402")

	constType, _ := con.ConstantType()

	assert.Equal(t, "SInt", constType)
}

func TestConstant_ConstantType(t *testing.T) {
	con, _ := NewConstant("100102")

	costType, _ := con.ConstantType()

	assert.Equal(t, "SColl(SInt)", costType)
}
