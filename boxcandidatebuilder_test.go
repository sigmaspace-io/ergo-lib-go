package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoxCandidateBuilder(t *testing.T) {
	addr, _ := NewAddress("3WvsT2Gm4EpsM9Pg18PdY6XyhNNMqXDsvJTbbf6ihLvAmSb7u5RN")
	contr, _ := NewContractPayToAddress(addr)
	bxVal, _ := NewBoxValue(10000000)

	boxBuilder := NewBoxCandidateBuilder(bxVal, contr, 0)

	con, _ := NewConstant("0402")
	con2, _ := NewConstant("100102")

	boxBuilder.SetRegisterValue(R4, con)
	boxBuilder.SetRegisterValue(R5, con2)

	conRes, _ := boxBuilder.RegisterValue(R4)
	assert.Equal(t, con, conRes)

	conRes2, _ := boxBuilder.RegisterValue(R5)
	assert.Equal(t, con2, conRes2)

	boxBuilder.DeleteRegisterValue(R5)
	conRes3, _ := boxBuilder.RegisterValue(R5)
	assert.Nil(t, conRes3)

	boxCand, buildErr := boxBuilder.Build()
	assert.Nil(t, buildErr)

	conRes4, _ := boxCand.RegisterValue(R4)
	assert.Equal(t, con, conRes4)
}
