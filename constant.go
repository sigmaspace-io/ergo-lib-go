package ergo

/*
   #include "ergo.h"
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// Constant represents Ergo constant(evaluated) values
type Constant interface {
	// Base16 encode as Base16-encoded ErgoTree serialized value or throw an error if serialization failed
	Base16() (string, error)
	// Type returns the Constant type as string
	Type() (string, error)
	// Value returns the Constant value as string
	Value() (string, error)
	// Int16 extracts int16 value and returns error if wrong Constant type
	Int16() (int16, error)
	// Int32 extracts int32 value and returns error if wrong Constant type
	Int32() (int32, error)
	// Int64 extracts int64 value and returns error if wrong Constant type
	Int64() (int64, error)
	// Bytes extracts byte array and returns error if wrong Constant type
	Bytes() ([]byte, error)
	bytesLength() (int, error)
	pointer() C.ConstantPtr
}

type constant struct {
	p C.ConstantPtr
}

func newConstant(c *constant) Constant {
	runtime.SetFinalizer(c, finalizeConstant)
	return c
}

// NewConstant creates a new Constant from Base16-encoded ErgoTree serialized value
func NewConstant(s string) (Constant, error) {
	base16ErgoTree := C.CString(s)
	defer C.free(unsafe.Pointer(base16ErgoTree))

	var p C.ConstantPtr

	errPtr := C.ergo_lib_constant_from_base16(base16ErgoTree, &p)
	err := newError(errPtr)

	if err.isError() {
		return nil, err.error()
	}

	c := &constant{p}

	return newConstant(c), nil
}

// NewConstantFromInt16 creates a new Constant from int16 value
func NewConstantFromInt16(i int16) Constant {
	var p C.ConstantPtr
	C.ergo_lib_constant_from_i16(C.int16_t(i), &p)
	c := &constant{p}
	return newConstant(c)
}

// NewConstantFromInt32 creates a new Constant from int32 value
func NewConstantFromInt32(i int32) Constant {
	var p C.ConstantPtr
	C.ergo_lib_constant_from_i32(C.int32_t(i), &p)
	c := &constant{p}
	return newConstant(c)
}

// NewConstantFromInt64 creates a new Constant from int64 value
func NewConstantFromInt64(i int64) Constant {
	var p C.ConstantPtr
	C.ergo_lib_constant_from_i64(C.int64_t(i), &p)
	c := &constant{p}
	return newConstant(c)
}

// NewConstantFromBytes creates a new Constant from byte array
func NewConstantFromBytes(b []byte) (Constant, error) {
	byteData := C.CBytes(b)
	defer C.free(unsafe.Pointer(byteData))
	var p C.ConstantPtr
	errPtr := C.ergo_lib_constant_from_bytes((*C.uchar)(byteData), C.uintptr_t(len(b)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	c := &constant{p}
	return newConstant(c), nil
}

// NewConstantFromECPointBytes parse from raw EcPoint value from bytes and make ProveDlog Constant
func NewConstantFromECPointBytes(b []byte) (Constant, error) {
	byteData := C.CBytes(b)
	defer C.free(unsafe.Pointer(byteData))
	var p C.ConstantPtr
	errPtr := C.ergo_lib_constant_from_ecpoint_bytes((*C.uchar)(byteData), C.uintptr_t(len(b)), &p)
	err := newError(errPtr)
	if err.isError() {
		return nil, err.error()
	}
	c := &constant{p}
	return newConstant(c), nil
}

// NewConstantFromBox creates a new Constant from Box
func NewConstantFromBox(box Box) Constant {
	var p C.ConstantPtr
	C.ergo_lib_constant_from_ergo_box(box.pointer(), &p)
	c := &constant{p}
	return newConstant(c)
}

func (c *constant) Base16() (string, error) {
	var constantStr *C.char

	errPtr := C.ergo_lib_constant_to_base16(c.p, &constantStr)
	defer C.ergo_lib_delete_string(constantStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantStr), nil
}

func (c *constant) Type() (string, error) {
	var constantTypeStr *C.char

	errPtr := C.ergo_lib_constant_type_to_dbg_str(c.p, &constantTypeStr)
	defer C.ergo_lib_delete_string(constantTypeStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return C.GoString(constantTypeStr), nil
}

func (c *constant) Value() (string, error) {
	var constantValueStr *C.char

	errPtr := C.ergo_lib_constant_value_to_dbg_str(c.p, &constantValueStr)
	defer C.ergo_lib_delete_string(constantValueStr)
	err := newError(errPtr)

	if err.isError() {
		return "", err.error()
	}

	return strings.ReplaceAll(C.GoString(constantValueStr), " ", ""), nil
}

func (c *constant) Int16() (int16, error) {
	res := C.ergo_lib_constant_to_i16(c.p)
	err := newError(res.error)
	if err.isError() {
		return 0, err.error()
	}
	return int16(res.value), nil
}

func (c *constant) Int32() (int32, error) {
	res := C.ergo_lib_constant_to_i32(c.p)
	err := newError(res.error)
	if err.isError() {
		return 0, err.error()
	}
	return int32(res.value), nil
}

func (c *constant) Int64() (int64, error) {
	res := C.ergo_lib_constant_to_i64(c.p)
	err := newError(res.error)
	if err.isError() {
		return 0, err.error()
	}
	return int64(res.value), nil
}

func (c *constant) bytesLength() (int, error) {
	var returnNum C.ReturnNum_usize
	returnNum = C.ergo_lib_constant_bytes_len(c.p)
	err := newError(returnNum.error)

	if err.isError() {
		return 0, err.error()
	}
	size := C.ulong(returnNum.value)

	return int(size), nil
}

func (c *constant) Bytes() ([]byte, error) {
	bytesLength, bytesLengthErr := c.bytesLength()
	if bytesLengthErr != nil {
		return []byte{}, bytesLengthErr
	}

	output := C.malloc(C.uintptr_t(bytesLength))
	defer C.free(unsafe.Pointer(output))

	errPtr := C.ergo_lib_constant_to_bytes(c.p, (*C.uint8_t)(output))
	err := newError(errPtr)

	if err.isError() {
		return []byte{}, err.error()
	}

	result := C.GoBytes(unsafe.Pointer(output), C.int(bytesLength))
	return result, nil
}

func (c *constant) pointer() C.ConstantPtr {
	return c.p
}

func finalizeConstant(c *constant) {
	C.ergo_lib_constant_delete(c.p)
}
