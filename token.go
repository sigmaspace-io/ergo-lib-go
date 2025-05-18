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
type TokenId struct {
	p C.TokenIdPtr
}

func newTokenId(t *TokenId) *TokenId {
	runtime.AddCleanup(t, finalizeTokenId, t.p)
	return t
}

// NewTokenId creates a TokenId from a base16-encoded string (32 byte digest)
func NewTokenId(s string) (*TokenId, error) {
	tokenIdStr := C.CString(s)
	defer C.free(unsafe.Pointer(tokenIdStr))

	var p C.TokenIdPtr

	errPtr := C.ergo_lib_token_id_from_str(tokenIdStr, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &TokenId{p: p}

	return newTokenId(t), nil
}

// NewTokenIdFromBoxId creates a TokenId from ergo Box id (32 byte digest)
func NewTokenIdFromBoxId(boxId *BoxId) *TokenId {
	var p C.TokenIdPtr
	C.ergo_lib_token_id_from_box_id(boxId.pointer(), &p)
	runtime.KeepAlive(boxId)
	t := &TokenId{p: p}
	return newTokenId(t)
}

// Equals checks if provided TokenId is same
func (t *TokenId) Equals(tokenId *TokenId) bool {
	res := C.ergo_lib_token_id_eq(t.p, tokenId.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(tokenId)
	return bool(res)
}

func finalizeTokenId(p C.TokenIdPtr) {
	C.ergo_lib_token_id_delete(p)
}

// Base16 returns the TokenId as base16 encoded string
func (t *TokenId) Base16() string {
	var outStr *C.char

	C.ergo_lib_token_id_to_str(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(t)

	result := C.GoString(outStr)

	return result
}

func (t *TokenId) pointer() C.TokenIdPtr {
	return t.p
}

// TokenAmount is Token amount with bound checks
type TokenAmount struct {
	p C.TokenAmountPtr
}

func newTokenAmount(t *TokenAmount) *TokenAmount {
	runtime.AddCleanup(t, finalizeTokenAmount, t.p)
	return t
}

// NewTokenAmount creates TokenAmount from int64
func NewTokenAmount(amount int64) (*TokenAmount, error) {
	var p C.TokenAmountPtr

	errPtr := C.ergo_lib_token_amount_from_i64(C.int64_t(amount), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	t := &TokenAmount{p: p}

	return newTokenAmount(t), nil
}

// Int64 converts TokenAmount to int64
func (t *TokenAmount) Int64() int64 {
	amount := C.ergo_lib_token_amount_as_i64(t.p)
	return int64(amount)
}

// Equals checks if provided TokenAmount is same
func (t *TokenAmount) Equals(tokenAmount *TokenAmount) bool {
	res := C.ergo_lib_token_amount_eq(t.p, tokenAmount.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(tokenAmount)
	return bool(res)
}

func (t *TokenAmount) pointer() C.TokenAmountPtr {
	return t.p
}

func finalizeTokenAmount(p C.TokenAmountPtr) {
	C.ergo_lib_token_amount_delete(p)
}

// Token represented with TokenId paired with its TokenAmount
type Token struct {
	p C.TokenPtr
}

func newToken(t *Token) *Token {
	runtime.AddCleanup(t, finalizeToken, t.p)
	return t
}

// NewToken creates Token from provided TokenId and TokenAmount
func NewToken(tokenId *TokenId, tokenAmount *TokenAmount) *Token {
	var p C.TokenPtr

	C.ergo_lib_token_new(tokenId.pointer(), tokenAmount.pointer(), &p)
	runtime.KeepAlive(tokenId)
	runtime.KeepAlive(tokenAmount)

	t := &Token{p: p}

	return newToken(t)
}

// Id returns TokenId of the Token
func (t *Token) Id() *TokenId {
	var tokenIdPtr C.TokenIdPtr
	C.ergo_lib_token_get_id(t.p, &tokenIdPtr)

	tId := &TokenId{p: tokenIdPtr}

	return newTokenId(tId)
}

// Amount returns TokenAmount of the Token
func (t *Token) Amount() *TokenAmount {
	var tokenAmountPtr C.TokenAmountPtr
	C.ergo_lib_token_get_amount(t.p, &tokenAmountPtr)

	tAmount := &TokenAmount{p: tokenAmountPtr}

	return newTokenAmount(tAmount)
}

// JsonEIP12 returns json representation of Token as string according to EIP-12 https://github.com/ergoplatform/eips/pull/23
func (t *Token) JsonEIP12() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_token_to_json_eip12(t.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(t)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// Equals checks if provided Token is same
func (t *Token) Equals(token *Token) bool {
	res := C.ergo_lib_token_eq(t.p, token.pointer())
	runtime.KeepAlive(t)
	runtime.KeepAlive(token)
	return bool(res)
}

func (t *Token) pointer() C.TokenPtr {
	return t.p
}

func finalizeToken(p C.TokenPtr) {
	C.ergo_lib_token_delete(p)
}

// Tokens an ordered collection of Token
type Tokens struct {
	p C.TokensPtr
}

func newTokens(t *Tokens) *Tokens {
	runtime.AddCleanup(t, finalizeTokens, t.p)
	return t
}

// NewTokens creates an empty Tokens collection
func NewTokens() *Tokens {
	var p C.TokensPtr
	C.ergo_lib_tokens_new(&p)

	t := &Tokens{p: p}

	return newTokens(t)
}

// Len returns the length of the collection
func (t *Tokens) Len() int {
	res := C.ergo_lib_tokens_len(t.p)
	runtime.KeepAlive(t)
	return int(res)
}

// Get returns the Token at the provided index if it exists
func (t *Tokens) Get(index int) (*Token, error) {
	var p C.TokenPtr

	res := C.ergo_lib_tokens_get(t.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(t)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	to := &Token{p: p}

	return newToken(to), nil
}

// Add adds provided Token to the end of the collection
func (t *Tokens) Add(token *Token) {
	C.ergo_lib_tokens_add(token.pointer(), t.p)
	runtime.KeepAlive(t)
	runtime.KeepAlive(token)
}

// All returns an iterator over all Token inside the collection
func (t *Tokens) All() iter.Seq2[int, *Token] {
	return func(yield func(int, *Token) bool) {
		for i := 0; i < t.Len(); i++ {
			tk, err := t.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(t)
	}
}

func (t *Tokens) pointer() C.TokensPtr {
	return t.p
}

func finalizeTokens(p C.TokensPtr) {
	C.ergo_lib_tokens_delete(p)
}
