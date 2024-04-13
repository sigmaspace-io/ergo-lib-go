package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type NipopowProof interface {
	// IsBetterThan implementation of the ≥ algorithm from KMZ17, see Algorithm 4
	// https://fc20.ifca.ai/preproceedings/74.pdf
	IsBetterThan(otherProof NipopowProof) (bool, error)
	// SuffixHead returns suffix head
	SuffixHead() PoPowHeader
	// Json returns json representation of NipopowProof as text
	Json() (string, error)
	pointer() C.NipopowProofPtr
}

type nipopowProof struct {
	p C.NipopowProofPtr
}

func newNipopowProof(p *nipopowProof) NipopowProof {
	runtime.SetFinalizer(p, finalizeNipopowProof)
	return p
}

// NewNipopowProof parse NipopowProof from JSON
func NewNipopowProof(json string) (NipopowProof, error) {
	nipopowProofJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(nipopowProofJsonStr))

	var p C.NipopowProofPtr

	errPtr := C.ergo_lib_nipopow_proof_from_json(nipopowProofJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	n := &nipopowProof{p: p}

	return newNipopowProof(n), nil
}

func (p *nipopowProof) IsBetterThan(otherProof NipopowProof) (bool, error) {
	res := C.ergo_lib_nipopow_proof_is_better_than(p.p, otherProof.pointer())
	err := newError(res.error)
	if err.isError() {
		return false, err.error()
	}
	return bool(res.value), nil
}

func (p *nipopowProof) SuffixHead() PoPowHeader {
	var ptr C.PoPowHeaderPtr
	C.ergo_lib_nipopow_proof_suffix_head(p.p, &ptr)
	pp := &poPowHeader{p: ptr}
	return newPoPowHeader(pp)
}

func (p *nipopowProof) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_nipopow_proof_to_json(p.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (p *nipopowProof) pointer() C.NipopowProofPtr {
	return p.p
}

func finalizeNipopowProof(p *nipopowProof) {
	C.ergo_lib_nipopow_proof_delete(p.p)
}

// NipopowVerifier a verifier for PoPow proofs. During its lifetime, it processes many proofs with the aim of
// deducing at any given point what is the best (sub)chain rooted at the specified genesis
type NipopowVerifier interface {
	// BestProof returns the best NipopowProof
	BestProof() NipopowProof
	// BestChain returns chain of BlockHeaders from the best proof
	BestChain() BlockHeaders
	// Process given NipopowProof
	Process(newProof NipopowProof) error
}

type nipopowVerifier struct {
	p C.NipopowVerifierPtr
}

func newNipopowVerifier(n *nipopowVerifier) NipopowVerifier {
	runtime.SetFinalizer(n, finalizeNipopowVerifier)
	return n
}

// NewNipopowVerifier creates a new NipopowVerifier
func NewNipopowVerifier(genesisBlockId BlockId) NipopowVerifier {
	var p C.NipopowVerifierPtr
	C.ergo_lib_nipopow_verifier_new(genesisBlockId.pointer(), &p)
	np := &nipopowVerifier{p: p}
	return newNipopowVerifier(np)
}

func (n *nipopowVerifier) BestProof() NipopowProof {
	var p C.NipopowProofPtr
	C.ergo_lib_nipopow_verifier_best_proof(n.p, &p)
	np := &nipopowProof{p: p}
	return newNipopowProof(np)
}

func (n *nipopowVerifier) BestChain() BlockHeaders {
	var p C.BlockHeadersPtr
	C.ergo_lib_nipopow_verifier_best_chain(n.p, &p)
	bh := &blockHeaders{p: p}
	return newBlockHeaders(bh)
}

func (n *nipopowVerifier) Process(newProof NipopowProof) error {
	errPtr := C.ergo_lib_nipopow_verifier_process(n.p, newProof.pointer())
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func finalizeNipopowVerifier(n *nipopowVerifier) {
	C.ergo_lib_nipopow_verifier_delete(n.p)
}

type PoPowHeader interface {
	// Header returns BlockHeader of PoPowHeader
	Header() (BlockHeader, error)
	// Interlinks returns BlockIds of PoPowHeader
	Interlinks() (BlockIds, error)
	// InterlinksProof returns BatchMerkleProof of PoPowHeader
	InterlinksProof() (BatchMerkleProof, error)
	// CheckInterlinksProof checks interlinks proof
	CheckInterlinksProof() bool
	// Json returns json representation of PoPowHeader as string
	Json() (string, error)
}

type poPowHeader struct {
	p C.PoPowHeaderPtr
}

func newPoPowHeader(p *poPowHeader) PoPowHeader {
	runtime.SetFinalizer(p, finalizePoPowHeader)
	return p
}

// NewPoPowHeader parses PoPowHeader from json string
func NewPoPowHeader(json string) (PoPowHeader, error) {
	poPowHeaderJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(poPowHeaderJsonStr))

	var p C.PoPowHeaderPtr

	errPtr := C.ergo_lib_popow_header_from_json(poPowHeaderJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	n := &poPowHeader{p: p}

	return newPoPowHeader(n), nil
}

func (p *poPowHeader) Header() (BlockHeader, error) {
	var ptr C.BlockHeaderPtr
	errPtr := C.ergo_lib_popow_header_get_header(p.p, &ptr)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bh := &blockHeader{p: ptr}
	return newBlockHeader(bh), nil
}

func (p *poPowHeader) Interlinks() (BlockIds, error) {
	var ptr C.BlockIdsPtr
	errPtr := C.ergo_lib_popow_header_get_interlinks(p.p, &ptr)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bi := &blockIds{p: ptr}
	return newBlockIds(bi), nil
}

func (p *poPowHeader) InterlinksProof() (BatchMerkleProof, error) {
	var ptr C.BatchMerkleProofPtr
	errPtr := C.ergo_lib_popow_header_get_interlinks_proof(p.p, &ptr)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	b := &batchMerkleProof{p: ptr}
	return newBatchMerkleProof(b), nil
}

func (p *poPowHeader) CheckInterlinksProof() bool {
	res := C.ergo_lib_popow_header_check_interlinks_proof(p.p)
	return bool(res)
}

func (p *poPowHeader) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_popow_header_to_json(p.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func finalizePoPowHeader(p *poPowHeader) {
	C.ergo_lib_popow_header_delete(p.p)
}
