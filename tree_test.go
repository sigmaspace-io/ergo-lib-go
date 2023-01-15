package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree, _ := NewTree("0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041")

	assert.NotNil(t, tree)
}

func TestNewTree_Invalid(t *testing.T) {
	_, err := NewTree("1111108zzxczbkkk")

	assert.Error(t, err)
}

func TestTree_Base16(t *testing.T) {
	tree, _ := NewTree("0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041")
	s, _ := tree.Base16()

	assert.Equal(t, "0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041", *s)
}

func TestTree_Address(t *testing.T) {
	tree, _ := NewTree("0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041")
	a, _ := tree.Address()

	assert.Equal(t, "9gscej8Kyzvy7AE3DMhBVpkU2CEAAM9fC6zs5dVwjAmcszPCjEr", a.Base58(MainnetPrefix))
}

func TestTree_ErgoTreeTemplateHash(t *testing.T) {
	tree1, _ := NewTree("10060e2002d1541415c323527f19ef5b103eb33c220ea8b66fcb711806b0037d115d63f204000402040004040e201a6a8c16e4b1cc9d73d03183565cfb8e79dd84198cb66beeed7d3463e0da2b98d803d601e4c6a70507d602d901026393cbc27202e4c6a7070ed6037300ea02eb02cd7201cedb6a01dde4c6a70407e4c6a706077201d1ececedda720201b2a573010093cbc2b2a47302007203edda720201b2a473030093cbc2b2a47304007203afa5d9010463afdb63087204d901064d0e948c7206017305")
	length1, _ := tree1.ErgoTreeTemplateBytesLength()
	hash1, _ := tree1.ErgoTreeTemplateHash()
	tree2, _ := NewTree("100204a00b08cd03c2b58229e1306e34a3f29ae74206bf44ee91e49d35ef0ff476951fab8593323eea02d192a39a8cc7a70173007301")
	length2, _ := tree2.ErgoTreeTemplateBytesLength()
	hash2, _ := tree2.ErgoTreeTemplateHash()

	assert.Equal(t, 120, length1)
	assert.Equal(t, "b05e92b987d2a8d8cbf875d4e0bad58c4bfc36d9e2263a4527afaa24f9d7f84e", hash1)
	assert.Equal(t, 14, length2)
	assert.Equal(t, "961e872f7ab750cb77ad75ea8a32d0ea3472bd0c230de09329b802801b3d1817", hash2)
}
