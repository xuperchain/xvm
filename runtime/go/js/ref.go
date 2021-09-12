package js

import (
	"fmt"
	"unsafe"
)

const (
	nanHead = 0x7FF80000
)

const (
	// ValueUndefined is the ref of Undefined
	ValueUndefined = 0
)

const (
	// the type flags need to be in sync with wasm_exec.js
	typeFlagNone = iota
	typeFlagObject
	typeFlagString
	typeFlagSymbol
	typeFlagFunction
)

var (
	valueUndefined = Value{ref: 0}
	ValueNaN       = predefValue(0, typeFlagNone)
	ValueZero      = predefValue(1, typeFlagNone)
	ValueNull      = predefValue(2, typeFlagNone)
	ValueTrue      = predefValue(3, typeFlagNone)
	ValueFalse     = predefValue(4, typeFlagNone)
	ValueGlobal    = predefValue(5, typeFlagObject)
	ValueGo        = predefValue(6, typeFlagObject) // instance of the Go class in JavaScript
)

func predefValue(id uint32, typeFlag byte) *Value {
	return &Value{ref: (nanHead|Ref(typeFlag))<<32 | Ref(id)}
}

// ref represents the unique id of a js object
type Ref int64

// Number return ref as a number, if ref not a number, false will be returned
func (r Ref) Number() (int64, bool) {
	f := *(*float64)(unsafe.Pointer(&r))
	if f == f {
		return int64(f), true
	}
	return 0, false
}

// ID return the id of ref
func (r Ref) ID() int64 {
	id := uint32(r)
	return int64(id)
}

// String return the debug string of ref
func (r Ref) String() string {
	n, ok := r.Number()
	if ok {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("0x%x", int64(r))
}
