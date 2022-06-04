package js

import (
	"errors"
	"fmt"
)

// Value is the internal representing of a js object
type Value struct {
	name  string // for debug
	value interface{}
	ref   Ref
}

// String return the string representing of a value
func (v *Value) String() string {
	return fmt.Sprintf("%s", v.value)
}

// Name is used for debugging
func (v *Value) Name() string {
	return v.name
}
func (v *Value) Ref() Ref {
	return v.ref
}

// Bytes return ref as []byte, error will return if type of ref is not Uint8Array
func (v *Value) Bytes() ([]byte, error) {
	slice, ok := v.value.([]byte)
	if !ok {
		return nil, errors.New("not Uint8Array")
	}
	return slice, nil
}
