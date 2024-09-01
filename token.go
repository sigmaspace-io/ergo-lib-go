package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"iter"
	"runtime"
	"unsafe"
)

// TokenId (32-byte digest)
type TokenId interface {
	// Base16 returns the TokenId as base16 encoded string
	Base16() string
	// Equals checks if provided TokenId is same
	Equals(tokenId TokenId) bool
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

// NewTokenIdFromBoxId creates a TokenId from ergo box id (32 byte digest)
func NewTokenIdFromBoxId(boxId BoxId) TokenId {
	var p C.TokenIdPtr
	C.ergo_lib_token_id_from_box_id(boxId.pointer(), &p)
	t := &tokenId{p: p}
	return newTokenId(t)
}

func (t *tokenId) Equals(tokenId TokenId) bool {
	res := C.ergo_lib_token_id_eq(t.p, tokenId.pointer())
	return bool(res)
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

// TokenAmount is token amount with bound checks
type TokenAmount interface {
	// Int64 converts TokenAmount to int64
	Int64() int64
	// Equals checks if provided TokenAmount is same
	Equals(tokenAmount TokenAmount) bool
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

func (t *tokenAmount) Int64() int64 {
	amount := C.ergo_lib_token_amount_as_i64(t.p)
	return int64(amount)
}

func (t *tokenAmount) Equals(tokenAmount TokenAmount) bool {
	res := C.ergo_lib_token_amount_eq(t.p, tokenAmount.pointer())
	return bool(res)
}

func (t *tokenAmount) pointer() C.TokenAmountPtr {
	return t.p
}

func finalizeTokenAmount(t *tokenAmount) {
	C.ergo_lib_token_amount_delete(t.p)
}

// Token represented with TokenId paired with its TokenAmount
type Token interface {
	// Id returns TokenId of the Token
	Id() TokenId
	// Amount returns TokenAmount of the Token
	Amount() TokenAmount
	// JsonEIP12 returns json representation of Token as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
	JsonEIP12() (string, error)
	// Equals checks if provided Token is same
	Equals(token Token) bool
	pointer() C.TokenPtr
}

type token struct {
	p C.TokenPtr
}

func newToken(t *token) Token {
	runtime.SetFinalizer(t, finalizeToken)
	return t
}

// NewToken creates Token from provided TokenId and TokenAmount
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

func (t *token) Equals(token Token) bool {
	res := C.ergo_lib_token_eq(t.p, token.pointer())
	return bool(res)
}

func (t *token) pointer() C.TokenPtr {
	return t.p
}

func finalizeToken(t *token) {
	C.ergo_lib_token_delete(t.p)
}

// Tokens an ordered collection of Token
type Tokens interface {
	// Len returns the length of the collection
	Len() int
	// Get returns the Token at the provided index if it exists
	Get(index int) (Token, error)
	// Add adds provided Token to the end of the collection
	Add(token Token)
	// All returns an iterator over all Token inside the collection
	All() iter.Seq2[int, Token]
	pointer() C.TokensPtr
}

type tokens struct {
	p C.TokensPtr
}

func newTokens(t *tokens) Tokens {
	runtime.SetFinalizer(t, finalizeTokens)
	return t
}

// NewTokens creates an empty Tokens collection
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

	res := C.ergo_lib_tokens_get(t.p, C.uintptr_t(index), &p)
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

func (t *tokens) All() iter.Seq2[int, Token] {
	return func(yield func(int, Token) bool) {
		for i := 0; i < t.Len(); i++ {
			tk, err := t.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
	}
}

func (t *tokens) pointer() C.TokensPtr {
	return t.p
}

func finalizeTokens(t *tokens) {
	C.ergo_lib_tokens_delete(t.p)
}
