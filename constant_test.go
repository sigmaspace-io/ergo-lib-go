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

	constType, _ := con.Type()

	assert.Equal(t, "SInt", constType)
}

func TestConstant_ConstantType(t *testing.T) {
	con, _ := NewConstant("100102")

	costType, _ := con.Type()

	assert.Equal(t, "SColl(SInt)", costType)
}

func TestConstant_ConstantValue(t *testing.T) {
	con, _ := NewConstant("100204a00b")

	constValue, _ := con.Value()

	assert.Equal(t, "[2,720]", constValue)
}

func TestNewConstant_Invalid(t *testing.T) {
	_, err := NewConstant("0402,")

	assert.Error(t, err)
}

func TestNewConstant_TupleExpression(t *testing.T) {
	con, err := NewConstant("860202660263")

	assert.Nil(t, err)
	constValue, _ := con.Value()
	assert.Equal(t, "BoundedVec{inner:[102,99]}", constValue)
}
