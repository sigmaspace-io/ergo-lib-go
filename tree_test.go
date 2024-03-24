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

	assert.Equal(t, "0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041", s)
}

func TestTree_Address(t *testing.T) {
	tree, _ := NewTree("0008cd0336100ef59ced80ba5f89c4178ebd57b6c1dd0f3d135ee1db9f62fc634d637041")
	a, _ := tree.Address()

	assert.Equal(t, "9gscej8Kyzvy7AE3DMhBVpkU2CEAAM9fC6zs5dVwjAmcszPCjEr", a.Base58(MainnetPrefix))
}

func TestTree_ErgoTreeTemplateHash(t *testing.T) {
	tree1, _ := NewTree("10060e2002d1541415c323527f19ef5b103eb33c220ea8b66fcb711806b0037d115d63f204000402040004040e201a6a8c16e4b1cc9d73d03183565cfb8e79dd84198cb66beeed7d3463e0da2b98d803d601e4c6a70507d602d901026393cbc27202e4c6a7070ed6037300ea02eb02cd7201cedb6a01dde4c6a70407e4c6a706077201d1ececedda720201b2a573010093cbc2b2a47302007203edda720201b2a473030093cbc2b2a47304007203afa5d9010463afdb63087204d901064d0e948c7206017305")
	length1, _ := tree1.TemplateBytesLength()
	hash1, _ := tree1.TemplateHash()
	tree2, _ := NewTree("100204a00b08cd03c2b58229e1306e34a3f29ae74206bf44ee91e49d35ef0ff476951fab8593323eea02d192a39a8cc7a70173007301")
	length2, _ := tree2.TemplateBytesLength()
	hash2, _ := tree2.TemplateHash()

	assert.Equal(t, 120, length1)
	assert.Equal(t, "b05e92b987d2a8d8cbf875d4e0bad58c4bfc36d9e2263a4527afaa24f9d7f84e", hash1)
	assert.Equal(t, 14, length2)
	assert.Equal(t, "961e872f7ab750cb77ad75ea8a32d0ea3472bd0c230de09329b802801b3d1817", hash2)
}

func TestTree_ErgoTreeConstantsLength(t *testing.T) {
	tree, _ := NewTree("101004020e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a7017300730110010204020404040004c0fd4f05808c82f5f6030580b8c9e5ae040580f882ad16040204c0944004c0f407040004000580f882ad16d19683030191a38cc7a7019683020193c2b2a57300007473017302830108cdeeac93a38cc7b2a573030001978302019683040193b1a5730493c2a7c2b2a573050093958fa3730673079973089c73097e9a730a9d99a3730b730c0599c1a7c1b2a5730d00938cc7b2a5730e0001a390c1a7730f")

	length, _ := tree.ConstantsLength()

	assert.Equal(t, 16, length)
}

func TestTree_ErgoTreeGetConstant(t *testing.T) {
	tree, _ := NewTree("101004020e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a7017300730110010204020404040004c0fd4f05808c82f5f6030580b8c9e5ae040580f882ad16040204c0944004c0f407040004000580f882ad16d19683030191a38cc7a7019683020193c2b2a57300007473017302830108cdeeac93a38cc7b2a573030001978302019683040193b1a5730493c2a7c2b2a573050093958fa3730673079973089c73097e9a730a9d99a3730b730c0599c1a7c1b2a5730d00938cc7b2a5730e0001a390c1a7730f")

	cons1, _ := tree.Constant(0)
	consType1, _ := cons1.Type()

	cons8, _ := tree.Constant(8)
	consType8, _ := cons8.Type()

	cons16, _ := tree.Constant(15)
	consType16, _ := cons16.Type()

	assert.Equal(t, "SInt", consType1)
	assert.Equal(t, "SLong", consType8)
	assert.Equal(t, "SLong", consType16)
}

func TestTree_ErgoTreeGetConstants(t *testing.T) {
	tree, _ := NewTree("101004020e36100204a00b08cd0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798ea02d192a39a8cc7a7017300730110010204020404040004c0fd4f05808c82f5f6030580b8c9e5ae040580f882ad16040204c0944004c0f407040004000580f882ad16d19683030191a38cc7a7019683020193c2b2a57300007473017302830108cdeeac93a38cc7b2a573030001978302019683040193b1a5730493c2a7c2b2a573050093958fa3730673079973089c73097e9a730a9d99a3730b730c0599c1a7c1b2a5730d00938cc7b2a5730e0001a390c1a7730f")

	constants, _ := tree.Constants()
	assert.Equal(t, 16, len(constants))

	consType3, _ := constants[2].Type()
	consType14, _ := constants[13].Type()

	assert.Equal(t, "SColl(SInt)", consType3)
	assert.Equal(t, "SInt", consType14)
}
