package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContractPayToAddress(t *testing.T) {
	testAddress, _ := NewAddress("9hdxkYakTHWXR992umPcvh8bAEGG9Sdoi7uW8TKXk1enXCDFBVJ")
	testContract, _ := NewContractPayToAddress(testAddress)

	testTree, _ := testContract.Tree().Address()

	assert.Equal(t, testAddress, testTree)
}

func TestNewContractFromTree(t *testing.T) {
	testTree, _ := NewTree("0008cd039ac1cb7dd70a151a6c42948c9f4f405c92f70dbfbba2d895b48e4ae7746ba6a6")
	testContract := NewContractFromTree(testTree)

	assert.Equal(t, testTree, testContract.Tree())
}
