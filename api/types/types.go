package types

import "github.com/zeta-io/zctl/errors"

type Type string

const (
	Int     Type = "int"
	Int8    Type = "int8"
	Int16   Type = "int16"
	Int32   Type = "int32"
	Int64   Type = "int64"
	UInt    Type = "uint"
	UInt8   Type = "uint8"
	UInt16  Type = "uint16"
	UInt32  Type = "uint32"
	UInt64  Type = "uint64"
	Bool    Type = "bool"
	Float32 Type = "float32"
	Float64 Type = "float64"
	Time    Type = "time"
	String  Type = "string"
)

var GoTypes = map[Type]string{
	Int:     "int",
	Int8:    "int8",
	Int16:   "int16",
	Int32:   "int32",
	Int64:   "int64",
	UInt:    "uint",
	UInt8:   "uint6",
	UInt16:  "uint16",
	UInt32:  "uint32",
	UInt64:  "uint64",
	Bool:    "bool",
	Float32: "float32",
	Float64: "float64",
	Time:    "time.Time",
	String:  "string",
}

func Parse(s string) (Type, error) {
	if _, ok := GoTypes[Type(s)]; !ok {
		return "", errors.WrapperError(errors.ErrFieldTypeError, s)
	}
	return Type(s), nil
}
