package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Parameters interface {
	pointer() C.ParametersPtr
}

type parameters struct {
	p C.ParametersPtr
}

func newParameters(p *parameters) Parameters {
	runtime.SetFinalizer(p, finalizeParameters)
	return p
}

// DefaultParameters returns default blockchain parameters that were set at genesis
func DefaultParameters() Parameters {
	var p C.ParametersPtr
	C.ergo_lib_parameters_default(&p)
	pa := &parameters{p: p}
	return newParameters(pa)
}

// NewParameters creates new Parameters from provided blockchain parameters
func NewParameters(
	blockVersion int32,
	storageFeeFactor int32,
	minValuePerByte int32,
	maxBlockSize int32,
	maxBlockCost int32,
	tokenAccessCost int32,
	inputCost int32,
	dataInputCost int32,
	outputCost int32) Parameters {
	var p C.ParametersPtr
	C.ergo_lib_parameters_new(
		C.int32_t(blockVersion),
		C.int32_t(storageFeeFactor),
		C.int32_t(minValuePerByte),
		C.int32_t(maxBlockSize),
		C.int32_t(maxBlockCost),
		C.int32_t(tokenAccessCost),
		C.int32_t(inputCost),
		C.int32_t(dataInputCost),
		C.int32_t(outputCost),
		&p)
	pa := &parameters{p: p}
	return newParameters(pa)
}

// NewParametersFromJson parses parameters from JSON. Support Ergo Node API/Explorer API
func NewParametersFromJson(json string) (Parameters, error) {
	parametersJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(parametersJsonStr))

	var p C.ParametersPtr

	errPtr := C.ergo_lib_parameters_from_json(parametersJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	pa := &parameters{p: p}
	return newParameters(pa), nil
}

func (p *parameters) pointer() C.ParametersPtr {
	return p.p
}

func finalizeParameters(p *parameters) {
	C.ergo_lib_parameters_delete(p.p)
}
