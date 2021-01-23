package types

import "github.com/zeta-io/zctl/errors"

type Type int

const (
	Int     Type = 1
	Int8	Type = 2
	Int16	Type = 3
	Int32	Type = 4
	Int64	Type = 5
	UInt	Type = 6
	UInt8	Type = 7
	UInt16	Type = 8
	UInt32	Type = 9
	UInt64	Type = 10
	Bool	Type = 11
	Float32	Type = 12
	Float64	Type = 13
	Time	Type = 14
	String	Type = 15
	Any		Type = 16
	Struct	Type = 17
	Array	Type = 18
	Map		Type = 19
)

var primitiveTypes = map[string]Type{
	"int":Int,
	"int8":Int8,
	"int16":Int16,
	"int32":Int32,
	"int64":Int64,
	"uint":UInt,
	"uint6":UInt8,
	"uint16":UInt16,
	"uint32":UInt32,
	"uint64":UInt64,
	"bool":Bool,
	"float32":Float32,
	"float64":Float64,
	"time":Time,
	"string":String,
}

func ParsePrimitive(s string) (Type, error) {
	if t, ok := primitiveTypes[s]; ok {
		return t, nil
	}
	return 0, errors.ErrTypesFormat
}

func (t Type) IsPrimitive() bool{
	return t == Int || t == Int8 ||
		t == Int16 || t == Int32 ||
		t == Int64 || t == UInt ||
		t == UInt8 || t == UInt16 ||
		t == UInt32 || t == UInt64 ||
		t == Bool || t == Float32 ||
		t == Float64 || t == Time ||
		t == String
}