package ergo

import (
	"encoding/hex"
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

func TestNewConstantFromInt16(t *testing.T) {
	testValue := int16(127)
	c := NewConstantFromInt16(testValue)
	res, _ := c.Int16()
	assert.Equal(t, testValue, res)
}

func TestConstant_Int32(t *testing.T) {
	testValue := int32(999999999)
	c := NewConstantFromInt32(testValue)
	res, _ := c.Int32()
	assert.Equal(t, testValue, res)
}

func TestConstant_Int64(t *testing.T) {
	testValue := int64(9223372036854775807)
	c := NewConstantFromInt64(testValue)
	res, _ := c.Int64()
	assert.Equal(t, testValue, res)
}

func TestConstant_Bytes(t *testing.T) {
	b := []byte{1, 1, 2, 255}
	c, _ := NewConstantFromBytes(b)
	res, _ := c.Bytes()
	assert.Equal(t, b, res)
}

func TestNewConstantFromECPointBytes(t *testing.T) {
	str, _ := hex.DecodeString("02d6b2141c21e4f337e9b065a031a6269fb5a49253094fc6243d38662eb765db00")
	_, err := NewConstantFromECPointBytes(str)
	assert.NoError(t, err)
}

func TestNewConstantFromBox(t *testing.T) {
	json := `{
              "boxId": "e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e",
              "value": 67500000000,
              "ergoTree": "100204a00b08cd021dde34603426402615658f1d970cfa7c7bd92ac81a8b16eeebff264d59ce4604ea02d192a39a8cc7a70173007301",
              "assets": [],
              "creationHeight": 284761,
              "additionalRegisters": {},
              "transactionId": "9148408c04c2e38a6402a7950d6157730fa7d49e9ab3b9cadec481d7769918e9",
              "index": 1
            }`
	testBox, _ := NewBoxFromJson(json)
	c := NewConstantFromBox(testBox)
	encoded, _ := c.Base16()
	decoded, _ := NewConstant(encoded)
	assert.Equal(t, c, decoded)
}
