package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type NipopowProof struct {
	p C.NipopowProofPtr
}

func newNipopowProof(p *NipopowProof) *NipopowProof {
	runtime.AddCleanup(p, finalizeNipopowProof, p.p)
	return p
}

// NewNipopowProof parse NipopowProof from JSON
func NewNipopowProof(json string) (*NipopowProof, error) {
	nipopowProofJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(nipopowProofJsonStr))

	var p C.NipopowProofPtr

	errPtr := C.ergo_lib_nipopow_proof_from_json(nipopowProofJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	n := &NipopowProof{p: p}

	return newNipopowProof(n), nil
}

// IsBetterThan implementation of the ≥ algorithm from KMZ17, see Algorithm 4
// https://fc20.ifca.ai/preproceedings/74.pdf
func (p *NipopowProof) IsBetterThan(otherProof *NipopowProof) (bool, error) {
	res := C.ergo_lib_nipopow_proof_is_better_than(p.p, otherProof.pointer())
	runtime.KeepAlive(p)
	runtime.KeepAlive(otherProof)
	err := newError(res.error)
	if err.isError() {
		return false, err.error()
	}
	return bool(res.value), nil
}

// SuffixHead returns suffix head
func (p *NipopowProof) SuffixHead() *PoPowHeader {
	var ptr C.PoPowHeaderPtr
	C.ergo_lib_nipopow_proof_suffix_head(p.p, &ptr)
	runtime.KeepAlive(p)
	pp := &PoPowHeader{p: ptr}
	return newPoPowHeader(pp)
}

// Json returns json representation of NipopowProof as text
func (p *NipopowProof) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_nipopow_proof_to_json(p.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(p)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

func (p *NipopowProof) pointer() C.NipopowProofPtr {
	return p.p
}

func finalizeNipopowProof(p C.NipopowProofPtr) {
	C.ergo_lib_nipopow_proof_delete(p)
}

// NipopowVerifier a verifier for PoPow proofs. During its lifetime, it processes many proofs with the aim of
// deducing at any given point what is the best (sub)chain rooted at the specified genesis
type NipopowVerifier struct {
	p C.NipopowVerifierPtr
}

func newNipopowVerifier(n *NipopowVerifier) *NipopowVerifier {
	runtime.AddCleanup(n, finalizeNipopowVerifier, n.p)
	return n
}

// NewNipopowVerifier creates a new NipopowVerifier
func NewNipopowVerifier(genesisBlockId *BlockId) *NipopowVerifier {
	var p C.NipopowVerifierPtr
	C.ergo_lib_nipopow_verifier_new(genesisBlockId.pointer(), &p)
	runtime.KeepAlive(genesisBlockId)
	np := &NipopowVerifier{p: p}
	return newNipopowVerifier(np)
}

// BestProof returns the best NipopowProof
func (n *NipopowVerifier) BestProof() *NipopowProof {
	var p C.NipopowProofPtr
	C.ergo_lib_nipopow_verifier_best_proof(n.p, &p)
	runtime.KeepAlive(n)
	np := &NipopowProof{p: p}
	return newNipopowProof(np)
}

// BestChain returns chain of BlockHeaders from the best proof
func (n *NipopowVerifier) BestChain() *BlockHeaders {
	var p C.BlockHeadersPtr
	C.ergo_lib_nipopow_verifier_best_chain(n.p, &p)
	runtime.KeepAlive(n)
	bh := &BlockHeaders{p: p}
	return newBlockHeaders(bh)
}

// Process given NipopowProof
func (n *NipopowVerifier) Process(newProof *NipopowProof) error {
	errPtr := C.ergo_lib_nipopow_verifier_process(n.p, newProof.pointer())
	runtime.KeepAlive(n)
	runtime.KeepAlive(newProof)
	err := newError(errPtr)
	if err.isError() {
		return err.error()
	}
	return nil
}

func finalizeNipopowVerifier(p C.NipopowVerifierPtr) {
	C.ergo_lib_nipopow_verifier_delete(p)
}

type PoPowHeader struct {
	p C.PoPowHeaderPtr
}

func newPoPowHeader(p *PoPowHeader) *PoPowHeader {
	runtime.AddCleanup(p, finalizePoPowHeader, p.p)
	return p
}

// NewPoPowHeader parses PoPowHeader from json string
func NewPoPowHeader(json string) (*PoPowHeader, error) {
	poPowHeaderJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(poPowHeaderJsonStr))

	var p C.PoPowHeaderPtr

	errPtr := C.ergo_lib_popow_header_from_json(poPowHeaderJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}

	n := &PoPowHeader{p: p}

	return newPoPowHeader(n), nil
}

// Header returns BlockHeader of PoPowHeader
func (p *PoPowHeader) Header() (*BlockHeader, error) {
	var ptr C.BlockHeaderPtr
	errPtr := C.ergo_lib_popow_header_get_header(p.p, &ptr)
	runtime.KeepAlive(p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bh := &BlockHeader{p: ptr}
	return newBlockHeader(bh), nil
}

// Interlinks returns BlockIds of PoPowHeader
func (p *PoPowHeader) Interlinks() (*BlockIds, error) {
	var ptr C.BlockIdsPtr
	errPtr := C.ergo_lib_popow_header_get_interlinks(p.p, &ptr)
	runtime.KeepAlive(p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	bi := &BlockIds{p: ptr}
	return newBlockIds(bi), nil
}

// InterlinksProof returns BatchMerkleProof of PoPowHeader
func (p *PoPowHeader) InterlinksProof() (*BatchMerkleProof, error) {
	var ptr C.BatchMerkleProofPtr
	errPtr := C.ergo_lib_popow_header_get_interlinks_proof(p.p, &ptr)
	runtime.KeepAlive(p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	b := &BatchMerkleProof{p: ptr}
	return newBatchMerkleProof(b), nil
}

// CheckInterlinksProof checks interlinks proof
func (p *PoPowHeader) CheckInterlinksProof() bool {
	res := C.ergo_lib_popow_header_check_interlinks_proof(p.p)
	runtime.KeepAlive(p)
	return bool(res)
}

// Json returns json representation of PoPowHeader as string
func (p *PoPowHeader) Json() (string, error) {
	var outStr *C.char

	errPtr := C.ergo_lib_popow_header_to_json(p.p, &outStr)
	defer C.ergo_lib_delete_string(outStr)
	runtime.KeepAlive(p)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	result := C.GoString(outStr)

	return result, nil
}

// Equals checks if provided PoPowHeader is same
func (p *PoPowHeader) Equals(poPowHeader *PoPowHeader) bool {
	res := C.ergo_lib_po_pow_header_eq(p.p, poPowHeader.pointer())
	runtime.KeepAlive(p)
	runtime.KeepAlive(poPowHeader)
	return bool(res)
}

func (p *PoPowHeader) pointer() C.PoPowHeaderPtr {
	return p.p
}

func finalizePoPowHeader(p C.PoPowHeaderPtr) {
	C.ergo_lib_popow_header_delete(p)
}
