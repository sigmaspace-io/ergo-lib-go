package ergo

/*
#include "ergo.h"
*/
import "C"
import (
	"iter"
	"runtime"
)

// DataInput represent Inputs that are used to enrich script context, but won't be spent by the Transaction
type DataInput struct {
	p C.DataInputPtr
}

func newDataInput(d *DataInput) *DataInput {
	runtime.AddCleanup(d, finalizeDataInput, d.p)
	return d
}

// NewDataInput create DataInput from BoxId
func NewDataInput(boxId *BoxId) *DataInput {
	var p C.DataInputPtr
	C.ergo_lib_data_input_new(boxId.pointer(), &p)
	runtime.KeepAlive(boxId)
	d := &DataInput{p: p}
	return newDataInput(d)
}

// BoxId returns the BoxId of the DataInput
func (d *DataInput) BoxId() *BoxId {
	var p C.BoxIdPtr
	C.ergo_lib_data_input_box_id(d.p, &p)
	runtime.KeepAlive(d)
	bi := &BoxId{p: p}
	return newBoxId(bi)
}

func (d *DataInput) pointer() C.DataInputPtr {
	return d.p
}

func finalizeDataInput(p C.DataInputPtr) {
	C.ergo_lib_data_input_delete(p)
}

// DataInputs an ordered collection if DataInput
type DataInputs struct {
	p C.DataInputsPtr
}

func newDataInputs(d *DataInputs) *DataInputs {
	runtime.AddCleanup(d, finalizeDataInputs, d.p)
	return d
}

// NewDataInputs creates an empty DataInputs collection
func NewDataInputs() *DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_data_inputs_new(&p)
	d := &DataInputs{p: p}
	return newDataInputs(d)
}

// Len returns the length of the collection
func (d *DataInputs) Len() int {
	res := C.ergo_lib_data_inputs_len(d.p)
	runtime.KeepAlive(d)
	return int(res)
}

// Get returns the Input at the provided index if it exists
func (d *DataInputs) Get(index int) (*DataInput, error) {
	var p C.DataInputPtr

	res := C.ergo_lib_data_inputs_get(d.p, C.uintptr_t(index), &p)
	runtime.KeepAlive(d)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		di := &DataInput{p: p}
		return newDataInput(di), nil
	}

	return nil, nil
}

// Add adds provided DataInput to the end of the collection
func (d *DataInputs) Add(dataInput *DataInput) {
	C.ergo_lib_data_inputs_add(dataInput.pointer(), d.p)
	runtime.KeepAlive(d)
	runtime.KeepAlive(dataInput)
}

// All returns an iterator over all DataInput inside the collection
func (d *DataInputs) All() iter.Seq2[int, *DataInput] {
	return func(yield func(int, *DataInput) bool) {
		for i := 0; i < d.Len(); i++ {
			tk, err := d.Get(i)
			if err != nil {
				return
			}
			if !yield(i, tk) {
				return
			}
		}
		runtime.KeepAlive(d)
	}
}

func (d *DataInputs) pointer() C.DataInputsPtr {
	return d.p
}

func finalizeDataInputs(p C.DataInputsPtr) {
	C.ergo_lib_data_inputs_delete(p)
}
