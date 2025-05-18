package ergo

/*
#include "ergo.h"
*/
import "C"
import "runtime"

// BoxSelection represents selected Boxes with change Boxes. Instances are created by SimpleBoxSelector
type BoxSelection struct {
	p C.BoxSelectionPtr
}

func newBoxSelection(b *BoxSelection) *BoxSelection {
	runtime.AddCleanup(b, finalizeBoxSelection, b.p)
	return b
}

// NewBoxSelection creates a selection to easily inject custom selection algorithms
func NewBoxSelection(ergoBoxes *Boxes, changeErgoBoxes *BoxAssetsDataList) *BoxSelection {
	var p C.BoxSelectionPtr
	C.ergo_lib_box_selection_new(ergoBoxes.pointer(), changeErgoBoxes.pointer(), &p)
	runtime.KeepAlive(ergoBoxes)
	runtime.KeepAlive(changeErgoBoxes)
	bs := &BoxSelection{p: p}
	return newBoxSelection(bs)
}

// Boxes returns selected Boxes to spend as Transaction Inputs
func (b *BoxSelection) Boxes() *Boxes {
	var p C.ErgoBoxesPtr
	C.ergo_lib_box_selection_boxes(b.p, &p)
	runtime.KeepAlive(b)
	bo := &Boxes{p: p}
	return newBoxes(bo)
}

// ChangeBoxes returns selected Boxes to use as change
func (b *BoxSelection) ChangeBoxes() *BoxAssetsDataList {
	var p C.ErgoBoxAssetsDataListPtr
	C.ergo_lib_box_selection_change(b.p, &p)
	runtime.KeepAlive(b)
	ba := &BoxAssetsDataList{p: p}
	return newBoxAssetsDataList(ba)
}

// Equals checks if provided BoxSelection is same
func (b *BoxSelection) Equals(boxSelection *BoxSelection) bool {
	res := C.ergo_lib_box_selection_eq(b.p, boxSelection.pointer())
	runtime.KeepAlive(b)
	runtime.KeepAlive(boxSelection)
	return bool(res)
}

func (b *BoxSelection) pointer() C.BoxSelectionPtr {
	return b.p
}

func finalizeBoxSelection(p C.BoxSelectionPtr) {
	C.ergo_lib_box_selection_delete(p)
}

// SimpleBoxSelector is a naive Box selector, collects Inputs until target balance is reached
type SimpleBoxSelector struct {
	p C.SimpleBoxSelectorPtr
}

func newSimpleBoxSelector(s *SimpleBoxSelector) *SimpleBoxSelector {
	runtime.AddCleanup(s, finalizeSimpleBoxSelector, s.p)
	return s
}

// NewSimpleBoxSelector creates a new SimpleBoxSelector
func NewSimpleBoxSelector() *SimpleBoxSelector {
	var p C.SimpleBoxSelectorPtr
	C.ergo_lib_simple_box_selector_new(&p)
	s := &SimpleBoxSelector{p: p}
	return newSimpleBoxSelector(s)
}

// Select selects Inputs to satisfy target balance and Tokens
// Parameters:
// Inputs - available Inputs (returns an error, if empty)
// targetBalance - coins (in nanoERGs) needed
// targetTokens - amount of Tokens needed
// Returns: selected Inputs and Box assets(value+Tokens) with change
func (b *SimpleBoxSelector) Select(inputs *Boxes, targetBalance *BoxValue, targetTokens *Tokens) (*BoxSelection, error) {
	var p C.BoxSelectionPtr
	errPtr := C.ergo_lib_simple_box_selector_select(b.p, inputs.pointer(), targetBalance.pointer(), targetTokens.pointer(), &p)
	runtime.KeepAlive(b)
	runtime.KeepAlive(inputs)
	runtime.KeepAlive(targetBalance)
	runtime.KeepAlive(targetTokens)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	bs := &BoxSelection{p: p}
	return newBoxSelection(bs), nil
}

func finalizeSimpleBoxSelector(p C.SimpleBoxSelectorPtr) {
	C.ergo_lib_simple_box_selector_delete(p)
}
