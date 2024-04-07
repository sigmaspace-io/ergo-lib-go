package ergo

/*
#include "ergo.h"
*/
import "C"
import "runtime"

// BoxSelection represents selected boxes with change boxes. Instance are created by SimpleBoxSelector
type BoxSelection interface {
	// Boxes returns selected boxes to spend as transaction inputs
	Boxes() Boxes
	// ChangeBoxes returns selected boxes to use as change
	ChangeBoxes() BoxAssetsDataList
	pointer() C.BoxSelectionPtr
}

type boxSelection struct {
	p C.BoxSelectionPtr
}

func newBoxSelection(b *boxSelection) BoxSelection {
	runtime.SetFinalizer(b, finalizeBoxSelection)
	return b
}

// NewBoxSelection creates a selection to easily inject custom selection algorithms
func NewBoxSelection(ergoBoxes Boxes, changeErgoBoxes BoxAssetsDataList) BoxSelection {
	var p C.BoxSelectionPtr
	C.ergo_lib_box_selection_new(ergoBoxes.pointer(), changeErgoBoxes.pointer(), &p)
	bs := &boxSelection{p: p}
	return newBoxSelection(bs)
}

func (b *boxSelection) Boxes() Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_box_selection_boxes(b.p, &p)
	bo := &boxes{p: p}
	return newBoxes(bo)
}

func (b *boxSelection) ChangeBoxes() BoxAssetsDataList {
	var p C.ErgoBoxAssetsDataListPtr
	C.ergo_lib_box_selection_change(b.p, &p)
	ba := &boxAssetsDataList{p: p}
	return newBoxAssetsDataList(ba)
}

func (b *boxSelection) pointer() C.BoxSelectionPtr {
	return b.p
}

func finalizeBoxSelection(b *boxSelection) {
	C.ergo_lib_box_selection_delete(b.p)
}

// SimpleBoxSelector is a naive box selector, collects inputs until target balance is reached
type SimpleBoxSelector interface {
	// Select selects inputs to satisfy target balance and tokens
	// Parameters:
	// inputs - available inputs (returns an error, if empty)
	// targetBalance - coins (in nanoERGs) needed
	// targetTokens - amount of tokens needed
	// Returns: selected inputs and box assets(value+tokens) with change
	Select(inputs Boxes, targetBalance BoxValue, targetTokens Tokens) (BoxSelection, error)
}

type simpleBoxSelector struct {
	p C.SimpleBoxSelectorPtr
}

func newSimpleBoxSelector(s *simpleBoxSelector) SimpleBoxSelector {
	runtime.SetFinalizer(s, finalizeSimpleBoxSelector)
	return s
}

// NewSimpleBoxSelector creates a new SimpleBoxSelector
func NewSimpleBoxSelector() SimpleBoxSelector {
	var p C.SimpleBoxSelectorPtr
	C.ergo_lib_simple_box_selector_new(&p)
	s := &simpleBoxSelector{p: p}
	return newSimpleBoxSelector(s)
}

func (b *simpleBoxSelector) Select(inputs Boxes, targetBalance BoxValue, targetTokens Tokens) (BoxSelection, error) {
	var p C.BoxSelectionPtr
	errPtr := C.ergo_lib_simple_box_selector_select(b.p, inputs.pointer(), targetBalance.pointer(), targetTokens.pointer(), &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	bs := &boxSelection{p: p}
	return newBoxSelection(bs), nil
}

func finalizeSimpleBoxSelector(s *simpleBoxSelector) {
	C.ergo_lib_simple_box_selector_delete(s.p)
}
