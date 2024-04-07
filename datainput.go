package ergo

/*
#include "ergo.h"
*/
import "C"
import "runtime"

// DataInput represent inputs that are used to enrich script context, but won't be spent by the transaction
type DataInput interface {
	// BoxId returns the BoxId of the DataInput
	BoxId() BoxId
	pointer() C.DataInputPtr
}

type dataInput struct {
	p C.DataInputPtr
}

func newDataInput(d *dataInput) DataInput {
	runtime.SetFinalizer(d, finalizeDataInput)
	return d
}

// NewDataInput create DataInput from BoxId
func NewDataInput(boxId BoxId) DataInput {
	var p C.DataInputPtr
	C.ergo_lib_data_input_new(boxId.pointer(), &p)
	d := &dataInput{p: p}
	return newDataInput(d)
}

func (d *dataInput) BoxId() BoxId {
	var p C.BoxIdPtr
	C.ergo_lib_data_input_box_id(d.p, &p)
	bi := &boxId{p: p}
	return newBoxId(bi)
}

func (d *dataInput) pointer() C.DataInputPtr {
	return d.p
}

func finalizeDataInput(d *dataInput) {
	C.ergo_lib_data_input_delete(d.p)
}

// DataInputs an ordered collection if DataInput
type DataInputs interface {
	// Len returns the length of the collection
	Len() uint32
	// Get returns the Input at the provided index if it exists
	Get(index uint32) (DataInput, error)
	// Add adds provided DataInput to the end of the collection
	Add(dataInput DataInput)
}

type dataInputs struct {
	p C.DataInputsPtr
}

func newDataInputs(d *dataInputs) DataInputs {
	runtime.SetFinalizer(d, finalizeDataInputs)
	return d
}

// NewDataInputs creates an empty DataInputs collection
func NewDataInputs() DataInputs {
	var p C.DataInputsPtr
	C.ergo_lib_data_inputs_new(&p)
	d := &dataInputs{p: p}
	return newDataInputs(d)
}

func (d *dataInputs) Len() uint32 {
	res := C.ergo_lib_data_inputs_len(d.p)
	return uint32(res)
}

func (d *dataInputs) Get(index uint32) (DataInput, error) {
	var p C.DataInputPtr

	res := C.ergo_lib_data_inputs_get(d.p, C.ulong(index), &p)
	err := newError(res.error)
	if err.isError() {
		return nil, err.error()
	}

	if res.is_some {
		di := &dataInput{p: p}
		return newDataInput(di), nil
	}

	return nil, nil
}

func (d *dataInputs) Add(dataInput DataInput) {
	C.ergo_lib_data_inputs_add(dataInput.pointer(), d.p)
}

func finalizeDataInputs(d *dataInputs) {
	C.ergo_lib_data_inputs_delete(d.p)
}
