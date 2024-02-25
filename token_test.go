package ergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTokenIdFromString(t *testing.T) {
	tokenId, _ := NewTokenId("19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac")
	assert.NotNil(t, tokenId)
}

func TestTokenId_Base16(t *testing.T) {
	tokenId, _ := NewTokenId("19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac")
	tokenIdStr := tokenId.Base16()
	assert.Equal(t, "19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac", tokenIdStr)
}

func TestTokenAmount(t *testing.T) {
	amount := int64(12345678)
	tokenAmount, _ := NewTokenAmount(amount)
	assert.Equal(t, amount, tokenAmount.Int64())
}

func TestNewToken(t *testing.T) {
	tokenIdStr := "19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac"
	tokenId, _ := NewTokenId(tokenIdStr)
	tokenAmountNr := int64(12345678)
	tokenAmount, _ := NewTokenAmount(tokenAmountNr)

	token := NewToken(tokenId, tokenAmount)

	newTokenId := token.Id()
	newTokenAmount := token.Amount()

	assert.Equal(t, tokenIdStr, newTokenId.Base16())
	assert.Equal(t, tokenAmountNr, newTokenAmount.Int64())

}

func TestToken_JsonEIP12(t *testing.T) {
	tokenJsonStr := `{"tokenId":"19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac","amount":"12345678"}`

	tokenId, _ := NewTokenId("19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac")
	tokenAmount, _ := NewTokenAmount(12345678)

	token := NewToken(tokenId, tokenAmount)

	tokenJson, _ := token.JsonEIP12()

	assert.Equal(t, tokenJsonStr, tokenJson)
}

func TestNewTokens(t *testing.T) {
	tokens := NewTokens()
	assert.Equal(t, 0, tokens.Len())

	tokenIdStr := "19475d9a78377ff0f36e9826cec439727bea522f6ffa3bda32e20d2f8b3103ac"
	tokenId, _ := NewTokenId(tokenIdStr)
	tokenAmountNr := int64(12345678)
	tokenAmount, _ := NewTokenAmount(tokenAmountNr)

	token := NewToken(tokenId, tokenAmount)

	tokens.Add(token)
	assert.Equal(t, 1, tokens.Len())

	newToken, _ := tokens.Get(0)

	assert.Equal(t, tokenIdStr, newToken.Id().Base16())
}
