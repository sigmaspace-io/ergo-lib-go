package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type TokenId interface {
	Base16() string
	pointer() C.TokenIdPtr
}

type tokenId struct {
	p C.TokenIdPtr
}

func newTokenId(t *tokenId) TokenId {
	runtime.SetFinalizer(t, finalizeTokenId)
	return t
}

// NewTokenId creates a TokenId from a base16-encoded string (32 byte digest)
func NewTokenId(s string) (TokenId, error) {
	tokenIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(tokenIdStr))

	var p C.TokenIdPtr

	errPtr := C.ergo_lib_token_id_from_str(tokenIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &tokenId{p: p}

	return newTokenId(t), nil
}

func finalizeTokenId(t *tokenId) {
	C.ergo_lib_token_id_delete(t.p)
}

func (t *tokenId) Base16() string {
	var outStr *C.char

	C.ergo_lib_token_id_to_str(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)

	result := C.GoString(outStr)

	return result
}

func (t *tokenId) pointer() C.TokenIdPtr {
	return t.p
}

type TokenAmount interface {
	Int64() int64
	pointer() C.TokenAmountPtr
}

type tokenAmount struct {
	p C.TokenAmountPtr
}

func newTokenAmount(t *tokenAmount) TokenAmount {
	runtime.SetFinalizer(t, finalizeTokenAmount)
	return t
}

// NewTokenAmount creates TokenAmount from int64
func NewTokenAmount(amount int64) (TokenAmount, error) {
	var p C.TokenAmountPtr

	errPtr := C.ergo_lib_token_amount_from_i64(C.int64_t(amount), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &tokenAmount{p: p}

	return newTokenAmount(t), nil
}

// Int64 converts TokenAmount to int64
func (t *tokenAmount) Int64() int64 {
	amount := C.ergo_lib_token_amount_as_i64(t.p)
	return int64(amount)
}

func (t *tokenAmount) pointer() C.TokenAmountPtr {
	return t.p
}

func finalizeTokenAmount(t *tokenAmount) {
	C.ergo_lib_token_amount_delete(t.p)
}

type Token interface {
	Id() TokenId
	Amount() TokenAmount
	JsonEIP12() (string, error)
	pointer() C.TokenPtr
}

type token struct {
	p C.TokenPtr
}

func newToken(t *token) Token {
	runtime.SetFinalizer(t, finalizeToken)
	return t
}

func NewToken(tokenId TokenId, tokenAmount TokenAmount) Token {
	var p C.TokenPtr

	C.ergo_lib_token_new(tokenId.pointer(), tokenAmount.pointer(), &p)

	t := &token{p: p}

	return newToken(t)
}

func (t *token) Id() TokenId {
	var tokenIdPtr C.TokenIdPtr
	C.ergo_lib_token_get_id(t.p, &tokenIdPtr)

	tId := &tokenId{p: tokenIdPtr}

	return newTokenId(tId)
}

func (t *token) Amount() TokenAmount {
	var tokenAmountPtr C.TokenAmountPtr
	C.ergo_lib_token_get_amount(t.p, &tokenAmountPtr)

	tAmount := &tokenAmount{p: tokenAmountPtr}

	return newTokenAmount(tAmount)
}

func (t *token) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_token_to_json_eip12(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (t *token) pointer() C.TokenPtr {
	return t.p
}

func finalizeToken(t *token) {
	C.ergo_lib_token_delete(t.p)
}

type Tokens interface {
	Len() int
	Get(index int) (Token, error)
	Add(token Token)
	pointer() C.TokensPtr
}

type tokens struct {
	p C.TokensPtr
}

func newTokens(t *tokens) Tokens {
	runtime.SetFinalizer(t, finalizeTokens)
	return t
}

func NewTokens() Tokens {
	var p C.TokensPtr
	C.ergo_lib_tokens_new(&p)

	t := &tokens{p: p}

	return newTokens(t)
}

func (t *tokens) Len() int {
	res := C.ergo_lib_tokens_len(t.p)
	return int(res)
}

func (t *tokens) Get(index int) (Token, error) {
	var p C.TokenPtr

	res := C.ergo_lib_tokens_get(t.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	to := &token{p: p}

	return newToken(to), nil
}

func (t *tokens) Add(token Token) {
	C.ergo_lib_tokens_add(token.pointer(), t.p)
}

func (t *tokens) pointer() C.TokensPtr {
	return t.p
}

func finalizeTokens(t *tokens) {
	C.ergo_lib_tokens_delete(t.p)
}
