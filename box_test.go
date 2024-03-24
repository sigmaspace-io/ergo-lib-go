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

func TestNewBoxValue(t *testing.T) {
	testBoxValue, err := NewBoxValue(1000000000)
	assert.Nil(t, err)
	assert.Equal(t, int64(1000000000), testBoxValue.Int64())
}

func TestSumOfBoxValues(t *testing.T) {
	testBoxValue1, _ := NewBoxValue(1000000000)
	testBoxValue2, _ := NewBoxValue(3000000000)

	testBoxValueSum, _ := SumOfBoxValues(testBoxValue1, testBoxValue2)

	assert.Equal(t, int64(4000000000), testBoxValueSum.Int64())
}

func TestNewBox(t *testing.T) {
	testBoxValue, _ := NewBoxValue(67500000000)
	testBoxId, _ := NewBoxId("e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e")
	testTxId, _ := NewTxId("9148408c04c2e38a6402a7950d6157730fa7d49e9ab3b9cadec481d7769918e9")
	testTokens := NewTokens()
	testCreationHeight := uint32(284761)
	testErgoTree, _ := NewTree("100204a00b08cd021dde34603426402615658f1d970cfa7c7bd92ac81a8b16eeebff264d59ce4604ea02d192a39a8cc7a70173007301")
	testContract := NewContractFromTree(testErgoTree)

	testErgoBox, boxErr := NewBox(testBoxValue, testCreationHeight, testContract, testTxId, 1, testTokens)

	assert.Nil(t, boxErr)
	assert.Equal(t, testCreationHeight, testErgoBox.CreationHeight())
	assert.Equal(t, testBoxId, testErgoBox.BoxId())
	assert.Equal(t, testBoxValue, testErgoBox.BoxValue())
	assert.Equal(t, testErgoTree, testErgoBox.Tree())
}

func TestNewBoxFromJson(t *testing.T) {
	testBoxValue, _ := NewBoxValue(67500000000)
	testBoxId, _ := NewBoxId("e56847ed19b3dc6b72828fcfb992fdf7310828cf291221269b7ffc72fd66706e")
	testErgoTree, _ := NewTree("100204a00b08cd021dde34603426402615658f1d970cfa7c7bd92ac81a8b16eeebff264d59ce4604ea02d192a39a8cc7a70173007301")
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

	testErgoBox, boxErr := NewBoxFromJson(json)

	assert.Nil(t, boxErr)
	assert.Equal(t, uint32(284761), testErgoBox.CreationHeight())
	assert.Equal(t, testBoxId, testErgoBox.BoxId())
	assert.Equal(t, testBoxValue, testErgoBox.BoxValue())
	assert.Equal(t, testErgoTree, testErgoBox.Tree())
}
