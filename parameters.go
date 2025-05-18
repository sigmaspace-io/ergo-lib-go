package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Parameters struct {
	p C.ParametersPtr
}

func newParameters(p *Parameters) *Parameters {
	runtime.AddCleanup(p, finalizeParameters, p.p)
	return p
}

// DefaultParameters returns default blockchain Parameters that were set at genesis
func DefaultParameters() *Parameters {
	var p C.ParametersPtr
	C.ergo_lib_parameters_default(&p)
	pa := &Parameters{p: p}
	return newParameters(pa)
}

// NewParameters creates new Parameters from provided blockchain Parameters
func NewParameters(
	blockVersion int32,
	storageFeeFactor int32,
	minValuePerByte int32,
	maxBlockSize int32,
	maxBlockCost int32,
	tokenAccessCost int32,
	inputCost int32,
	dataInputCost int32,
	outputCost int32) *Parameters {
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
	pa := &Parameters{p: p}
	return newParameters(pa)
}

// NewParametersFromJson parses Parameters from JSON. Support Ergo Node API/Explorer API
func NewParametersFromJson(json string) (*Parameters, error) {
	parametersJsonStr := C.CString(json)
	defer C.free(unsafe.Pointer(parametersJsonStr))

	var p C.ParametersPtr

	errPtr := C.ergo_lib_parameters_from_json(parametersJsonStr, &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	pa := &Parameters{p: p}
	return newParameters(pa), nil
}

func (p *Parameters) pointer() C.ParametersPtr {
	return p.p
}

func finalizeParameters(p C.ParametersPtr) {
	C.ergo_lib_parameters_delete(p)
}
